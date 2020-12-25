package config

import (
	"database/sql"
	"fmt"

	"github.com/go-ini/ini"
	_ "github.com/lib/pq"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"log"
)

var cfg, _ = ini.Load("conf/setting.ini")

var StartPort = cfg.Section("pasvport").Key("startport").String()
var RangePort = cfg.Section("pasvport").Key("rangeport").String()
var Dbname = cfg.Section("db").Key("dbname").String()

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

func export(dbname string) map[string]string {
	var user = cfg.Section(dbname).Key("user").String()
	var passwd = cfg.Section(dbname).Key("passwd").String()
	var ip = cfg.Section(dbname).Key("ip").String()
	var port = cfg.Section(dbname).Key("port").String()
	var database = cfg.Section(dbname).Key("database").String()

	config := make(map[string]string)
	config["user"] = user
	config["passwd"] = passwd
	config["ip"] = ip
	config["port"] = port
	config["database"] = database
	return config
}


func Db_mongo() *mongo.Client {
	var config = export("mongodb")

	// Set client options
	mongodb_url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", config["user"], config["passwd"], config["ip"], config["port"], config["database"])
	fmt.Println(mongodb_url)
	clientOptions := options.Client().ApplyURI(mongodb_url)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}
	return client
}

func Db() *sql.DB {
	var config = export("postgres")

	conn := fmt.Sprintf("host=%s  user=%s  dbname=%s  sslmode=disable", config["ip"], config["user"], config["database"])
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil
	}
	return db
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
