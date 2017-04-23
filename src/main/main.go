package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/lib/pq"
	"bytes"
	"gogoing"
	"time"
)

type Video struct {
	video_id int
	video_name string
	video_map_key string
}

type appContext struct {
	db *sql.DB
}

func connectDB(driverName string, dbName string) (c *appContext, errMsg error) {
	db, err := sql.Open(driverName, dbName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err !=nil {
		return nil, err
	}
	return &appContext{db}, err
}

func (c *appContext) insert(v Video) {
	lastInsertId := 0
	err := c.db.QueryRow("INSERT INTO video(video_name, video_map_key) VALUES($1,$2) RETURNING video_id", v.video_name, v.video_map_key).Scan(&lastInsertId)
	checkErr(err)
	fmt.Println("id: ", lastInsertId)
}

func (c *appContext) update(v Video) {
	var buf bytes.Buffer
	if v.video_id < 1 {
		return
	}
	buf.WriteString("UPDATE video SET ")
	if v.video_name != "" {
		buf.WriteString("video_name=$1,")
	}
	if v.video_map_key != "" {
		buf.WriteString("video_map_key=$2,")
	}
	buf.Truncate(buf.Len()-1)
	buf.WriteString(" WHERE video_id=$3")

	stmt, err := c.db.Prepare(buf.String())
	checkErr(err)
	result, err:= stmt.Exec(v.video_name, v.video_map_key, v.video_id)
	affectNum, err:= result.RowsAffected()
	checkErr(err)
	fmt.Println("update affectNum: ", affectNum)
}

func (c *appContext) query() {
	rows, err:=c.db.Query("SELECT * FROM video")
	checkErr(err)
	defer rows.Close()
	for rows.Next() {
		v := new(Video)
		err := rows.Scan(&v.video_id, &v.video_name, &v.video_map_key)
		checkErr(err)
		fmt.Println(v.video_id, v.video_name, v.video_map_key)
	}
}

func (c *appContext) delete(video_id int) {
	stmt, err:= c.db.Prepare("DELETE FROM video WHERE video_id=$1")
	checkErr(err)
	result, err:= stmt.Exec(video_id)
	affectNum, err:= result.RowsAffected()
	checkErr(err)
	fmt.Println("delete affectNum: ", affectNum)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	/*c, err:= connectDB("postgres", "user=bilibili password=0663 dbname=nodejs sslmode=disable")
	defer c.db.Close()
	checkErr(err)
	//c.insert(Video{video_id: 0,video_name:"崩坏3rd", video_map_key:"av79849534"})

	c.query()
	c.update(Video{video_id: 2,video_name:"崩坏3rd 0.98", video_map_key:"av79849534"})*/
	/*d_list := structure.NewDoubleLinkList()
	d_list.LinkFirst(&structure.Node{Value:"123lsh"})
	d_list.LinkFirst(&structure.Node{Value:"node_first"})
	fmt.Println(d_list.Len(), d_list.GetFirst().Value, d_list.GetLast().Value)
	u1 := uuid.NewV4()*/

	acceptor := gogoing.NewAcceptor().Start("192.168.0.181:7000")
	fmt.Println(acceptor.Name())
	time.Sleep(time.Hour*1)

	//acceptor.EventDispatcher().AddHandler(gogoing.CONNECT_EVENT, event.NewEventHandler(gogoing.CONNECT_EVENT))
}