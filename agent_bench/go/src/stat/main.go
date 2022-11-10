package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var fileDir = flag.String("dir", "", "日志输出的目录，统计每个文件的tps")

// 获取文件的行数
func getFileLine(filename string) (line int) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("getFileLine error, file name:", filename, " error:", err)
		return 0
	}
	defer file.Close()
	line = 0
	reader := bufio.NewReader(file)
	for {
		_, isPrefix, err := reader.ReadLine()
		if err != nil {
			break
		}
		if !isPrefix {
			line++
		}
	}
	return
}

func statAllFiles(lastFileInfo map[string](int), lastStatTime *time.Time) {
	files, _ := ioutil.ReadDir(*fileDir)
	for _, file := range files {
		filePath := *fileDir + "/" + file.Name()
		lastStatLine := 0
		ok := false
		diffLine := 1
		if !file.IsDir() {
			line := getFileLine(filePath)
			if lastStatLine, ok = lastFileInfo[filePath]; ok {
				diffLine = line - lastStatLine
			} else {
				diffLine = line
			}
			fmt.Printf("%s line:%d,tps :%d\n", filePath, line, diffLine/int(time.Since(*lastStatTime).Seconds()+1))
			lastFileInfo[filePath] = line
		}
	}
	*lastStatTime = time.Now()
}

func main() {
	lastStatTime := time.Now()
	lastFileInfo := map[string](int){}
	flag.Parse()
	isTerminal := make(chan bool)
	quitMain := make(chan bool)
	go func() {
		var ticker = time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				if len(*fileDir) > 0 {
					statAllFiles(lastFileInfo, &lastStatTime)
				}
			case <-isTerminal:
				quitMain <- true
			}
		}
	}()
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("程序退出", s)
				isTerminal <- true
			}
		}
	}()
	<-quitMain
}
