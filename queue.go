package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nbcx/cronjob/ext"
	"plugin"
	"sync"
)

func NewQueue(conf *Config) (ext.QueueInterface, error) {

	queue := conf.GetString("queue")

	if queue != "" {
		has, err := ext.FileExists(queue)

		if !has {
			return nil, err
		}

		//打开动态库
		p, err := plugin.Open(queue)
		if err != nil {
			panic(err)
		}

		//接口验证
		f, err := p.Lookup("NewQueue")
		if err != nil {
			panic(err)
		}
		d := f.(func(conf ext.ConfigInterface, log ext.LoggerInterface) ext.QueueInterface)(conf, logger)
		return d, nil
	}

	dbPath := "./db"

	exists, _ := ext.FileExists(dbPath)
	db, err := sql.Open("sqlite3", dbPath)

	sqlite := &sqliteQueue{
		db: *db,
	}

	if !exists {
		sqlite.init()
	}

	return sqlite, err
}

type sqliteQueue struct {
	db   sql.DB
	lock sync.Mutex
}

func (this *sqliteQueue) Push(cmd string, args string) error {
	//插入数据
	stmt, err := this.db.Prepare("INSERT INTO command(`key`, `value`, args, ct) values(?,?,?)")
	this.errlog(err)

	res, err := stmt.Exec(cmd, args, "2012-12-09")
	this.errlog(err)

	_, err = res.LastInsertId()
	this.errlog(err)
	return err
}

func (this *sqliteQueue) Pop() *ext.Task {
	this.lock.Lock()
	defer this.lock.Unlock()

	task := &ext.Task{}
	//查询数据
	rows, err := this.db.Query("SELECT * FROM command limit 1")

	if !rows.Next() {
		return nil
	}

	rows.Scan(&task.Id, &task.Cmd, &task.Args, &task.Ct)
	rows.Close()

	//删除数据
	stmt, err := this.db.Prepare("delete from command where id=?")
	this.errlog(err)

	res, err := stmt.Exec(task.Id)
	this.errlog(err)

	_, err = res.RowsAffected()
	this.errlog(err)

	return task
}

func (this *sqliteQueue) init() {
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

func (this *sqliteQueue) errlog(err error) {
	if err != nil {
		logger.Error(err)
		panic(err)
	}
}
