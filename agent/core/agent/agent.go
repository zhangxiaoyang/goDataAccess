package agent

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhangxiaoyang/goDataAccess/agent/core/updater"
	"log"
	"os"
	"time"
)

type Agent struct {
	dbPath   string
	rulePath string
}

func NewAgent(dbPath string, rulePath string) *Agent {
	os.Remove(dbPath)
	os.Mkdir(dbPath, os.ModePerm)

	return &Agent{
		dbPath:   dbPath,
		rulePath: rulePath,
	}
}

func (this *Agent) Start() {
	log.Println("started updater")
	updater.NewUpdater(this.dbPath, this.rulePath).Start()
	log.Println("finished updater")
	log.Println("started validator")
	//validator.NewValidator().Start()
	log.Println("finished validator")

	updaterTicker := time.NewTicker(24 * time.Hour)
	validatorTicker := time.NewTicker(30 * time.Minute)
	quit := make(chan struct{})
	for {
		select {
		case <-updaterTicker.C:
			log.Println("started updater")
			updater.NewUpdater(this.dbPath, this.rulePath).Start()
			log.Println("finished updater")
		case <-validatorTicker.C:
			log.Println("started validator")
			//validator.NewValidator().Start()
			log.Println("finished validator")
		case <-quit:
			updaterTicker.Stop()
			//validatorTicker.Stop()
			return
		}
	}
}
