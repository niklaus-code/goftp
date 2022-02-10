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

//ftpuser table
//tag is important , don't change or delete
type Ftptable struct {
	Uuser  string `json:"user"`
	Rprasswd string `json:"rpasswd"`
	Wwpasswd string `json:"wpasswd"`
	Ddatapath  string `json:"datapath"`
}

var cfg, _ = ini.Load("conf/setting.ini")
var Dbsort = "mysql"

var StartPort = cfg.Section("pasvport").Key("startport").String()
var RangePort = cfg.Section("pasvport").Key("rangeport").String()


func Db_mongo() (*mongo.Client, error) {
	// Set client options
	mongodb_url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", "user", "passwd", "127.0.0.1", "27017", "goftp")
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
	var db *gorm.DB
	var err error

	if Dbsort == "mysql" {
		dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local","nicloud", "nicloud", "127.0.0.1", "3306", "goftp")
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			})
		if err != nil {
			return nil, err
		}
	} else {
		dsn := fmt.Sprintf("host=%s user=postgres password='' database=gscloud_web port=5432 sslmode=disable TimeZone=Asia/Shanghai","127.0.0.1")
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
