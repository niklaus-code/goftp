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
		return checkpass(name, pwd)
	}
}

func mapuser() (map[string]string, error) {
	usertable := config.Ftpuserobj()

	var usermap = make(map[string]string)
	val := reflect.Indirect(reflect.ValueOf(usertable))
	fieldnum := val.Type().NumField()
	for i := 0; i < fieldnum; i++ {
		usermap[val.Type().Field(i).Tag.Get("json")] = val.Type().Field(i).Name
	}

	return usermap, nil
}

func checkpass(user string, pwd string) (string, int, error) {
	mapu, _ := mapuser()
    dbs, err := config.Db()
    if err != nil {
    	return "", 0, err
	}

	var u interface{}
	if string(pwd[0]) == "c" {
		u = config.Ftpvdiruserobj()
	} else {
		u = config.Ftpuserobj()
	}

	sql := fmt.Sprintf("%s='%s'", mapu["user"], user)
	fmt.Println(sql)
    errdb := dbs.Where(sql).First(&u)
    if errdb.Error != nil {
    	fmt.Println(213123)
		return "", 0, errdb.Error
    }
	val := reflect.Indirect(reflect.ValueOf(u))
	fmt.Println(val.FieldByName(mapu["rpasswd"]).String())
	fmt.Println(pwd)

    if val.FieldByName(mapu["rpasswd"]).String() == pwd {
    	return val.FieldByName(mapu["datapath"]).String(), 0, nil
	}
	if val.FieldByName(mapu["wpasswd"]).String()  == pwd {
		return val.FieldByName(mapu["datapath"]).String(), 1, nil
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
