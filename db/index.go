package db

import (
	"app/library/log"
	"app/library/types/times"
	"app/models"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"strings"
	"time"
)

type Table string

//refresh db struct to sql file
func Init() {
	log.StdOut("init", "db")
	var result = []string{"-- " + times.FormatDatetime(time.Now()) + "\n\n"}
	var db = models.DB()
	rows, err := db.Raw("SHOW TABLES;").Rows()
	if err != nil {
		log.StdFatal("show tables", err)
	}
	for rows.Next() {
		var i Table
		rows.Scan(&i)
		fmt.Println(i, err)
		result = append(result, describeTable(db, i))
	}
	log.StdOut(
		"Init",
		"Create SQL",
		ioutil.WriteFile("local/create.sql", []byte(strings.Join(result, "\n\n")), 0665),
	)
}

func describeTable(db *gorm.DB, table Table) (result string) {
	sqlStr := fmt.Sprintf("SHOW CREATE TABLE %s;", table)
	//fmt.Println(sqlStr)
	rows, err := db.Raw(sqlStr).Rows()
	if err != nil {
		log.StdFatal("describeTable", "query", err)
	}
	for rows.Next() {
		var param1 []byte
		var param2 []byte
		err := rows.Scan(&param1, &param2)
		if err == nil {
			return strings.ReplaceAll(string(param2)+";", "CREATE TABLE", "CREATE TABLE IF NOT EXISTS")
		} else {
			fmt.Println(err)
		}
		break
	}
	return
}
