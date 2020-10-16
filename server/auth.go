// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"crypto/subtle"
	"fmt"

	"context"

	"github.com/niklaus-code/goftp-mongo/config"
	"go.mongodb.org/mongo-driver/bson"
)

// Auth is an interface to auth your ftp user login.
type Auth interface {
	CheckPasswd(string, string) (int, error)
}

type Ftpuser struct {
	Username  string
	Rpassword string
	Wpassword string
	Datapath  string
}

func CheckPasswd(dbsort string, name string, pass string) Ftpuser {
	var ftpuser Ftpuser
	switch {
	case dbsort == "mongo":
		return check_mongo(name, pass)
	case dbsort == "mysql":
		return check_sql(name, pass)
	case dbsort == "postgres":
		return check_sql(name, pass)
	default:
		return ftpuser
	}

}

func check_sql(name string, pass string) Ftpuser {
	c := config.Db()

	var ftpuser Ftpuser
	err := c.QueryRow("select username, rpassword, wpassword, datapath from goftp where username = $1", name).Scan(&ftpuser.Username, &ftpuser.Rpassword, &ftpuser.Wpassword, &ftpuser.Datapath)

	if err != nil {
		fmt.Println("--------------------")
		fmt.Println(err)
		return ftpuser
	}

	return ftpuser
}

//mongo auth
func check_mongo(name string, pass string) Ftpuser {
	mongoclient := config.Db_mongo()
	collection := mongoclient.Database("bs_data").Collection("tb_user_ftp")

	rpasswd_filter := bson.D{{"username", name}, {"rpassword", pass}}
	wpasswd_filter := bson.D{{"username", name}, {"wpassword", pass}}

	var user Ftpuser
	err := collection.FindOne(context.TODO(), rpasswd_filter).Decode(&user)
	if err == nil {
		return user
	}

	err = collection.FindOne(context.TODO(), wpasswd_filter).Decode(&user)
	if err == nil {
		return user
	}

	return user
}

// CheckPasswd will check user's password
// func (a *SimpleAuth) CheckPasswd(name, pass string) (int, error) {
// func CheckPasswd(name, pass string) (Ftpuser, error) {
// 	return check(name, pass), nil
// 	// return constantTimeEquals(name, a.Name) && constantTimeEquals(pass, a.Password), nil
// }

func constantTimeEquals(a, b string) bool {
	return len(a) == len(b) && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
