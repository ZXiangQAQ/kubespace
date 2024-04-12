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

func DeleteHandler(conf *config.ServerConfig) api.Handler {
	return &deleteHandler{models: conf.Models}
}

type deleteHandler struct {
	models *model.Models
}

func (h *deleteHandler) Auth(c *api.Context) (bool, *api.AuthPerm, error) {
	return true, &api.AuthPerm{
		Scope:   types.ScopePlatform,
		ScopeId: 0,
		Role:    types.RoleEditor,
	}, nil
}

func (h *deleteHandler) Handle(c *api.Context) *utils.Response {
	sonarQubeId, err := utils.ParseUint(c.Param("id"))
	if err != nil {
		return c.ResponseError(errors.New(code.ParamsError, err))
	}
	sonarQube, err := h.models.SonarQubeManager.Get(sonarQubeId)
	if err != nil {
		return c.ResponseError(errors.New(code.DataNotExists, "获取SonarQube失败: "+err.Error()))
	}
	err = h.models.SonarQubeManager.Delete(sonarQube)
	if err != nil {
		err = errors.New(code.DBError, "删除SonarQube失败: "+err.Error())
	}
	resp := c.ResponseError(err)
	c.CreateAudit(&types.AuditOperate{
		Operation:            types.AuditOperationDelete,
		OperateDetail:        fmt.Sprintf("删除SonarQube：%s", sonarQube.Name),
		Scope:                types.ScopePlatform,
		ResourceId:           sonarQube.ID,
		ResourceType:         types.AuditResourceSonarQube,
		ResourceName:         sonarQube.Name,
		Code:                 resp.Code,
		Message:              resp.Msg,
		OperateDataInterface: nil,
	})
	return resp
}
