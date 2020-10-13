package hack

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func Init(filePath string){
	host := os.Getenv("PG_HOST")
	connectionParams := os.Getenv("CONNECTION_STRING")
	connectionString := fmt.Sprintf(`host=%v %v`, host, connectionParams)

	path, err := filepath.Abs(filePath)
	if err != nil {
		log.Fatalf("cannot get path")
	}

	query, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot get path")
	}

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("can't connect to db")
	}

	_, err = db.Exec(string(query))
	if err != nil {
		log.Fatalf("can't execute query")
	}
	err = db.Close()
	if err != nil {
		log.Fatalf("can't close db")
	}
}

func TearDown(){
	host := os.Getenv("PG_HOST")
	connectionParams := os.Getenv("CONNECTION_STRING")
	connectionString := fmt.Sprintf(`host=%v %v`, host, connectionParams)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("can't connect to db")
	}

	_, err = db.Exec(`TRUNCATE TABLE "user" CASCADE; TRUNCATE TABLE "user_role" CASCADE; TRUNCATE TABLE "role" CASCADE;`)
	if err != nil {
		log.Fatalf("can't execute query")
	}
	err = db.Close()
	if err != nil {
		log.Fatalf("can't close db")
	}
}