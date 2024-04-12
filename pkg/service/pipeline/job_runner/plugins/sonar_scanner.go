package plugins

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/kubespace/kubespace/pkg/model/types"
	utilgit "github.com/kubespace/kubespace/pkg/third/git"
	"github.com/kubespace/kubespace/pkg/utils"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type SonarScanner struct{}

func (s SonarScanner) Executor(params *ExecutorParams) (Executor, error) {
	return newSonarScannerExecutor(params)
}

type sonarScannerParams struct {
	JobId       uint `json:"job_id"`
	WorkspaceId uint `json:"workspace_id"`

	CodeUrl      string        `json:"code_url"`
	CodeApiUrl   string        `json:"code_api_url"`
	CodeType     string        `json:"code_type"`
	CodeBranch   string        `json:"code_branch"`
	CodeCommitId string        `json:"code_commit_id"`
	CodeSecret   *types.Secret `json:"code_secret"`

	SonarScannerImage  PipelineResource `json:"sonar_scanner_image"`
	SonarScannerType   string           `json:"sonar_scanner_type"`
	SonarScannerFile   string           `json:"sonar_scanner_file"`
	SonarScannerScript string           `json:"sonar_scanner_script"`
	SonarQube          types.SonarQube  `json:"sonar_qube"`
}

type SonarScannerExecutorResult struct {
}

func newSonarScannerExecutor(params *ExecutorParams) (*SonarScannerExecutor, error) {
	var scanParams sonarScannerParams
	if err := utils.ConvertTypeByJson(params.Params, &scanParams); err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	sonarScannerExecutor := &SonarScannerExecutor{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		Params:     &scanParams,
		Logger:     params.Logger,
		Result:     &SonarScannerExecutorResult{},
	}
	codeDir := utils.GetCodeRepoName(scanParams.CodeUrl)
	if codeDir == "" {
		klog.Errorf("job=%d get empty code repo name", scanParams.JobId)
		return nil, fmt.Errorf("get empty code repo name")
	}
	sonarScannerExecutor.codeDir, _ = filepath.Abs(filepath.Join(params.RootDir, codeDir))

	return sonarScannerExecutor, nil
}

const (
	SonarScannerTypeNone   = "none"
	SonarScannerTypeFile   = "file"
	SonarScannerTypeScript = "script"
)

const (
	SonarScannerProjectFileName = "sonar-project.properties"
)

type SonarScannerExecutor struct {
	Logger                                 // 日志
	Params     *sonarScannerParams         // 插件参数
	Result     *SonarScannerExecutorResult // 扫描结果
	codeDir    string                      // 代码目录
	ctx        context.Context             // 上下文
	cancelFunc context.CancelFunc          // 上下文取消方法
	canceled   bool                        // 是否取消
}

func (s *SonarScannerExecutor) Execute() (any, error) {
	steps := []stepFunc{s.clone, s.scanner}
	for _, step := range steps {
		err := step()
		if s.canceled {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
	}
	return s.Result, nil
}

func (s *SonarScannerExecutor) Cancel() error {
	s.canceled = true
	s.cancelFunc()
	return nil
}

// clone 克隆代码
func (s *SonarScannerExecutor) clone() error {
	if err := os.RemoveAll(s.codeDir); err != nil {
		return err
	}
	s.Log("git clone %v", s.Params.CodeUrl)
	time.Sleep(1)
	var err error

	client, err := utilgit.NewClient(s.Params.CodeType, s.Params.CodeApiUrl, s.Params.CodeSecret)
	if err != nil {
		return err
	}
	r, err := client.Clone(s.ctx, s.codeDir, false, &git.CloneOptions{
		URL:      s.Params.CodeUrl,
		Progress: s.Logger,
	})
	if err != nil {
		s.Log("克隆代码仓库失败：%v", err)
		klog.Errorf("job=%d clone %s error: %v", s.Params.JobId, s.Params.CodeUrl, err)
		return fmt.Errorf("git clone %s error: %v", s.Params.CodeUrl, err)
	}
	w, err := r.Worktree()
	if err != nil {
		s.Log("克隆代码仓库失败：%v", err)
		klog.Errorf("job=%d clone %s error: %v", s.Params.JobId, s.Params.CodeUrl, err)
		return fmt.Errorf("git clone %s error: %v", s.Params.CodeUrl, err)
	}
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(s.Params.CodeCommitId),
	})
	if err != nil {
		s.Log("git checkout %s 失败：%v", s.Params.CodeCommitId, err)
		klog.Errorf("job=%d git checkout %s error: %v", s.Params.JobId, s.Params.CodeCommitId, err)
		return fmt.Errorf("git checkout %s error: %v", s.Params.CodeCommitId, err)
	}
	return nil
}

func (s *SonarScannerExecutor) scanner() error {
	if s.Params.SonarScannerFile == SonarScannerTypeNone {
		s.Log("跳过代码检测")
		return nil
	}

	if s.Params.SonarScannerImage.Value == "" {
		s.Log("代码扫描镜像为空，请检查流水线配置")
		return fmt.Errorf("sonar scanner image is empty")
	}

	if s.Params.SonarQube.HostUrl == "" || s.Params.SonarQube.Token == "" {
		s.Log("SonarQube服务器配置不正确，请检查平台配置")
		return fmt.Errorf("sonar qube settings is empty")
	}

	// 默认读仓库中的文件
	sonarScannerType := SonarScannerTypeFile
	if s.Params.SonarScannerType != "" {
		sonarScannerType = s.Params.SonarScannerType
	}
	switch sonarScannerType {
	case SonarScannerTypeFile:
		sonarScannerProjectFileName := SonarScannerProjectFileName
		if s.Params.SonarScannerFile != "" {
			sonarScannerProjectFileName = s.Params.SonarScannerFile
		}
		_, err := os.Stat(filepath.Join(s.codeDir, sonarScannerProjectFileName))
		if err != nil {
			s.Log("代码仓库中不存在sonar配置文件")
			return fmt.Errorf("sonar scanner project file %s error: %v", SonarScannerProjectFileName, err)
		}
	case SonarScannerTypeScript:
		if s.Params.SonarScannerScript == "" {
			s.Log("代码扫描的配置文件为空")
			return nil
		}
		err := os.WriteFile(
			filepath.Join(s.codeDir, SonarScannerProjectFileName),
			[]byte(s.Params.SonarScannerScript),
			0666,
		)
		if err != nil {
			s.Log("写配置文件%s错误：%v", SonarScannerProjectFileName, err)
			klog.Errorf("job=%d write scan error: %v", s.Params.JobId, err)
			return fmt.Errorf("write scan file error: %s", err.Error())
		}
	}

	dockerRunCmd := fmt.Sprintf(
		"docker run --net=host --rm -v %s:/usr/src -e SONAR_HOST_URL=%s -e SONAR_TOKEN=%s %s",
		s.codeDir,
		s.Params.SonarQube.HostUrl,
		s.Params.SonarQube.Token,
		s.Params.SonarScannerImage.Value,
	)
	klog.Infof("job=%d sonar scanner cmd: %s", s.Params.JobId, dockerRunCmd)
	cmd := exec.CommandContext(s.ctx, "bash", "-xc", dockerRunCmd)
	cmd.Stdout = s.Logger
	cmd.Stderr = s.Logger
	if err := cmd.Run(); err != nil {
		klog.Errorf("job=%d scan error: %v", s.Params.JobId, err)
		s.Log("scan error: %s", err.Error())
		return fmt.Errorf("scan code error: %v", err)
	}
	return nil
}
