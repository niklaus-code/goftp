// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package server

import (
	//"context"
	"errors"
	"fmt"
	"goftp/config"
	//"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

// Auth is an interface to auth your ftp user login.
type Auth interface {
	CheckPasswd(string, string) (int, error)
}

func CheckPasswd(name string, pwd string) (string,  int, error) {
    var dbsort = config.Dbsort

	switch {
	case dbsort == "mongodb":
		return "", 0, nil
		//return check_mongo(name, pwd)
	default:
		return check_sql(name, pwd)
	}
}

func mapuser() (map[string]string, error) {
	usertable := config.Fuobj()

	var usermap = make(map[string]string)
	val := reflect.Indirect(reflect.ValueOf(usertable))
	fieldnum := val.Type().NumField()
	for i := 0; i < fieldnum; i++ {
		usermap[val.Type().Field(i).Tag.Get("json")] = val.Type().Field(i).Name
	}

	return usermap, nil
}

func check_sql(user string, pwd string) (string, int, error) {
	mapu, _ := mapuser()
    dbs, err := config.Db()
    if err != nil {
    	return "", 0, err
	}

	sql := fmt.Sprintf("%s='%s'", mapu["user"], user)

	u := config.Fuobj()
    errdb := dbs.Where(sql).First(&u)
    if errdb.Error != nil {
		return "", 0, errdb.Error
    }
	val := reflect.Indirect(reflect.ValueOf(u))

    mu, _ := mapuser()

    if val.FieldByName(mu["rpasswd"]).String() == pwd {
    	return val.FieldByName(mu["datapath"]).String(), 0, nil
	}
	if val.FieldByName(mu["wpasswd"]).String()  == pwd {
		return val.FieldByName(mu["datapath"]).String(), 1, nil
	}

	return "", 0, errors.New("认证失败")
}

//mongo auth
//func check_mongo(user string, pwd string) (*config.Ftptable, int, error) {
//	var usertable config.Ftptable
//	mongoclient, err := config.Db_mongo()
//	if err != nil {
//		return nil, 0, err
//	}
//	collection := mongoclient.Database("bs_data").Collection("tb_user_ftp")
//
//	filter := bson.D{{"username", user}}
//
//	err = collection.FindOne(context.TODO(), filter).Decode(&user)
//	if err != nil {
//		return nil, 0, err
//	}
//
//	if usertable.Rpasswd == pwd {
//		return &usertable, 0, nil
//	}
//	if usertable.Wpasswd == pwd {
//		return &usertable, 1, nil
//	}
//
//	return nil, 0, errors.New("认证失败")
//}
