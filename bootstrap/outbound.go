package bootstrap

import (
	log "github.com/sirupsen/logrus"
	"github.com/web-page-analysis/util"
)

type OutboundConfig struct {
	DialTimeout   int64 `yaml:"dial_timeout"`
	RemoteTimeout int64 `yaml:"remote_timeout"`
}

func initOutboundConfig() error {
	err := util.YamlReader(`bootstrap/config/outbound.yaml`, &OutboundConf)
	if err != nil {
		log.Errorf("init app config error: %v", err)
		return err
	}
	return nil
}
