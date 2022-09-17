package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

/*_________________ DB setup start ______________________ */
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

/*_________________ DB setup end ______________________ */


func main(){

    // Send true or false for spots.sql execution or to skip respectively
    dataSetup(false);

    /*  __________________TASK 1 | Query________________________*/

    task1();


    /*  __________________TASK 2 | Endpoint________________________*/

    router:=gin.Default(); 
    router.GET("/proximity",locationProximityRoute);

    router.Run("localhost:8080");

}
type proximity struct{
    Longitude string `json:"longitude"`
    Latitude string `json:"latitude"`
    Radius float64 `json:"radius"`
    Type string `json:"type"`
}

func locationProximityRoute(c *gin.Context ){
   var newProximity proximity; 

   if err := c.BindJSON(&newProximity);err!=nil{
       c.IndentedJSON(http.StatusNotFound,gin.H{"error":"Error with passed args"})
       return; 
   }
   if newProximity.Type == "circle"{
       c.IndentedJSON(http.StatusAccepted,gin.H{"message":"Sending to circle controller"})
       circleController();
   }else if newProximity.Type=="square"{
       c.IndentedJSON(http.StatusAccepted,gin.H{"message":"Sending to square controller"})
       squareController();
   }else{
       c.IndentedJSON(http.StatusNotFound,gin.H{"error":"Not a valid type. circle or square"})
   }
 
}

func circleController(){
        
}

func squareController(){

}

func task1(){

    /*  __________________ Part 1 | Query ________________________

    /* Change the website field so that it only contains the domain*/
    
    websiteFieldToDomainQuery := `
        UPDATE "MY_TABLE"
        SET website = substring(website from '(?:.*://)?(?:www\.)?([^/?]*)')
    `;
    executeSQL(websiteFieldToDomainQuery);
    /* Essentially we are applying a regex pattern that is matching the domain of a given url */

    
    /*  __________________Part 2 | Query________________________ */

    /* Count how many spots contain the same domain */
    
    multipleSpotsCountQuery :=`
        select COUNT(*) OVER () AS TotalRecords
        from "MY_TABLE"
        group by website
        having count(website)>1
        limit 1;
    `;
    executeSQL(multipleSpotsCountQuery); 


    /*  __________________TASK 1 | Query________________________ */

    /* Part 3) Return spots which have a domain with a count greater than 1 */
    
    spotsWithMultipleCountQuery :=`
    SELECT * FROM "MY_TABLE"
    WHERE website IN (SELECT website
        FROM "MY_TABLE"
        GROUP BY website
        HAVING COUNT(website) > 1 
    );
    `;
    executeSQL(spotsWithMultipleCountQuery); 



    /*  __________________TASK 1 | Query________________________ */

            /* Part 4) Make a PL/SQL function for point 1. */
    
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
func dataSetup(setup bool){
    /*________________ INITIAL DATA SETUP ________________ */

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
}

func executeSQL(sqlStatement string){
    _,error := db.Exec(sqlStatement);
    
    if error != nil{
        fmt.Println("Error executing sql Statement: ",error);
    }else{
        fmt.Println("Successfully Executed!");
    }
}





