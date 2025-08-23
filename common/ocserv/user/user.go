package user

import (
	"bytes"
	"common/models"
	"common/ocserv"
	"common/pkg"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
)

type OcservUser struct{}

type OcservUserInterface interface {
	Create(username, group, password string, config *models.OcservUserConfig) error
	Lock(username string) (string, error)
	UnLock(username string) (string, error)
	Delete(username string) (string, error)
	CreateConfig(username string, config *models.OcservUserConfig) error
	DeleteConfig(username string) error
}

func NewOcservUser() *OcservUser {
	return &OcservUser{}
}

func (u *OcservUser) Create(username, group, password string, config *models.OcservUserConfig) error {
	args := []string{"-c", ocserv.OcpasswdPath, username}
	if group != "" {
		args = append([]string{"-g", group}, args...)
	}
	cmd := exec.Command(ocserv.OcpasswdExec, args...)
	cmd.Stdin = bytes.NewBufferString(password + "\n" + password + "\n")
	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		return err
	}
	if output == "" {
		return errors.New("create User Failed")
	}

	if config != nil {
		filename := filepath.Join(ocserv.ConfigUserBaseDir, username)

		var file *os.File

		// Open file with create, truncate, write-only flags and permission 0640
		file, err = os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
		if err != nil {
			return err
		}
		defer file.Close()

		err = pkg.ConfigWriter(file, pkg.ToMap(config))
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *OcservUser) Lock(username string) (string, error) {
	output, err := pkg.RunOcpasswd("-l", "-c", ocserv.OcpasswdPath, username)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (u *OcservUser) UnLock(username string) (string, error) {
	output, err := pkg.RunOcpasswd("-u", "-c", ocserv.OcpasswdPath, username)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (u *OcservUser) Delete(username string) (string, error) {
	output, err := pkg.RunOcpasswd("-d", "-c", ocserv.OcpasswdPath, username)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (u *OcservUser) CreateConfig(username string, config *models.OcservUserConfig) error {
	filename := pkg.ConfigFilePathCreator(username)
	// Open file with create, truncate, write-only flags and permission 0640
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0640)
	if err != nil {
		return err
	}
	err = pkg.ConfigWriter(file, pkg.ToMap(config))
	if err != nil {
		return err
	}
	return nil
}

func (u *OcservUser) DeleteConfig(username string) error {
	filename := pkg.ConfigFilePathCreator(username)
	if err := os.Remove(filename); err != nil {
		return err
	}
	return nil
}
