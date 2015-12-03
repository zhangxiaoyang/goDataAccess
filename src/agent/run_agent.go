package main

import (
	"io/ioutil"
	"log"
	"os"
	"spider/core/engine"
)

func main() {
	dbDir := "db/"
	ruleDir := "rule/"

	file, _ := os.Create(dbDir + "agent.json")
	defer file.Close()

	if fileInfos, err := ioutil.ReadDir(ruleDir); err == nil {
		for _, f := range fileInfos {
			log.Printf("Crawling %s\n", ruleDir+f.Name())
			engine.NewQuickEngine(ruleDir + f.Name()).SetOutputFile(file).Start()
			log.Printf("Finished %s\n", ruleDir+f.Name())
		}
	}

}
