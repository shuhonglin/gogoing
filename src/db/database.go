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

var Mysql *sqlx.DB
var Postgre *sqlx.DB

func init() {
	workDir,_ := os.Getwd()
	config := ini.SetConfig(workDir+"/src/conf/conf.ini")
	if config.GetSection("mysql")!=nil {
		mysqldsn := config.GetValue("mysql", "username")+":"+config.GetValue("mysql", "password")+"@tcp("+
			config.GetValue("mysql", "hostname")+":"+config.GetValue("mysql", "port")+
			")/"+config.GetValue("mysql", "dbname")
		ConnectMysql(mysqldsn)
	}
	if config.GetSection("postgres")!=nil {
		pgdsn := "postgres://"+config.GetValue("postgres", "username")+":"+config.GetValue("postgres", "password")+"@"+
			config.GetValue("postgres", "hostname")+":"+config.GetValue("postgres", "port")+
			"/"+config.GetValue("postgres", "dbname")+"?sslmode="+config.GetValue("postgres", "sslmode")
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
	Mysql, err = sqlx.Connect("mysql", mysqldsn)
	if err != nil {
		fmt.Printf("Can't connect to MySQL for:    %v\n", err)
		panic(err)
	}
}

func ConnectPostgre(pgdsn string) {
	var err error
	// "postgres://postgres:postgres@192.168.56.101/test_db?sslmode=disable"
	fmt.Println(pgdsn)
	Postgre, err = sqlx.Connect("postgres", pgdsn)
	if err != nil {
		fmt.Printf("Can't connect to PostgreSQL for:    %v\n", err)
		panic(err)
	}
}

func CloseAll() {
	if Mysql!=nil {
		Mysql.Close()
	}
	if Postgre!=nil {
		Postgre.Close()
	}
}

