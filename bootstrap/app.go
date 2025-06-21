package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"github.com/web-page-analysis/util"
)

type AppConfig struct {
	Port int64 `json:"port"`
}

func initAppConfig() error {
	err := util.YamlReader(`bootstrap/config/app.yaml`, &AppConf)
	if err != nil {
		log.Errorf("init app config error: %v", err)
		return err
	}
	return nil
}
