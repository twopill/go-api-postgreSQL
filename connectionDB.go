package main

import (
	mydata "./app"
	"database/sql"
	_ "encoding/json"
	"fmt"
	_ "github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/rs/cors"
	cors "github.com/rs/cors/wrapper/gin"
	"log"
	"net/http"
	_ "os"
)
//commit test
type personTable struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
}

type accessdata struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/auth/useraccess", func(c *gin.Context) {

		var user accessdata
		c.BindJSON(&user)

		u := accessdata{
			User:     user.User,
			Password: user.Password}
		if u.User == "" && u.Password == "" { // only for test
			fmt.Println("entrato")
			c.JSON(200, gin.H{"ACCESS_TOKEN": mydata.GenerateToken()})
		}

	})

	router.GET("/getData", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"payload": selectAllFromPersonTable(),
		})
	})

	router.GET("/dataGet/:name", func(c *gin.Context) {
		name := c.Param("name")
		if len(name) > 0 {
			c.String(http.StatusOK, "Hello %s", name)
		} else {
			c.String(http.StatusNotFound, "Error 	")
		}
	})

	router.POST("/postData", func(c *gin.Context) {

		var user personTable

		c.BindJSON(&user)

		e := personTable{
			ID:      user.ID,
			Name:    user.Name,
			Age:     user.Age,
			Address: user.Address}

		err := insert(e)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("creato")
		c.JSON(200, gin.H{"id": user.ID, "name": user.Name, "age": user.Age, "address": user.Address, "status": http.StatusOK})

	})

	router.Run(":8080")

}

func selectByName(name string) { //TODO: da vedere
	var person2 []personTable
	q := `SELECT * FROM person_table WHERE name = ?`
	db := getConnection()
	defer db.Close()

	rows, err := db.Query(q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		e := personTable{}
		err = rows.Scan(
			&e.ID,
			&e.Name,
			&e.Age,
			&e.Address)
		if err != nil {
			return
		}
		person2 = append(person2, e)
	}
	fmt.Println(person2)

}

func selectAllFromPersonTable() (person []personTable) {
	q := `SELECT * FROM person_table`
	db := getConnection()
	defer db.Close()
	rows, err := db.Query(q)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		e := personTable{}
		err = rows.Scan(
			&e.ID,
			&e.Name,
			&e.Age,
			&e.Address)
		if err != nil {
			return
		}
		person = append(person, e)
	}
	fmt.Println(person)
	return person
}

func insert(p personTable) error {

	q := `INSERT INTO
	person_table(id, name, age, address)
	VALUES ($1, $2, $3, $4)`

	db := getConnection()
	defer db.Close()

	stmt, err := db.Prepare(q)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	r, err := stmt.Exec(p.ID, p.Name, p.Age, p.Address)
	if err != nil {
		return err
	}
	i, _ := r.RowsAffected()
	if i != 1 {
		return fmt.Errorf("Error 1")
	}
	return nil

}

func getConnection() *sql.DB {
	dns := mydata.DataAccess
	db, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal("error: ", err)
	}
	return db
}
