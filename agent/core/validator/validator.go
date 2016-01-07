package validator

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhangxiaoyang/goDataAccess/agent/util"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/pipeline"
	"log"
	"path"
)

type Validator struct {
	dbPath   string
	rulePath string
}

func NewValidator(dbPath string, rulePath string) *Validator {
	return &Validator{
		dbPath:   dbPath,
		rulePath: rulePath,
	}
}

func (this *Validator) Start() {
	ruleFilePath := path.Join(this.rulePath, "validate.json")
	dbFilePath := path.Join(this.dbPath, "agent.db")
	updateTableName := `"update"`
	validateTableName := `"validate"`

	db, err := util.InitTable(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(ip TEXT, port TEXT, domain TEXT, level INTEGER)",
		validateTableName,
	), dbFilePath)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	level := util.GetLastLevel(validateTableName, db) + 1

	qe := engine.NewQuickEngine(ruleFilePath)
	reqs := this.genRequests(qe.GetQuickEngineConfig().StartUrls, updateTableName, level, db)

	qe.
		GetEngine().
		SetStartRequests(reqs).
		AddPlugin(util.NewModifyResponsePlugin()).
		AddPlugin(util.NewAddLevelPlugin(level)).
		SetPipeline(pipeline.NewSqlPipeline(db, validateTableName)).
		Start()
}

func (this *Validator) genRequests(urls []string, tableName string, level int, db *sql.DB) []*common.Request {
	proxies := this.getProxyByLevel(tableName, db)
	reqs := []*common.Request{}
	for _, url := range urls {
		for _, proxy := range proxies {
			req := common.NewRequest(url)
			req.ProxyUrl = proxy
			reqs = append(reqs, req)
		}
	}
	return reqs
}

func (this *Validator) getProxyByLevel(tableName string, db *sql.DB) []string {
	proxies := []string{}
	level := util.GetLastLevel(tableName, db)

	rows, err := db.Query(fmt.Sprintf(
		"SELECT ip, port FROM %s WHERE level=%d",
		tableName,
		level,
	))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ip string
			var port string
			err := rows.Scan(&ip, &port)
			if err == nil {
				proxies = append(proxies, fmt.Sprintf("%s:%s", ip, port))
			}
		}
	} else {
		log.Fatal(err)
	}
	return proxies
}
