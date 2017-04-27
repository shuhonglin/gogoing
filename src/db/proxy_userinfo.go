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
}

func (self *ProxyUserinfo) LazyLoad(uid int64) {
	if self.userinfo == nil {
		self.userinfo = new(entity.Userinfo)
		err := GogoingDB.Select(self.userinfo, "SELECT * FROM userinfo WHERE uid=$1", uid)
		if err!=nil {
			fmt.Println("Load userinfo error:", err)
		}
	}
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