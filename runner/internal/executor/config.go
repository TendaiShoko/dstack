package executor

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dstackai/dstackai/runner/consts"
	"github.com/dstackai/dstackai/runner/internal/common"
	"github.com/dstackai/dstackai/runner/internal/gerrors"
	"github.com/dstackai/dstackai/runner/internal/log"
	"github.com/dstackai/dstackai/runner/internal/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Id         string           `yaml:"id"`
	Resources  *models.Resource `yaml:"resources"`
	Hostname   *string          `yaml:"hostname"`
	ExposePort *string          `yaml:"expose_ports,omitempty"`
}

func (ex *Executor) loadConfig() error {
	thePathConfig := filepath.Join(common.HomeDir(), consts.DSTACK_DIR_PATH, consts.RUNNER_FILE_NAME)
	if _, err := os.Stat(thePathConfig); os.IsNotExist(err) {
		log.Error(context.Background(), "Failed to load config", "err", gerrors.Wrap(err))
		return err
	}
	theConfigFile, err := ioutil.ReadFile(thePathConfig)
	if err != nil {
		log.Error(context.Background(), "Unexpected error, please try to rerun", "err", gerrors.Wrap(err))
		return err
	}
	if err = yaml.Unmarshal(theConfigFile, &ex.config); err != nil {
		log.Error(context.Background(), "Config file is corrupted or does not exists", "err", gerrors.Wrap(err))
		return err
	}
	if ex.config.Resources == nil {
		ex.config.Resources, err = updateResourceInfo()
		if err != nil {
			return gerrors.Wrap(err)
		}
		err = ex.saveConfig()
		if err != nil {
			return gerrors.Wrap(err)
		}
	}
	return nil
}
func (ex *Executor) saveConfig() error {
	log.Info(context.Background(), "Save runner's config")
	thePathConfig := filepath.Join(common.HomeDir(), consts.DSTACK_DIR_PATH, consts.RUNNER_FILE_NAME)
	theConfigFile, err := yaml.Marshal(ex.config)
	if err != nil {
		log.Error(context.Background(), "Unexpected error, please try to rerun", "err", gerrors.Wrap(err))
		return err
	}
	err = ioutil.WriteFile(thePathConfig, theConfigFile, 0o644)
	if err != nil {
		log.Error(context.Background(), "Unexpected error, please try to rerun", "err", gerrors.Wrap(err))
		return err
	}
	return nil
}
func (c *Config) ExposePorts() (defaultStart int, defaultEnd int) {
	defaultStart = consts.EXPOSE_PORT_START
	defaultEnd = consts.EXPOSE_PORT_END
	if c.ExposePort != nil {
		sPorts := strings.Split(*c.ExposePort, "-")
		if len(sPorts) != 2 {
			logrus.Errorf("Expose_port is not the correct format. Requires INT-INT. Set to default values")
			return
		}
		startPort, err := strconv.Atoi(sPorts[0])
		if err != nil {
			logrus.Errorf("Expose_port is not the correct format. Requires INT-INT. Set to default values")
			return
		}
		endPort, err := strconv.Atoi(sPorts[1])
		if err != nil {
			logrus.Errorf("Expose_port is not the correct format. Requires INT-INT. Set to default values")
			return
		}
		if startPort > endPort {
			logrus.Errorf("Expose_port is not the correct format. Requires INT-INT. Set to default values")
			return
		}
		return startPort, endPort
	}
	return
}
