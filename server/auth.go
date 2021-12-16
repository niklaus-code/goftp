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

var usertable config.Ftp_user

func CheckPasswd(name string, pwd string) (*config.Ftp_user, int, error) {
    var dbsort = config.Dbsort

	switch {
	case dbsort == "mongodb":
		return check_mongo(name, pwd)
	default:
		return check_sql(name, pwd)
	}

}

func check_sql(user string, pwd string) (*config.Ftp_user, int, error) {
    dbs, err := config.Db()
    if err != nil {
    	return nil, 0, err
	}

    errdb := dbs.Where("id=?", user).First(&usertable)
    if errdb.Error != nil {
		return nil, 0, errdb.Error
    }
    if usertable.Rpassword == pwd {
    	return &usertable, 0, nil
	}
	if usertable.Wpassword == pwd {
		return &usertable, 1, nil
	}

	return nil, 0, errors.New("认证失败")
}

//mongo auth
func check_mongo(user string, pwd string) (*config.Ftp_user, int, error) {
	mongoclient, err := config.Db_mongo()
	if err != nil {
		return nil, 0, err
	}
	collection := mongoclient.Database("bs_data").Collection("tb_user_ftp")

	filter := bson.D{{"username", user}}

	err = collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, 0, err
	}

	if usertable.Rpassword == pwd {
		return &usertable, 0, nil
	}
	if usertable.Wpassword == pwd {
		return &usertable, 1, nil
	}

	return nil, 0, errors.New("认证失败")
}