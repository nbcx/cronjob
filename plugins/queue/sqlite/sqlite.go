package main

import (
	"database/sql"
	"fmt"
	"github.com/nbcx/cronjob/ext"
	"gox/plugin/log/logger"
	"sync"
)

func init() {
	fmt.Println("sqlite plugin enable!")
}

type sqlite struct {
	db   sql.DB
	lock sync.Mutex
	log  ext.LoggerInterface
}

func NewQueue(conf ext.ConfigInterface, log ext.LoggerInterface) (ext.QueueInterface, error) {
	log.Error(conf.GetSectionString("save", "paths", "没有值啊"))

	dbPath := "./db"

	exists, _ := ext.FileExists(dbPath)
	db, err := sql.Open("sqlite3", dbPath)

	sqlite := &sqlite{
		db:  *db,
		log: log,
	}

	if !exists {
		sqlite.init()
	}

	return sqlite, err
}

func (this *sqlite) Push(key string, value string, args string) error {
	//插入数据
	stmt, err := this.db.Prepare("INSERT INTO cronjob(key, value, args, ct) values(?,?,?,?)")
	defer stmt.Close()
	this.errlog(err)

	res, err := stmt.Exec(key, value, args, "2012-12-09")
	this.errlog(err)

	_, err = res.LastInsertId()
	this.errlog(err)
	return err
}

func (this *sqlite) Pop() *ext.Task {
	this.lock.Lock()
	defer this.lock.Unlock()

	task := &ext.Task{}
	//查询数据
	rows, err := this.db.Query("SELECT * FROM cronjob limit 1")

	if !rows.Next() {
		return nil
	}
	var id string
	var ct string
	rows.Scan(&id, &task.Key, &task.Value, &task.Args, &ct)
	rows.Close()

	//删除数据
	stmt, err := this.db.Prepare("delete from cronjob where id=?")
	this.errlog(err)
	defer stmt.Close()

	res, err := stmt.Exec(id)
	this.errlog(err)

	_, err = res.RowsAffected()
	this.errlog(err)

	return task
}

func (this *sqlite) init() {
	//创建表
	sql_table := `
    CREATE TABLE IF NOT EXISTS cronjob(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        key text NULL,
        value text NULL,
        args text NULL,
        ct DATE NULL
    );
    `
	this.db.Exec(sql_table)
}

func (this *sqlite) errlog(err error) {
	if err != nil {
		logger.Error(err)
		panic(err)
	}
}
