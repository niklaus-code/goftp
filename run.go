// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// This is a very simple ftpd server using this library as an example
// and as something to run tests against.
package main

import (
	"flag"
	filedriver "goftp/file-driver"
	"goftp/server"
	"log"
	"github.com/go-ini/ini"
	"strconv"
)

func main() {
	var cfg, _ = ini.Load("conf/setting.ini")
	var logpath = cfg.Section("log").Key("logpath").String()
	var debug = cfg.Section("debug").Key("debug").String()
	DEBUG, err := strconv.ParseBool(debug)
	if err != nil {
		DEBUG = true
	}

	var (
		port = flag.Int("port", 21, "Port")
		host = flag.String("host", "0.0.0.0", "Host")
	)

	flag.Parse()

	factory := &filedriver.FileDriverFactory{
		Perm:     server.NewSimplePerm("user", "group"),
	}

	opts := &server.ServerOpts{
		Factory:  factory,
		Port:     *port,
		Hostname: *host,
		Logpath: logpath,
		Debug: DEBUG,
	}

	log.Printf("请使用root用户启动项目")
	log.Printf("Starting ftp server on %v:%v", opts.Hostname, opts.Port)

	server := server.NewServer(opts)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
