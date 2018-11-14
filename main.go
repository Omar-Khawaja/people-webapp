package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var files = []string{"index.html"}
var t = template.Must(template.ParseFiles(files...))

type Name struct {
	FirstName string
	LastName  string
}

type configuration struct {
	Host     string
	Port     int
	Username string
	Password string
	DB       string
}

func main() {
	port := 8080
	http.HandleFunc("/names", namesHandler)
	fmt.Printf("Starting server on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil))
}

func namesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := dbConn()
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err)
		return
	}
	defer db.Close()

	selDB, err := db.Query("SELECT * FROM people;")
	if err != nil {
		log.Println(err)
		fmt.Fprint(w, err)
		return
	}
	name := Name{}
	nameArray := []Name{}
	for selDB.Next() {
		var firstname, lastname string
		err = selDB.Scan(&firstname, &lastname)
		if err != nil {
			log.Println(err)
			return
		}
		name.FirstName = firstname
		name.LastName = lastname
		nameArray = append(nameArray, name)
	}
	t.Execute(w, nameArray)
}

func dbConn() (db *sql.DB, err error) {
	file, err := os.Open("config.json")
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		return db, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.Username, conf.Password, conf.DB)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}
	return db, nil
}
