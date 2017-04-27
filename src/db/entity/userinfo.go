package entity

import "time"

type Userinfo struct {
	Uid int  `db:"uid"`
	Username string  `db:"username"`
	Departname string  `db:"departname"`
	Created time.Time  `db:"created"`
}