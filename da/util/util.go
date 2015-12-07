package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func LoadUrlsFromFile(filePath string) []string {
	var urls = []string{}
	file, _ := os.Open(filePath)
	defer file.Close()

	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		urls = append(urls, strings.TrimSpace(line))
	}
	return urls
}
