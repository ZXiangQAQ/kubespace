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
	"time"
)

func UpdateHandler(conf *config.ServerConfig) api.Handler {
	return &updateHandler{models: conf.Models}
}

type updateHandler struct {
	models *model.Models
}

func (h *updateHandler) Auth(c *api.Context) (bool, *api.AuthPerm, error) {
	return true, &api.AuthPerm{
		Scope:   types.ScopePlatform,
		ScopeId: 0,
		Role:    types.RoleEditor,
	}, nil
}

func (h *updateHandler) Handle(c *api.Context) *utils.Response {
	var body createSonarQubeRequest
	if err := c.ShouldBind(&body); err != nil {
		return c.ResponseError(errors.New(code.ParamsError, err))
	}
	sonarQubeId, err := utils.ParseUint(c.Param("id"))
	if err != nil {
		return c.ResponseError(errors.New(code.ParamsError, err))
	}
	sonarQube, err := h.models.SonarQubeManager.Get(sonarQubeId)
	if err != nil {
		return c.ResponseError(errors.New(code.DataNotExists, "获取SonarQube失败: "+err.Error()))
	}
	sonarQube.Description = body.Description
	sonarQube.HostUrl = body.HostUrl
	sonarQube.Token = body.Token
	sonarQube.UpdateUser = c.User.Name
	sonarQube.UpdateTime = time.Now()
	_, err = h.models.SonarQubeManager.Update(sonarQube)
	if err != nil {
		err = errors.New(code.DBError, "更新SonarQube失败: "+err.Error())
	}
	resp := c.ResponseError(err)
	c.CreateAudit(&types.AuditOperate{
		Operation:            types.AuditOperationUpdate,
		OperateDetail:        fmt.Sprintf("更新SonarQube：%s", sonarQube.Name),
		Scope:                types.ScopePlatform,
		ResourceId:           sonarQube.ID,
		ResourceType:         types.AuditResourceSonarQube,
		ResourceName:         sonarQube.Name,
		Code:                 resp.Code,
		Message:              resp.Msg,
		OperateDataInterface: body,
	})
	return resp
}
