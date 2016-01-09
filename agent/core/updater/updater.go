package updater

import (
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/agent/util"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/pipeline"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
	"log"
	"path"
)

type Updater struct {
	dbPath   string
	rulePath string
}

func NewUpdater(dbPath string, rulePath string) *Updater {
	return &Updater{
		dbPath:   dbPath,
		rulePath: rulePath,
	}
}

func (this *Updater) Start() {
	ruleFilePath := path.Join(this.rulePath, "update.json")
	dbFilePath := path.Join(this.dbPath, "agent.db")
	tableName := `"update"`

	db, err := util.InitTable(fmt.Sprintf(
		"PRAGMA journal_mode = WAL; CREATE TABLE IF NOT EXISTS %s(ip TEXT, port TEXT, source TEXT, level INTEGER)",
		tableName,
	), dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	level := util.GetLastLevel(tableName, db) + 1

	e := engine.
		NewQuickEngine(ruleFilePath).
		GetEngine().
		AddPlugin(util.NewAddLevelPlugin(level)).
		SetPipeline(pipeline.NewSqlPipeline(db, tableName))

	var ok bool
	if ok = this.isAgentServerOK(); ok {
		e.AddPlugin(plugin.NewProxyPlugin())
	}
	log.Printf("started %s(isAgentServerOK: %v)\n", ruleFilePath, ok)
	e.Start()
}

func (this *Updater) isAgentServerOK() bool {
	_, err := common.NewProxy().GetOneProxy("http://example.com")
	if err != nil {
		return false
	}
	return true
}
