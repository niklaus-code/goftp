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

func Db_mongo() *mongo.Client {
	var ip = cfg.Section("mongodb").Key("ip").String()
	var port = cfg.Section("mongodb").Key("port").String()
	var user = cfg.Section("mongodb").Key("user").String()
	var passwd = cfg.Section("mongodb").Key("passwd").String()
	var database = cfg.Section("mongodb").Key("database").String()

	// Set client options
	mongodb_url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, passwd, ip, port, database)
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
	var ip = cfg.Section("postgres").Key("ip").String()
	// var port = cfg.Section("postgres").Key("port").String()
	var user = cfg.Section("postgres").Key("user").String()
	// var passwd = cfg.Section("postgres").Key("passwd").String()
	var database = cfg.Section("postgres").Key("database").String()

	conn := fmt.Sprintf("host=%s  user=%s  dbname=%s  sslmode=disable", ip, user, database)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("222222")
		fmt.Println(err)
		return nil
	}
	return db
}

func Download_rate() int {
	cfg, err := ini.Load("conf/setting.ini")
	if err != nil {
		fmt.Println(err)
		return 0
	}
	rate, err := cfg.Section("download").Key("rate").Int()
	if err != nil {
		return 0
	}
	return rate
}
