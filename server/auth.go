// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	"context"
	"errors"
	"github.com/niklaus-code/goftp/config"
	"go.mongodb.org/mongo-driver/bson"
)

// Auth is an interface to auth your ftp user login.
type Auth interface {
	CheckPasswd(string, string) (int, error)
}

type User_datasets struct {
	Id  string
	Rpassword string
	Wpassword string
	Datapath  string
}

func CheckPasswd(name string, pwd string) (*User_datasets, int, error) {
    var dbsort = config.Dbsort

	switch {
	case dbsort == "mongodb":
		return check_mongo(name, pwd)
	default:
		return check_sql(name, pwd)
	}

}

func check_sql(user string, pwd string) (*User_datasets, int, error) {
    dbs, err := config.Db()
    if err != nil {
    	return nil, 0, err
	}

    var ftpuser User_datasets
    errdb := dbs.Where("id=?", user).First(&ftpuser)
    if errdb.Error != nil {
		return nil, 0, errdb.Error
    }
    if ftpuser.Rpassword == pwd {
    	return &ftpuser, 0, nil
	}
	if ftpuser.Wpassword == pwd {
		return &ftpuser, 1, nil
	}

	return nil, 0, errors.New("认证失败")
}

//mongo auth
func check_mongo(user string, pwd string) (*User_datasets, int, error) {
	mongoclient, err := config.Db_mongo()
	if err != nil {
		return nil, 0, err
	}
	collection := mongoclient.Database("bs_data").Collection("tb_user_ftp")

	filter := bson.D{{"username", user}}

	var u User_datasets
	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, 0, err
	}

	if u.Rpassword == pwd {
		return &u, 0, nil
	}
	if u.Wpassword == pwd {
		return &u, 1, nil
	}

	return nil, 0, errors.New("认证失败")
}