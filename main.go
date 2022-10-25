package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/cihub/seelog"
)

// log_format  '$log_time $remote_addr - $remote_user [$time_local] '
//     	       '"$request" $status $body_bytes_sent '
//             '"$http_referer" "$http_user_agent" '
//             '"$http_x_forwarded_for" $log_count $request_id '
//             '$geoip_country_name $geoip_country_code '
//             '$geoip_region_name $geoip_city ';
var nginxLog = `%s %s - - [21/Nov/2017:08:45:45 +0000] "POST /ngx_pagespeed_beacon?url=https%%3A%%2F%%2Fwww.example.com%%2Fads%%2Ffresh-oranges-1509260795 HTTP/2.0" 204 0 "https://www.example.com/ads/fresh-oranges-1509260795" "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:47.0) Gecko/20100101 Firefox/47.0" "-" %d %s Uganda UG Kampala Kampala
`

var javaStackLogCount = 3
var javaStackLog = `[%s] %d "BLOCKED_TEST pool-1-thread-1" prio=6 tid=0x0000000006904800 nid=0x28f4 runnable [0x000000000785f000]
java.lang.Thread.State: RUNNABLE
			 at java.io.FileOutputStream.writeBytes(Native Method)
			 at java.io.FileOutputStream.write(FileOutputStream.java:282)
			 at java.io.BufferedOutputStream.flushBuffer(BufferedOutputStream.java:65)
			 at java.io.BufferedOutputStream.flush(BufferedOutputStream.java:123)
			 - locked <0x0000000780a31778> (a java.io.BufferedOutputStream)
			 at java.io.PrintStream.write(PrintStream.java:432)
			 - locked <0x0000000780a04118> (a java.io.PrintStream)
			 at sun.nio.cs.StreamEncoder.writeBytes(StreamEncoder.java:202)
			 at sun.nio.cs.StreamEncoder.implFlushBuffer(StreamEncoder.java:272)
			 at sun.nio.cs.StreamEncoder.flushBuffer(StreamEncoder.java:85)
			 - locked <0x0000000780a040c0> (a java.io.OutputStreamWriter)
			 at java.io.OutputStreamWriter.flushBuffer(OutputStreamWriter.java:168)
			 at java.io.PrintStream.newLine(PrintStream.java:496)
			 - locked <0x0000000780a04118> (a java.io.PrintStream)
			 at java.io.PrintStream.println(PrintStream.java:687)
			 - locked <0x0000000780a04118> (a java.io.PrintStream)
			 at com.nbp.theplatform.threaddump.ThreadBlockedState.monitorLock(ThreadBlockedState.java:44)
			 - locked <0x0000000780a000b0> (a com.nbp.theplatform.threaddump.ThreadBlockedState)
			 at com.nbp.theplatform.threaddump.ThreadBlockedState$1.run(ThreadBlockedState.java:7)
			 at java.util.concurrent.ThreadPoolExecutor$Worker.runTask(ThreadPoolExecutor.java:886)
			 at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:908)
			 at java.lang.Thread.run(Thread.java:662)
Locked ownable synchronizers:
			 - <0x0000000780a31758> (a java.util.concurrent.locks.ReentrantLock$NonfairSync)

[%s] %d "BLOCKED_TEST pool-1-thread-2" prio=6 tid=0x0000000007673800 nid=0x260c waiting for monitor entry [0x0000000008abf000]
java.lang.Thread.State: BLOCKED (on object monitor)
			 at com.nbp.theplatform.threaddump.ThreadBlockedState.monitorLock(ThreadBlockedState.java:43)
			 - waiting to lock <0x0000000780a000b0> (a com.nbp.theplatform.threaddump.ThreadBlockedState)
			 at com.nbp.theplatform.threaddump.ThreadBlockedState\$2.run(ThreadBlockedState.java:26)
			 at java.util.concurrent.ThreadPoolExecutor$Worker.runTask(ThreadPoolExecutor.java:886)
			 at java.util.concurrent.ThreadPoolExecutor\$Worker.run(ThreadPoolExecutor.java:908)
			 at java.lang.Thread.run(Thread.java:662)
Locked ownable synchronizers:
			 - <0x0000000780b0c6a0> (a java.util.concurrent.locks.ReentrantLock$NonfairSync)

[%s] %d "BLOCKED_TEST pool-1-thread-3" prio=6 tid=0x00000000074f5800 nid=0x1994 waiting for monitor entry [0x0000000008bbf000]
java.lang.Thread.State: BLOCKED (on object monitor)
			 at com.nbp.theplatform.threaddump.ThreadBlockedState.monitorLock(ThreadBlockedState.java:42)
			 - waiting to lock <0x0000000780a000b0> (a com.nbp.theplatform.threaddump.ThreadBlockedState)
			 at com.nbp.theplatform.threaddump.ThreadBlockedState\$3.run(ThreadBlockedState.java:34)
			 at java.util.concurrent.ThreadPoolExecutor$Worker.runTask(ThreadPoolExecutor.java:886
			 at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:908)
			 at java.lang.Thread.run(Thread.java:662)
Locked ownable synchronizers:
			 - <0x0000000780b0e1b8> (a java.util.concurrent.locks.ReentrantLock$NonfairSync)
			 
`

