package main

import (
	"os"
	"log"
	"database/sql"

	_ "github.com/SAP/go-hdb/driver"
)

func main() {

	dbConnectionError := connectToHanaDB()
	if dbConnectionError != nil {
		log.Fatalln(dbConnectionError)
		return
	}
 
}


func connectToHanaDB() (error){
	user := os.Getenv("HDB_USER")
	password := os.Getenv("HDB_PASSWORD")
	host := os.Getenv("HDB_HOST")
	port := os.Getenv("HDB_PORT")
	pemFile := "DigiCertGlobalRootCA.pem"
	connectionString := "hdb://" + user + ":" + password + "@" + host + ":" + port + "?TLSServerName=" + host + "&TLSRootCAFile=" + pemFile

	db, err := sql.Open("hdb", connectionString)
	if err != nil {
		log.Fatalln(err)
		return err
	}

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
		return err
    }
    //fmt.Println("Connected!")
	return nil
}
