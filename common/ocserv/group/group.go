package group

import (
	"common/models"
	"common/ocserv"
	"common/pkg"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

type OcservGroup struct{}

type OcservGroupInterface interface {
	Create(name string, config *models.OcservGroupConfig) error
	Delete(name string) error
	DefaultsGroup() (*models.OcservGroupConfig, error)
	UpdateDefaultsGroup(config *models.OcservGroupConfig) error
}

func NewOcservGroup() *OcservGroup {
	return &OcservGroup{}
}

func (g *OcservGroup) Create(name string, config *models.OcservGroupConfig) error {
	filename := filepath.Join(ocserv.ConfigGroupBaseDir, name)

	// Open file with create, truncate, write-only flags and permission 0640
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	err = pkg.ConfigWriter(file, pkg.ToMap(config))
	if err != nil {
		return err
	}
	return nil
}

func (g *OcservGroup) Delete(name string) error {
	filename := filepath.Join(ocserv.ConfigGroupBaseDir, name)
	if err := os.Remove(filename); err != nil {
		return err
	}
	users, err := pkg.GetUsersByGroup(name)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	// Reset group assignment for each user to default (empty)
	for _, user := range users {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			_, _ = exec.Command(ocserv.OcpasswdExec, "-g", "", "-c", ocserv.OcpasswdPath, u).CombinedOutput()
		}(user)
	}
	wg.Wait()

	return nil
}

func (g *OcservGroup) DefaultsGroup() (*models.OcservGroupConfig, error) {
	configInterface, err := pkg.ParseOcservConfigFile(ocserv.DefaultGroupFile)
	if err != nil {
		return nil, err
	}

	configJson, err := json.Marshal(configInterface)
	if err != nil {
		return nil, err
	}

	var config models.OcservGroupConfig

	if err = json.Unmarshal(configJson, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (g *OcservGroup) UpdateDefaultsGroup(config *models.OcservGroupConfig) error {
	file, err := os.OpenFile(ocserv.DefaultGroupFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	defer file.Close()

	err = pkg.ConfigWriter(file, pkg.ToMap(config))
	if err != nil {
		return err
	}
	return nil
}
