package sonar_qube

import (
	"github.com/kubespace/kubespace/pkg/core/code"
	"github.com/kubespace/kubespace/pkg/core/errors"
	"github.com/kubespace/kubespace/pkg/model"
	"github.com/kubespace/kubespace/pkg/server/api/api"
	"github.com/kubespace/kubespace/pkg/server/config"
	"github.com/kubespace/kubespace/pkg/utils"
)

func ListHandler(conf *config.ServerConfig) api.Handler {
	return &listHandler{models: conf.Models}
}

type listHandler struct {
	models *model.Models
}

func (h *listHandler) Auth(c *api.Context) (bool, *api.AuthPerm, error) {
	return true, nil, nil
}

func (h *listHandler) Handle(c *api.Context) *utils.Response {
	sonarQubes, err := h.models.SonarQubeManager.List()
	if err != nil {
		return c.ResponseError(errors.New(code.ParamsError, err))
	}
	var data []map[string]interface{}

	for _, sonarQube := range sonarQubes {
		data = append(data, map[string]interface{}{
			"id":          sonarQube.ID,
			"name":        sonarQube.Name,
			"description": sonarQube.Description,
			"host_url":    sonarQube.HostUrl,
			"token":       sonarQube.Token,
			"create_user": sonarQube.CreateUser,
			"update_user": sonarQube.UpdateUser,
			"create_time": sonarQube.CreateTime,
			"update_time": sonarQube.UpdateTime,
		})
	}
	return c.ResponseOK(data)
}
