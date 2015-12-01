package main

import (
	"io/ioutil"
	"log"
	"os"
	"spider/core/engine"
)

func main() {
	file, _ := os.Create("agent.json")
	defer file.Close()

	baseDir := "rule/"
	if fileInfos, err := ioutil.ReadDir(baseDir); err == nil {
		for _, f := range fileInfos {
			log.Printf("Crawling %s\n", baseDir+f.Name())
			engine.NewQuickEngine(baseDir + f.Name()).SetOutputFile(file).Start()
			log.Printf("Finished %s\n", baseDir+f.Name())
		}
	}

}