var jsonLog string

var defaultConfig = `
<seelog type="asynctimer" asyncinterval="5000" minlevel="info" >
 <outputs formatid="common">
	 <rollingfile type="size" filename="%s" maxsize="%d" maxrolls="%d"/>
 </outputs>
 <formats>
	 <format id="common" format="%%Msg" />
 </formats>
</seelog>
`

var logger = seelog.Disabled

var stdoutFlag = flag.Bool("stdout", true, "output to stdout")
var stderrFlag = flag.Bool("stderr", false, "output to stderr")
var filePath = flag.String("path", "", "output to file path")
var perLogFileSize = flag.Int("log-file-size", 20*1024*1024, "max log size")
var maxLogFileCount = flag.Int("log-file-count", 10, "max rotated files")
var logsPerSec = flag.Int("logs-per-sec", 1, "logs per second upper limit")
var logType = flag.String("log-type", "java", "nginx java random json")
var logErrType = flag.String("log-err-type", "random", "nginx java random json")
var totalCount = flag.Int("total-count", 100, "total log count, set -1 for infinity")
var itemLen = flag.Int("item-length", 100, "value length in json log and nginx log")
var keyCount = flag.Int("key-count", 10, "key count in json log")

var nowCount = 0
var ip = "127.0.0.1"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func getMachineIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

func mockJsonLog() string {
	kv := make(map[string]string)
	for i := 0; i < *keyCount; i++ {
		kv[RandString(10)] = RandString(*itemLen)
	}
	kv["count"] = "%d"
	kv["log_time"] = "%s"
	val, _ := json.Marshal(kv)
	val = append(val, '\n')
	return string(val)
}

func mockOneLog(timeStr, logType string) string {
	nowCount++
	switch logType {
	case "nginx":
		return fmt.Sprintf(nginxLog, timeStr, ip, nowCount, RandString(*itemLen))
	case "java":
		return fmt.Sprintf(javaStackLog, timeStr, nowCount, timeStr, nowCount, timeStr, nowCount)
	case "json":
		return fmt.Sprintf(jsonLog, nowCount, timeStr)
	}
	return fmt.Sprintf("%s %d %s\n", timeStr, nowCount, RandString(*itemLen))
}

func dumpOneLog(timeStr string) {
	if len(*filePath) > 0 {
		logger.Info(mockOneLog(timeStr, *logType))
	}
	if *stdoutFlag {
		os.Stdout.WriteString(mockOneLog(timeStr, *logType))
	}
	if *stderrFlag {
		os.Stderr.WriteString(mockOneLog(timeStr, *logErrType))
	}
}

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	if len(*filePath) > 0 {
		log.Println("use file output, path : ", *filePath)
		logConfig := fmt.Sprintf(defaultConfig, *filePath, *perLogFileSize, *maxLogFileCount)
		fmt.Println("log config, : " + logConfig)
		var err error
		logger, err = seelog.LoggerFromConfigAsString(logConfig)
		if err != nil {
			panic(err)
		}
	}
	ip, _ = getMachineIP()
	if *logType == "json" {
		jsonLog = mockJsonLog()
	}
	i := 0
	for i < *totalCount || *totalCount < 0 {
		startTime := time.Now()
		timeStr := startTime.Format(time.RFC3339Nano)
		for j := 0; j < *logsPerSec && (i < *totalCount || *totalCount < 0); j++ {
			dumpOneLog(timeStr)
			i++
		}
		endTime := time.Now()
		sleep_time := time.Second - endTime.Sub(startTime)
		if sleep_time > 0 {
			time.Sleep(sleep_time)
		}
	}
	logger.Flush()
}
