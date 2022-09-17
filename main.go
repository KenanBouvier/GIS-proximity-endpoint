package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

	_ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "password"
    dbname   = "spotlas"
)

func main(){
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        panic(err)
    }

    path := filepath.Join("spots.sql");
    data, ioError := ioutil.ReadFile(path)
    if ioError != nil {
        fmt.Println("Error retrieving .sql file: ", ioError);
    }
    sql := string(data)
    _, error := db.Exec(sql);

    if error != nil {
        fmt.Println("Error executing sql: ",error)
    }else{
        fmt.Println("Successfully inserted!")
    }

    fmt.Println("DB Connection was established")
    
}
