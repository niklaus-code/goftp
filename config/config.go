package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/schema"

	"github.com/go-ini/ini"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

//ftpuser table
//field can be change for match databse,but tag con't change or delete , tag is important
type Ftpuser struct {
	User  string `json:"user"`
	Rpasswd string `json:"rpasswd"`
	Wpasswd string `json:"wpasswd"`
	Datapath  string `json:"datapath"`
}

func Fuobj ()  *Ftpuser {
	fu := Ftpuser{}
	return &fu
}

var cfg, _ = ini.Load("conf/setting.ini")
var Dbsort = cfg.Section("dbsort").Key("db").String()

var StartPort = cfg.Section("pasvport").Key("startport").String()
var RangePort = cfg.Section("pasvport").Key("rangeport").String()
var dbname = cfg.Section("database").Key("dbname").String()
var host = cfg.Section("database").Key("dbhost").String()
var port = cfg.Section("database").Key("port").String()
var user = cfg.Section("database").Key("user").String()
var passwd = cfg.Section("database").Key("passwd").String()


//func Db_mongo() (*mongo.Client, error) {
//	// Set client options
//	mongodb_url := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", "user", "passwd", "127.0.0.1", "27017", "goftp")
//	clientOptions := options.Client().ApplyURI(mongodb_url)
//
//	// Connect to MongoDB
//	client, err := mongo.Connect(context.TODO(), clientOptions)
//
//	if err != nil {
//		return nil, err
//	}
//
//	// Check the connection
//	err = client.Ping(context.TODO(), nil)
//
//	if err != nil {
//		return nil, err
//	}
//	return client, nil
//}

func Db() (*gorm.DB,error) {
	var db *gorm.DB
	var err error

	if Dbsort == "mysql" {
		dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",user, passwd, host, port, dbname)
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
