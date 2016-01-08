package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zhangxiaoyang/goDataAccess/agent/util"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"
)

type Server struct {
	port      int
	tableName string
	db        *sql.DB
}

func NewServer(dbPath string, port int) *Server {
	dbFilePath := path.Join(dbPath, "agent.db")
	validateTableName := `"validate"`
	db, err := util.InitTable(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(ip TEXT, port TEXT, domain TEXT, level INTEGER)",
		validateTableName,
	), dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		log.Println("Closing database...")
		db.Close()
		os.Exit(0)
	}()

	return &Server{
		port:      port,
		tableName: validateTableName,
		db:        db,
	}
}

func (this *Server) Start() {
	http.HandleFunc("/", this.getAllProxies)
	http.HandleFunc("/getAllProxies", this.getAllProxies)
	http.HandleFunc("/getOneProxy", this.getOneProxy)
	log.Printf("Served at port %d\n", this.port)
	http.ListenAndServe(fmt.Sprintf(":%d", this.port), nil)
}

func (this *Server) getAllProxies(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	var domain string
	if url == "" {
		domain = "m.baidu.com"
	} else {
		domain = util.ExtractDomain(url)
	}

	proxies, domain, level := util.GetLastProxiesByDomain(this.tableName, domain, this.db)
	var result map[string]interface{}
	if len(proxies) == 0 {
		result = map[string]interface{}{
			"num":    0,
			"result": proxies,
		}
	} else {
		result = map[string]interface{}{
			"num":          len(proxies),
			"level":        level,
			"domain_match": domain,
			"result":       proxies,
		}
	}
	jsonResp, _ := json.Marshal(result)
	io.WriteString(w, string(jsonResp))
}

func (this *Server) random(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (this *Server) getOneProxy(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	var domain string
	if url == "" {
		domain = "m.baidu.com"
	} else {
		domain = util.ExtractDomain(url)
	}

	proxies, domain, level := util.GetLastProxiesByDomain(this.tableName, domain, this.db)
	var result map[string]interface{}
	if len(proxies) == 0 {
		result = map[string]interface{}{
			"num":    0,
			"result": proxies,
		}
	} else {
		result = map[string]interface{}{
			"num":          len(proxies),
			"level":        level,
			"domain_match": domain,
			"result":       proxies,
		}
	}
	jsonResp, _ := json.Marshal(result)
	io.WriteString(w, string(jsonResp))
}
