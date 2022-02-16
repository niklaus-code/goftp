// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
)


type Logger interface {
	Print(sessionId string, message interface{})
	Printf(sessionId string, format string, v ...interface{})
	PrintCommand(sessionId string, command string, params string)
	PrintResponse(sessionId string, code int, message string)
}

type StdLogger struct{}

func (logger *StdLogger) Print(sessionId string, message interface{}) {
	log.Printf("%s  %s", sessionId, message)
}

func (logger *StdLogger) Printf(sessionId string, format string, v ...interface{}) {
	logger.Print(sessionId, fmt.Sprintf(format, v...))
}

func (logger *StdLogger) PrintCommand(sessionId string, command string, params string) {
	if command == "PASS" {
		log.Printf("%s > PASS ****", sessionId)
	} else {
		log.Printf("%s > %s %s", sessionId, command, params)
	}
}

func (logger *StdLogger) PrintResponse(sessionId string, code int, message string) {
	log.Printf("%s < %d %s", sessionId, code, message)
}


// Use an instance of this to log in a standard format
type Loggerfile struct{
	time string
	logfile string
}

func (logger *Loggerfile) Print(sessionId string, message interface{}) {
	logfile(fmt.Sprintf("%s %s  %s",logger.time, sessionId, message), logger.logfile)
}

func (logger *Loggerfile) Printf(sessionId string, format string, v ...interface{}) {
	logfile(logger.time + sessionId + fmt.Sprintf(format, v...), logger.logfile)
}

func (logger *Loggerfile) PrintCommand(sessionId string, command string, params string) {
	if command == "PASS" {
		logfile(fmt.Sprintf("%s %s > PASS ****", logger.time, sessionId), logger.logfile)
	} else {
		logfile(fmt.Sprintf("%s %s > %s %s",logger.time, sessionId, command, params), logger.logfile)
	}
}

func (logger *Loggerfile) PrintResponse(sessionId string, code int, message string) {
	logfile(fmt.Sprintf("%s %s < %d %s",logger.time, sessionId, code, message), logger.logfile)
}

// Silent logger, produces no output
type DiscardLogger struct{}

func (logger *DiscardLogger) Print(sessionId string, message interface{})                  {}
func (logger *DiscardLogger) Printf(sessionId string, format string, v ...interface{})     {}
func (logger *DiscardLogger) PrintCommand(sessionId string, command string, params string) {}
func (logger *DiscardLogger) PrintResponse(sessionId string, code int, message string)     {}


func logfile(data string, logpath string) {
	fileHandle, err := os.OpenFile(logpath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("open file error :", err)
		return
	}
	defer fileHandle.Close()
	// NewWriter 默认缓冲区大小是 4096
	// 需要使用自定义缓冲区的writer 使用 NewWriterSize()方法
	buf := bufio.NewWriter(fileHandle)
	// 字节写入
	buf.Write([]byte(data+"\n"))

	err = buf.Flush()
	if err != nil {
		log.Println("flush error :", err)
	}
}
