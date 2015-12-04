package agent

import (
	"bufio"
	"encoding/json"
	"github.com/zhangxiaoyang/goDataAccess/agent/core/validator"
	"github.com/zhangxiaoyang/goDataAccess/agent/util"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"github.com/zhangxiaoyang/goDataAccess/spider/core/engine"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Agent struct {
	dbDir            string
	ruleDir          string
	updateFilePath   string
	validateFilePath string
	pollingTime      time.Duration
	resourceManager  *common.ResourceManager
}

func NewAgent(dbDir string, ruleDir string, concurrency int) *Agent {
	return &Agent{
		dbDir:            dbDir,
		ruleDir:          ruleDir,
		updateFilePath:   dbDir + "agent.update.json",
		validateFilePath: dbDir + "agent.validate.json",
		pollingTime:      200 * time.Millisecond,
		resourceManager:  common.NewResourceManager(concurrency),
	}
}

func (this *Agent) Update() {
	file, _ := os.Create(this.updateFilePath)
	defer file.Close()

	if fileInfos, err := ioutil.ReadDir(this.ruleDir); err == nil {
		for _, f := range fileInfos {
			log.Printf("crawling %s\n", this.ruleDir+f.Name())
			engine.NewQuickEngine(this.ruleDir + f.Name()).SetOutputFile(file).Start()
			log.Printf("finished %s\n", this.ruleDir+f.Name())
		}
	}
}

func (this *Agent) Validate(validateUrl string, succ string) {
	validator := validator.NewValidator()
	if validateUrl != "" {
		validator.SetValidateUrl(validateUrl)
	}
	if succ != "" {
		validator.SetSucc(succ)
	}
	addrs := map[string]bool{}

	inFile, err := os.Open(this.updateFilePath)
	if err != nil {
		log.Printf("error %s\n", err)
		return
	}
	defer inFile.Close()
	outFile, err := os.Create(this.validateFilePath)
	if err != nil {
		log.Printf("error %s\n", err)
		return
	}
	defer outFile.Close()

	r := bufio.NewReader(inFile)
	count := 0
	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			if this.resourceManager.Count() == 0 {
				break
			} else {
				continue
			}
		}
		addr := &util.Addr{}
		json.Unmarshal([]byte(line), addr)

		if _, ok := addrs[addr.Serialize()]; !ok {
			if ok := this.resourceManager.Alloc(); !ok {
				log.Printf("waiting for resource\n")
				time.Sleep(this.pollingTime)
				continue
			}
			go func() {
				if validator.Validate(addr) {
					addrs[addr.Serialize()] = true
					outFile.WriteString(addr.Serialize() + "\n")
					log.Printf("valid %s\n", addr.Serialize())
				} else {
					log.Printf("invalid %s\n", addr.Serialize())
				}
				this.resourceManager.Free()
			}()
		} else {
			log.Printf("skip %s\n", addr.Serialize())
		}
		count++
	}
	log.Printf("generated %d agents from %d records\n", len(addrs), count)
}
