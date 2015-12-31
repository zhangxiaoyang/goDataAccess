package pipeline

import (
	"database/sql"
	"fmt"
	"github.com/zhangxiaoyang/goDataAccess/spider/common"
	"log"
	"strings"
	"sync"
)

type SqlPipeline struct {
	db        *sql.DB
	tableName string
	lock      *sync.Mutex
}

func NewSqlPipeline(db *sql.DB, tableName string) *SqlPipeline {
	return &SqlPipeline{db: db, tableName: tableName, lock: &sync.Mutex{}}
}

func (this *SqlPipeline) Pipe(items []*common.Item, merge bool) {
	this.lock.Lock()
	defer this.lock.Unlock()
	tx, err := this.db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		tx.Commit()
	}()

	for _, item := range items {
		keys := []string{}
		values := make([]interface{}, len(item.GetAll()))
		marks := []string{}
		i := 0
		for key, value := range item.GetAll() {
			keys = append(keys, key)
			values[i] = value
			marks = append(marks, "?")
			i++
		}

		stmt, err := tx.Prepare(fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES(%s)",
			this.tableName, strings.Join(keys, ","), strings.Join(marks, ","),
		))
		if err != nil {
			log.Fatal(err)
		}
		res, err := stmt.Exec(values...)
		if err != nil {
			log.Fatal(err)
		}
		_, err = res.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
	}
}
