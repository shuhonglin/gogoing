package db

import (
	"fmt"
)

type User struct {
	UserId int64 //`db:"userId`
	UserName string //`db:"userName"`
	LastLoginTime int32 //`db:"lastLoginTime"`
}

func (u User) String() string {
	return "User{" + "UserId=" + string(u.UserId) + ", UserName=" + u.UserName + ", LastLoginTime='" + string(u.LastLoginTime) +"}"
}

type UserProxy struct {

}

func (self *UserProxy) QueryUser() []User  {
	userList := []User{}
	GogoingDB.Select(&userList, "SELECT * FROM tb_user ORDER BY userId DESC")

	rows, _ := GogoingDB.Queryx("SELECT * FROM tb_user ORDER BY userId DESC")
	for rows.Next() {
		/*err := rows.StructScan(&u)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%#v\n", u)*/
		fmt.Println(rows.Columns())
	}

	return userList
}