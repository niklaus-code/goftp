// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// This is a very simple ftpd server using this library as an example
// and as something to run tests against.
package main

import (
	"flag"
	"log"

    "os"
	filedriver "github.com/niklaus-code/goftp/file-driver"
	"github.com/niklaus-code/goftp/server"
)

func init() {
        file := "./" +"message"+ ".txt"
        logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
        if err != nil {
                panic(err)
        }
        log.SetOutput(logFile) // 将文件设置为log输出的文件
        log.SetPrefix("[qSkipTool]")
        log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
        return
}

func main() {
	var (
		root = flag.String("root", "/tmp", "Root directory to serve")
		port = flag.Int("port", 21, "Port")
		host = flag.String("host", "0.0.0.0", "Host")
	)
	flag.Parse()
	if *root == "" {
		log.Fatalf("Please set a root to serve with -root")
	}

	factory := &filedriver.FileDriverFactory{
		RootPath: *root,
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		RootPath: *root,
	}

	log.Printf("请使用root用户启动项目")
	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)
	//log.Printf("Username %v, Password %v", *user, *pass)
	server := server.NewServer(opts)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
