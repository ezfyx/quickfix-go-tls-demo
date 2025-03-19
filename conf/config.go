package conf

import (
	"bytes"
	"github.com/quickfixgo/quickfix"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"quickfix-go-tls-demo/utils"
)

var (
	TestSettings *quickfix.Settings
)

func loadSettings(file string) (*quickfix.Settings, error) {
	cfg, err := os.Open(file)
	if err != nil {
		logrus.Errorf("error opening cfg %v, %v", file, err)
		return nil, err
	}
	defer cfg.Close()
	cfgStr, err := io.ReadAll(cfg)
	if err != nil {
		logrus.Errorf("error reading cfg: %s,", err)
		return nil, err
	}
	appSettings, err := quickfix.ParseSettings(bytes.NewReader(cfgStr))
	if err != nil {
		logrus.Errorf("error reading cfg: %s,", err)
		return nil, err
	}
	return appSettings, nil
}

func LoadConfig() error {
	logrus.SetFormatter(&utils.TextFormatter{
		UppercaseFirstMsgLetter: false,
	})

	settings, err := loadSettings(path.Join("config", "sample.cfg"))
	if err != nil {
		return err
	}
	TestSettings = settings

	return nil
}

func GetMdSettings() *quickfix.Settings {
	return TestSettings
}

func GetOmSettings() *quickfix.Settings {
	return TestSettings
}
