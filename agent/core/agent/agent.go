package agent

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhangxiaoyang/goDataAccess/agent/core/updater"
	"github.com/zhangxiaoyang/goDataAccess/agent/core/validator"
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
	u := updater.NewUpdater(this.dbPath, this.rulePath)
	v := validator.NewValidator(this.dbPath, this.rulePath)

	log.Println("started updater")
	u.Start()
	log.Println("finished updater")
	log.Println("started validator")
	v.Start()
	log.Println("finished validator")

	ticker := time.NewTicker(15 * time.Minute)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			log.Println("started updater")
			u.Start()
			log.Println("finished updater")
			log.Println("started validator")
			v.Start()
			log.Println("finished validator")
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
