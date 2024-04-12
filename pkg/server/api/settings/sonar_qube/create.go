package sonar_qube

import (
	"fmt"
	"github.com/kubespace/kubespace/pkg/core/code"
	"github.com/kubespace/kubespace/pkg/core/errors"
	"github.com/kubespace/kubespace/pkg/model"
	"github.com/kubespace/kubespace/pkg/model/types"
	"github.com/kubespace/kubespace/pkg/server/api/api"
	"github.com/kubespace/kubespace/pkg/server/config"
	"github.com/kubespace/kubespace/pkg/utils"
)

func CreateHandler(conf *config.ServerConfig) api.Handler {
	return &createHandler{models: conf.Models}
}

type createSonarQubeRequest struct {
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
	HostUrl     string `json:"host_url" form:"host_url"`
	Token       string `json:"token" form:"token"`
}

type createHandler struct {
	models *model.Models
}

func (h *createHandler) Auth(c *api.Context) (bool, *api.AuthPerm, error) {
	return true, &api.AuthPerm{
		Scope:   types.ScopePlatform,
		ScopeId: 0,
		Role:    types.RoleEditor,
	}, nil
}

func (h *createHandler) Handle(c *api.Context) *utils.Response {
	var req createSonarQubeRequest
	err := c.ShouldBind(&req)
	if err != nil {
		return c.ResponseError(errors.New(code.ParamsError, err))
	}

	instance := &types.SettingsSonarQube{
		Name:        req.Name,
		Description: req.Description,
		HostUrl:     req.HostUrl,
		Token:       req.Token,
		CreateUser:  c.User.Name,
		UpdateUser:  c.User.Name,
	}
	_, err = h.models.SonarQubeManager.Create(instance)
	if err != nil {
		err = errors.New(code.DBError, "创建SonarQube失败: "+err.Error())
	}
	resp := c.ResponseError(err)
	c.CreateAudit(&types.AuditOperate{
		Operation:            types.AuditOperationCreate,
		OperateDetail:        fmt.Sprintf("创建SonarQube：%s", req.Name),
		Scope:                types.ScopePlatform,
		ResourceId:           instance.ID,
		ResourceType:         types.AuditResourceSonarQube,
		ResourceName:         instance.Name,
		Code:                 resp.Code,
		Message:              resp.Msg,
		OperateDataInterface: instance,
	})
	return resp
}
