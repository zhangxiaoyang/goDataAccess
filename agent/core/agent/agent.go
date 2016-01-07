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

	ticker := 15 * time.Minute
	for {
		log.Println("started updater")
		u.Start()
		log.Println("finished updater")
		log.Println("started validator")
		v.Start()
		log.Println("finished validator")

		time.Sleep(ticker)
	}
}
