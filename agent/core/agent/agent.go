package agent

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/pipeline"
	"github.com/zhangxiaoyang/goDataAccess/spider/plugin"
	"log"
	"os"
	"path"
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
	ruleFilePath := path.Join(this.rulePath, "update.json")
	dbFilePath := path.Join(this.dbPath, "agent.db")
	db, err := sql.Open("sqlite3", dbFilePath)
	log.Println(ruleFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tableName := `"update"`
	_, err = db.Exec(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(ip text, port text, tags text, source text); delete from %s",
		tableName,
		tableName,
	))
	if err != nil {
		log.Fatal(err)
	}

	e := engine.
		NewQuickEngine(path.Join(this.rulePath, "update.json")).
		GetEngine().
		SetPipeline(pipeline.NewSqlPipeline(db, tableName))

	var ok bool
	if ok = this.isAgentServerOK(); ok {
		e.AddPlugin(plugin.NewProxyPlugin())
	}
	log.Printf("started %s(isAgentServerOK: %v)\n", ruleFilePath, ok)
	e.Start()
}

func (this *Agent) isAgentServerOK() bool {
	_, err := common.NewProxy().GetOneProxy("http://example.com")
	if err != nil {
		return false
	}
	return true
}
