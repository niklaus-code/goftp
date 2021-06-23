// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"crypto/subtle"
	"database/sql"
	"fmt"

	"context"

	"github.com/niklaus-code/goftp/config"
	"go.mongodb.org/mongo-driver/bson"
)

// Auth is an interface to auth your ftp user login.
type Auth interface {
	CheckPasswd(string, string) (int, error)
}

type Ftpuser struct {
	Id  string
	Rpasswd sql.NullString
	Wpasswd sql.NullString
	Datapath  string
}

func CheckPasswd(name string, pass string) (*Ftpuser, error) {
    var dbsort = config.Dbsort

	switch {
	case dbsort == "mongodb":
		return check_mongo(name, pass)
	default:
		return check_sql(name, pass)
	}

}

func check_sql(name string, pass string) (*Ftpuser, error) {
    c := config.Db()
    t := config.Ftpuser()

    title := fmt.Sprintf("select %s, %s, %s, %s from %s where %s = '%s'", t["user"], t["rpasswd"], t["wpasswd"], t["datapath"], t["table"], t["user"], name)
    var ftpuser Ftpuser
    err := c.QueryRow(title).Scan(&ftpuser.Id, &ftpuser.Rpasswd, &ftpuser.Wpasswd, &ftpuser.Datapath)

    c.Close()
    if err != nil {
		return nil, err
    }
    return &ftpuser, nil
}

//mongo auth
func check_mongo(name string, pass string) (*Ftpuser, error) {
	mongoclient := config.Db_mongo()
	collection := mongoclient.Database("bs_data").Collection("tb_user_ftp")

	filter := bson.D{{"username", name}}

	var user Ftpuser
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
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
