package db

import (
	"db/entity"
	"fmt"
	"reflect"
	"strings"
	"time"
	"strconv"
)

type ProxyUserinfo struct {
	userinfo *entity.Userinfo
	timer *time.Timer
}

func (self *ProxyUserinfo) LazyLoad(id int64) *entity.Userinfo {
	if self.userinfo == nil {
		self.userinfo = new(entity.Userinfo)
		err := GogoingDB.Get(self.userinfo, "SELECT * FROM go_schemas.userinfo WHERE uid=$1", id)
		if err!=nil {
			fmt.Println("Load userinfo error:", err)
		}
	}
	return self.userinfo
}

func (self *ProxyUserinfo) SetTimer(timer *time.Timer) {
	self.timer = timer
}

func (self *ProxyUserinfo) Timer() *time.Timer {
	return self.timer
}

func (self *ProxyUserinfo) Save() {
	self.userinfo = new(entity.Userinfo)
	self.userinfo.Uid = 1
	self.userinfo.Departname = "depart"
	self.userinfo.Username ="lsh"
	self.userinfo.Created = time.Now()
	if self.userinfo != nil {
		updateSql:="UPDATE "
		refValue := reflect.ValueOf(*self.userinfo)
		refType := refValue.Type()
		updateSql+=strings.ToLower(refType.Name()[strings.LastIndex(refType.Name(),".")+1:])+" SET "
		for i:=0;i<refValue.NumField();i++  {
			f := refValue.Field(i)
			fmt.Printf("%d. %s %s = %v \n", i, refType.Field(i).Name, f.Type().String(), f.Interface())
			if strings.ToLower(refType.Field(i).Name)!="uid" {
				updateSql+=strings.ToLower(refType.Field(i).Name)+"="
				switch f.Type().String() {
				case "string","time.Time":
					updateSql += "'"+fmt.Sprintf("%v", f.Interface())+"'"
				default:
					updateSql += fmt.Sprintf("%v", f.Interface())
				}
				if i!=refValue.NumField()-1 {
					updateSql+=","
				}
			}
		}
		updateSql+=" WHERE uid="+strconv.Itoa(self.userinfo.Uid)
		fmt.Println(updateSql)
	}
}