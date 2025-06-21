package bootstrap

var (
	AppConf      AppConfig
	OutboundConf OutboundConfig
)

type Config struct {
	AppConfig    AppConfig
	OutboundConf OutboundConfig
}

func InitConfig() (conf Config, err error) {
	err = initAppConfig()
	if err != nil {
		return conf, err
	}

	err = initOutboundConfig()
	if err != nil {
		return conf, err
	}
	conf = Config{
		AppConfig:    AppConf,
		OutboundConf: OutboundConf,
	}
	return conf, nil
}
