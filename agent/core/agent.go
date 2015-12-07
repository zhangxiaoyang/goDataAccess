package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/agent/util"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/processer"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

type Agent struct {
	ruleDir            string
	dbDir              string
	candidateAgentPath string
}

func NewAgent(ruleDir string, dbDir string) *Agent {
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		os.Mkdir(dbDir, os.ModePerm)
	}

	return &Agent{
		ruleDir:            ruleDir,
		dbDir:              dbDir,
		candidateAgentPath: path.Join(dbDir, "candidate.json"),
	}
}

func (this *Agent) Update() {
	file, _ := os.Create(this.candidateAgentPath)
	defer file.Close()

	fileInfos, err := ioutil.ReadDir(this.ruleDir)
	if err != nil {
		log.Printf("error %s\n", err)
		return
	}

	updateRulePaths := []string{}
	for _, f := range fileInfos {
		updateRulePath := path.Join(this.ruleDir, f.Name())
		if this.isUpdateRule(f.Name()) {
			updateRulePaths = append(updateRulePaths, updateRulePath)
		} else {
			log.Printf("skip %s\n", updateRulePath)
		}
	}

	var wg sync.WaitGroup
	wg.Add(len(updateRulePaths))
	for _, updateRulePath := range updateRulePaths {
		go func() {
			defer wg.Done()
			log.Printf("started %s\n", updateRulePath)
			engine.NewQuickEngine(updateRulePath).SetOutputFile(file).Start()
			log.Printf("finished %s\n", updateRulePath)
		}()
	}
	wg.Wait()
}

func (this *Agent) Validate(validateUrl string, succ string) {
	if !strings.HasPrefix(validateUrl, "http://") {
		validateUrl = "http://" + validateUrl
	}
	domain := util.ExtractDomain(validateUrl)
	validAgentPath := path.Join(this.dbDir, fmt.Sprintf("valid.%s.json", domain))
	validateRulePath := path.Join(this.ruleDir, "validate.json")

	file, _ := os.Create(validAgentPath)
	defer file.Close()

	engine.
		NewQuickEngine(validateRulePath).
		GetEngine().
		SetStartUrls(this.readAllCandidate()).
		SetDownloader(util.NewValidateDownloader(validateUrl, succ)).
		SetPipeline(util.NewFilePipeline(file)).
		SetProcesser(processer.NewLazyProcesser()).
		Start()
}

func (this *Agent) isUpdateRule(fileName string) bool {
	if strings.HasPrefix(fileName, "update.") {
		return true
	}
	return false
}

func (this *Agent) readAllCandidate() []string {
	file, err := os.Open(this.candidateAgentPath)
	if err != nil {
		log.Printf("error %s", err)
		return []string{}
	}
	defer file.Close()

	r := bufio.NewReader(file)
	addrs := map[string]bool{}
	for {
		line, err := r.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}

		addr := &util.Addr{}
		json.Unmarshal([]byte(line), addr)
		addrs[addr.Serialize()] = true
	}

	keys := []string{}
	for k := range addrs {
		keys = append(keys, k)
	}
	return keys
}
