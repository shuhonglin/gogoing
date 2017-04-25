package db

import (
	"os"
	"ini"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"fmt"
)

var mysql *sqlx.DB
var postgre *sqlx.DB

func init() {

}

func ConnectMysql() {
	var err error
	workDir,_ := os.Getwd()
	config := ini.SetConfig(workDir+"/src/conf/conf.ini")
	mysql, err = sqlx.Connect("mysql", "")
	if err != nil {
		fmt.Printf("Can't connect to MySQL for:    %v\n", err)
	}
}

