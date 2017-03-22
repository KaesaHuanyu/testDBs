package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

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

func initMySQL(address, username, password, database string) error {

	//连接到mysql
	db, err := sql.Open("mysql", username + ":" + password + "@tcp(" + address + ")/" + database)
	if err != nil {
		return err
	}
	defer db.Close()

	for i := 1; i < 1000; i++ {
		stmt, err := db.Prepare("INSERT test SET name=?")
		checkErr(err)

		_, err = stmt.Exec("number" + strconv.Itoa(i))
		checkErr(err)

	}



}
