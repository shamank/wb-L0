package main

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"github.com/shamank/wb-l0/app/config"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

const configDir = "configs"
const modelsDir = "./app/cmd/publish/models"

func main() {
	cfg, err := config.ConfigInit(configDir)
	if err != nil {
		logrus.Fatalf("error occurred in initial config: %s", err.Error())
		return
	}

	sc, err := stan.Connect(cfg.Nats.ClusterID, "publish-client", stan.NatsURL(cfg.Nats.URL))

	defer sc.Close()

	err = filepath.Walk(modelsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".json" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			var order interface{}

			if err := json.Unmarshal(data, &order); err != nil {
				return err
			}

			err = sc.Publish(cfg.Nats.Channel, data)
			if err != nil {

				return err
			}
			logrus.Infof("success to publish model: %s", path)
		}
		return nil
	})
	if err != nil {
		logrus.Fatalf("error with walking files: %s", err.Error())
		return
	}

}
