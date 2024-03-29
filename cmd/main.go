package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/SAP/go-hdb/driver"
	"github.com/gin-gonic/gin"

	"go-kyma-user-api/internal/user"
)

var dbConn *sql.DB

type DelStruct struct {
	UserName string `uri:"userName" binding:"required"`
}

func main() {

	dbConnectionError := connectToHanaDB()
	if dbConnectionError != nil {
		log.Fatalln(dbConnectionError)
		return
	} else {
		log.Println("Successfully connected to hdb!!!")
	}

	router := setupRouter()
	router.Run(":8080")

	defer dbConn.Close()
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/users", getUsers)
	r.POST("/users", createUser)
	r.DELETE("/users/:userName", deleteUser)

	return r
}

func connectToHanaDB() error {
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

	dbConn = db

	pingErr := dbConn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return err
	}
	return nil
}

func getUsers(c *gin.Context) {
	var users []user.User
	var user user.User

	log.Println("Execute Select Query!")
	rows, err := dbConn.Query("SELECT USERNAME, EMAIL, FIRSTNAME, LASTNAME, ADDRESS, MOBILE from USERAPI.USER")
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	log.Println("ROWWS!")

	for rows.Next() {
		if err := rows.Scan(&user.UserName, &user.Email, &user.FirstName,
			&user.LastName, &user.Address, &user.Mobile); err != nil {
			log.Println(err)
			return
		}

		users = append(users, user)
	}

	log.Println("USEERRS!")

	c.IndentedJSON(http.StatusOK, users)
}

func createUser(c *gin.Context) {

	var user user.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "JSON Binding Error",
				"message": "Invalid JSON!"})
		return
	}

	log.Println("Execute Create Query!")

	insertStatement, err := dbConn.Prepare("INSERT INTO USERAPI.USER (USERNAME, EMAIL, FIRSTNAME, LASTNAME, ADDRESS, MOBILE) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	defer insertStatement.Close()

	result, err := insertStatement.Exec(user.UserName, user.Email, user.FirstName, user.LastName, user.Address, user.Mobile)
	if err != nil {
		log.Fatal(err)
	}

	insertedID, err := result.LastInsertId()

	c.IndentedJSON(http.StatusOK, "User"+strconv.FormatInt(insertedID, 64)+" Created!")
}

func deleteUser(c *gin.Context) {

	var user DelStruct

	if err := c.ShouldBindUri(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   "URI Binding Error",
				"message": "Invalid URI!"})
		return
	}
	deleteStatement, err := dbConn.Prepare("DELETE FROM USERAPI.USER WHERE USERNAME = ?;")
	if err != nil {
		log.Fatal(err)
	}

	defer deleteStatement.Close()

	result, err := deleteStatement.Exec(user.UserName)
	if err != nil {
		log.Fatal(err)
	} else {
		count, err := result.RowsAffected()
		if err == nil && count > 0 {
			c.IndentedJSON(http.StatusOK, "User-> "+user.UserName+" deleted!")
		} else {
			c.IndentedJSON(http.StatusInternalServerError, "Deletion failed!")
		}
	}
}
