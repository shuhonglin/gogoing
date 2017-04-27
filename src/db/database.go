package db

import (
	"os"
	"ini"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"fmt"
	"strings"
)

var GogoingDB *sqlx.DB

func init() {
	workDir,_ := os.Getwd()
	config := ini.SetConfig(workDir+"/src/conf/conf.ini")
	if strings.ToLower(config.GetValue("db", "db_type")) == "mysql" {
		mysqldsn := config.GetValue("db", "username")+":"+config.GetValue("db", "password")+"@tcp("+
			config.GetValue("db", "hostname")+":"+config.GetValue("db", "port")+
			")/"+config.GetValue("db", "dbname")
		ConnectMysql(mysqldsn)
	} else if strings.ToLower(config.GetValue("db", "db_type")) == "postgres" {
		pgdsn := "postgres://"+config.GetValue("db", "username")+":"+config.GetValue("db", "password")+"@"+
			config.GetValue("db", "hostname")+":"+config.GetValue("db", "port")+
			"/"+config.GetValue("db", "dbname")+"?sslmode="+config.GetValue("db", "sslmode")
		ConnectPostgre(pgdsn)
	}

}

func ConnectMysql(mysqldsn string) {
	var err error
	/*workDir,_ := os.Getwd()
	config := ini.SetConfig(workDir+"/src/conf/conf.ini")*/
	if !strings.Contains(mysqldsn, "parseTime=true") {
		mysqldsn += "?parseTime=true"
	}
	// username:password@protocol(address)/dbname?param=value
	fmt.Println(mysqldsn)
	GogoingDB, err = sqlx.Connect("mysql", mysqldsn)
	if err != nil {
		fmt.Printf("Can't connect to MySQL for:    %v\n", err)
		panic(err)
	}
}

func ConnectPostgre(pgdsn string) {
	var err error
	// "postgres://postgres:postgres@192.168.56.101/test_db?sslmode=disable"
	fmt.Println(pgdsn)
	GogoingDB, err = sqlx.Connect("postgres", pgdsn)
	if err != nil {
		fmt.Printf("Can't connect to PostgreSQL for:    %v\n", err)
		panic(err)
	}
}

func CloseAll() {
	if GogoingDB!=nil {
		GogoingDB.Close()
	}
}

