// wrapping mysql database access
package mydb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const User = "root"
const Password = "root"

// const DB = "api.masmo"
const DB = "masmo"

type MyDB struct {
	db *sql.DB
}

func (this *MyDB) Connect() {
	var err error
	// this.db, err = sql.Open("mysql", User+":"+Password+"@tpc(http://192.168.59.103/:3306)/"+DB)
	this.db, err = sql.Open("mysql", User+"@tcp(192.168.59.103:3306)/"+DB)
	if err != nil {
		log.Fatal(err)
	}
}

func (this *MyDB) Close() {
	var err error
	err = this.db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func (this *MyDB) Query(sql string) *sql.Rows {
	rows, err := this.db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

func (this *MyDB) QueryRow(sql string) *sql.Row {
	row := this.db.QueryRow(sql)
	return row
}

func (this *MyDB) FetchAll(sql string) [][]interface{} {
	rows, err := this.db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	valuePtrs := make([]interface{}, count)

	ret := make([][]interface{}, 0)
	for rows.Next() {

		values := make([]interface{}, count)
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		for i, _ := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			values[i] = v
		}
		ret = append(ret, values)
	}

	return ret
}
