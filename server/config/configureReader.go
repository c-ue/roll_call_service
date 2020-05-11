package config

import (
	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"roll_call_service/server/logger"
)

type (
	SERVER struct {
		IP    string
		PORT  string
		DEBUG bool
	}

	MEMBER struct {
		NAME       string
		CalendarId string
		GROUP      int
		PRIVILEGE  int
	}

	EVENT struct {
		NAME string
		TIME []int
	}

	Config struct {
		SERVER  SERVER
		MEMBERS []MEMBER
		EVENTS  []EVENT
	}
)

func ReadConfig(fileName string) Config {
	var conf Config
	var log *zap.Logger = logger.Console()
	if _, err := toml.DecodeFile(fileName, &conf); err != nil {
		log.Debug("toml file reader error")
		panic(err)
	}
	return conf
}
