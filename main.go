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
    dbname   = "postgres"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
"password=%s dbname=%s sslmode=disable",
host, port, user, password, dbname);
var db, err = sql.Open("postgres", psqlInfo);

func main(){
    if err != nil {
        panic(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        panic(err)
    }
    fmt.Println("DB Connection established") ;


    // Master switch (var setup) for table initial state setup.

    setup := false;

    if setup{
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
    }


    /*  __________________TASK 1 | Query________________________

        Part 1) sql Query that will change the website fields to contain only domain

    */
    

    /* SOLUTION */

    websiteFieldToDomainQuery := `
        UPDATE "MY_TABLE"
        SET website = substring(website from '(?:.*://)?(?:www\.)?([^/?]*)')
    `;
    executeSQL(websiteFieldToDomainQuery);


    
    /*  __________________TASK 1 | Query________________________

        Part 2) Count how many spots contain the same domain

    */
    
    /* SOLUTION */     
    multipleSpotsCountQuery :=`
        select COUNT(*) OVER () AS TotalRecords
        from "MY_TABLE"
        group by website
        having count(website)>1
        limit 1;
    `;
    executeSQL(multipleSpotsCountQuery); 



    /*  __________________TASK 1 | Query________________________

        Part 3) Return spots which have a domain with a count greater than 1

    */
    
    /* SOLUTION */     
    spotsWithMultipleCountQuery :=`
    SELECT * FROM "MY_TABLE"
    WHERE website IN (SELECT website
        FROM "MY_TABLE"
        GROUP BY website
        HAVING COUNT(website) > 1 
    );      
    `;
    executeSQL(spotsWithMultipleCountQuery); 



    /*  __________________TASK 1 | Query________________________

        Part 4) Make a PL/SQL function for point 1.
        
    */
    
    /* SOLUTION */     
    // spotsWithMultipleCountQuery :=`
    // SELECT * FROM "MY_TABLE"
    // WHERE website IN (SELECT website
    //     FROM "MY_TABLE"
    //     GROUP BY website
    //     HAVING COUNT(website) > 1 
    // );      
    // `;
    // executeSQL(spotsWithMultipleCountQuery); 



}

func executeSQL(sqlStatement string){
    _,error := db.Exec(sqlStatement);
    
    if error != nil{
        fmt.Println("Error executing sql Statement: ",error);
    }else{
        fmt.Println("Successfully Executed!");
    }
}



