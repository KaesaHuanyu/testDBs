package mysql

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"fmt"
)

//func main() {
//	db, err := sql.Open("mysql", "root:767029384971ea6a3e1cdb0a38e3db38@tcp(192.168.2.170:60013)/testDB")
//	checkErr(err)
//	defer db.Close()
//
//	names := []string{"A", "B", "C", "D"}
//
//	tx, err := db.Begin()
//	checkErr(err)
//
//	for _, name := range names {
//		stmt, err := tx.Prepare("INSERT testTB SET name=?")
//		checkErr(err)
//
//		_, err = stmt.Exec(name)
//		checkErr(err)
//
//		time.Sleep(3 * time.Second)
//		fmt.Println("OK")
//	}
//
//	err = tx.Commit()
//	checkErr(err)
//}

//CREATE DATABASE testDB;
//CREATE TABLE testTB(name varchar(20));

func InsertData(address, username, password, database, table string) error {

	//连接到mysql
	db, err := sql.Open("mysql", username + ":" + password + "@tcp(" + address + ")/" + database)
	if err != nil {
		log.Println("InsertData: Failed to link the mysql")
		return err
	}
	if db == nil {
		err := fmt.Errorf("InsertData: sql.Open.db is nil")
		return err
	}
	defer db.Close()

	//插入数据
	for i := 1; i <= 100; i++ {
		stmt, err := db.Prepare("INSERT " + table + " SET name=?")
		if err != nil {
			log.Printf("InsertData.Prepare: Failed to instert in [ %d ] time: ", i, err)
			continue
		}

		_, err = stmt.Exec("number" + strconv.Itoa(i))
		if err != nil {
			log.Printf("InsertData.Exec: Failed to instert in [ %d ] time: ", i, err)
			continue
		}
	}
	log.Println("InsertData: Insert 100 datas completely")
	return nil
}

func FindData(address, username, password, database, table string, index int) error {

	//连接到mysql
	db, err := sql.Open("mysql", username + ":" + password + "@tcp(" + address + ")/" + database)
	if err != nil {
		log.Println("InsertData: Failed to link the mysql")
		return err
	}
	if db == nil {
		fmt.Errorf("InsertData: sql.Open.db is nil")
		return err
	}
	defer db.Close()

	//查询数据
	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		log.Println("FindData: db.Query: cannot get rows")
		return err
	}
	if rows == nil {
		err := fmt.Errorf("FindData: db.Query.rows is nil")
		return err
	}
	defer rows.Close()

	stmt, err := db.Prepare("SELECT name FROM " + table + " WHERE name = ?")
	if err != nil {
		log.Println("FindData: db.Prepare: Failed to select")
		return err
	}
	defer stmt.Close()

	var result string
	var count int
	for i := 1; i <= index * 100; i++ {
		err = stmt.QueryRow("number" + strconv.Itoa(i)).Scan(&result) // WHERE number = 13
		if err != nil {
			log.Printf("FindData: Failed to find data in [ %d ] time\n", i)
			continue
		}
		count++
	}

	log.Printf("FindData: there are [ %d ] datas\n", count)
	return nil
}
