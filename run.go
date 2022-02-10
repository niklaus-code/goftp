// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// This is a very simple ftpd server using this library as an example
// and as something to run tests against.
package main

import (
	"flag"
	filedriver "github.com/niklaus-code/goftp/file-driver"
	"github.com/niklaus-code/goftp/server"
	"log"
)



func main() {
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
