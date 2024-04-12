package settings

import (
	"github.com/kubespace/kubespace/pkg/model/manager"
	"github.com/kubespace/kubespace/pkg/model/types"
	"gorm.io/gorm"
)

type SonarQubeManager struct {
	*manager.CommonManager
}

func NewSonarQubeManager(db *gorm.DB) *SonarQubeManager {
	return &SonarQubeManager{
		CommonManager: manager.NewCommonManager(nil, db, "", false),
	}
}

func (m *SonarQubeManager) Create(sonarQube *types.SettingsSonarQube) (*types.SettingsSonarQube, error) {
	err := m.DB.Create(sonarQube).Error
	return sonarQube, err
}

func (m *SonarQubeManager) Delete(sonarQube *types.SettingsSonarQube) error {
	return m.DB.Delete(sonarQube).Error
}

func (m *SonarQubeManager) Update(sonarQube *types.SettingsSonarQube) (*types.SettingsSonarQube, error) {
	err := m.DB.Save(sonarQube).Error
	return sonarQube, err
}

func (m *SonarQubeManager) Get(id uint) (*types.SettingsSonarQube, error) {
	var sonarQube types.SettingsSonarQube
	err := m.DB.First(&sonarQube, id).Error
	return &sonarQube, err
}

func (m *SonarQubeManager) List() ([]types.SettingsSonarQube, error) {
	var sonarQubes []types.SettingsSonarQube
	err := m.DB.Find(&sonarQubes).Error
	return sonarQubes, err
}
