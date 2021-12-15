package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/schema"

	"github.com/go-ini/ini"
	"gorm.io/driver/mysql"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"gorm.io/gorm"
)

var cfg, _ = ini.Load("conf/setting.ini")

var StartPort = cfg.Section("pasvport").Key("startport").String()
var RangePort = cfg.Section("pasvport").Key("rangeport").String()
var Dbsort = cfg.Section("db").Key("dbname").String()

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
	dbsort := Dbsort
	dbinfo := export(dbsort)

	var db *gorm.DB
	var err error

	if dbsort == "mysql" {
		dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local", dbinfo["user"], dbinfo["passwd"], dbinfo["ip"], dbinfo["port"], dbinfo["database"])
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			})
		if err != nil {
			return nil, err
		}
	} else {
		dsn := fmt.Sprintf("host=%s user=postgres password='' database=gscloud_web port=5432 sslmode=disable TimeZone=Asia/Shanghai", dbinfo["ip"])
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			return nil, err
		}
	}

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
