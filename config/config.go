package config

import (
	"fmt"

	"github.com/go-ini/ini"
	_ "github.com/lib/pq"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/go-sql-driver/mysql"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

var cfg, _ = ini.Load("conf/setting.ini")

var StartPort = cfg.Section("pasvport").Key("startport").String()
var RangePort = cfg.Section("pasvport").Key("rangeport").String()
var Dbsort = cfg.Section("db").Key("dbsort").String()

func Ftpuser() map[string]string {
	var table = cfg.Section("ftpuser").Key("table").String()
	var user = cfg.Section("ftpuser").Key("user").String()
	var rpasswd = cfg.Section("ftpuser").Key("rpasswd").String()
	var wpasswd = cfg.Section("ftpuser").Key("wpasswd").String()
	var datadir = cfg.Section("ftpuser").Key("datapath").String()

	config := make(map[string]string)
	config["table"] = table
	config["datapath"] = datadir
	config["user"] = user
	config["rpasswd"] = rpasswd
	config["wpasswd"] = wpasswd
	return config
}

func export(Dbsort string) map[string]string {
	var user = cfg.Section(Dbsort).Key("user").String()
	var passwd = cfg.Section(Dbsort).Key("passwd").String()
	var ip = cfg.Section(Dbsort).Key("ip").String()
	var port = cfg.Section(Dbsort).Key("port").String()
	var database = cfg.Section(Dbsort).Key("database").String()

	config := make(map[string]string)
	config["user"] = user
	config["passwd"] = passwd
	config["ip"] = ip
	config["port"] = port
	config["database"] = database
	return config
}


func Db_mongo() (*mongo.Client, error) {
	var config = export("mongodb")

	// Set client options
	mongodb_url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", config["user"], config["passwd"], config["ip"], config["port"], config["database"])
	clientOptions := options.Client().ApplyURI(mongodb_url)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		return nil, err
	}
	return client, nil
}

func Db() (*gorm.DB,error) {
	db, errDb:=gorm.Open("mysql","nicloud:nicloud@(127.0.0.1:3306)/ftp?parseTime=true")
	if errDb != nil {
		fmt.Println(errDb.Error())
		return nil, errDb
	}

	sqlDB := db.DB()
	sqlDB.SetMaxIdleConns(100) //空闲连接数
	sqlDB.SetMaxOpenConns(1000)//最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * 360)

	return db, nil
}

func Download_rate() int {
	cfg, err := ini.Load("conf/setting.ini")
	if err != nil {
		return 0
	}
	rate, err := cfg.Section("download").Key("rate").Int()
	if err != nil {
		return 0
	}
	return rate
}
