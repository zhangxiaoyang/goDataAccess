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
	"math/rand"
	"path"
	"time"
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
		"PRAGMA journal_mode = WAL; CREATE TABLE IF NOT EXISTS %s(ip TEXT, port TEXT, domain TEXT, level INTEGER)",
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
	proxies, _ := util.GetLastProxies(tableName, db)
	reqs := []*common.Request{}
	for _, url := range urls {
		for _, proxy := range proxies {
			req := common.NewRequest(url)
			req.ProxyUrl = proxy
			reqs = append(reqs, req)
		}
	}

	rand.Seed(time.Now().Unix())
	this.Shuffle(reqs)
	return reqs
}

func (this *Validator) Shuffle(reqs []*common.Request) {
	for i := range reqs {
		j := rand.Intn(i + 1)
		reqs[i], reqs[j] = reqs[j], reqs[i]
	}
}
