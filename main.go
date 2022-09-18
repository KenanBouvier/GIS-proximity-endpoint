package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
    "log"

	"github.com/gin-gonic/gin"
    "github.com/omniscale/imposm3/geom/geos"

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

    // Send true/false for spots.sql execution yes/no respectively
    dataSetup(false);

    /*  __________________TASK 1 | Query________________________*/

    task1();


    /*  __________________TASK 2 | Endpoint________________________*/

    router:=gin.Default(); 
    router.GET("/proximity",locationProximityRoute);
    router.Run("localhost:8080");
}

type proximity struct{
    Longitude float64 `json:"longitude"`
    Latitude float64 `json:"latitude"`
    Radius float64 `json:"radius"`
    Type string `json:"type"`
}
type spot struct{
    
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
    // here we query all rows but the coordinates must be in WKT format to process.
    sqlQuery := `
    SELECT id,website,description,ST_AsText(coordinates),name,rating from "MY_TABLE"
    LIMIT 10;
    `
    rows,err := db.Query(sqlQuery);
 
    if err !=nil{
        log.Fatal(err);
    }
    defer rows.Close();
    println(rows);
    for rows.Next(){
        var rating sql.NullFloat64;
        var id,website,description,coordinates,name sql.NullString;

        if err := rows.Scan(&id,&name,&website,&coordinates,&description,&rating);err != nil{
            log.Fatal(err);
        }
        fmt.Println(id.String,name.String,website.String,coordinates.String,description.String,rating.Float64);

        g := geos.NewGeos();

        geom := g.FromWkt(coordinates.String);
        if geom == nil{
            log.Fatal("Error reading wkt");
        }
        fmt.Println(geom.Area());
    }
}

func squareController(){

}

func task1(){

    /*  __________________ Part 1 | Query ________________________

    /* Q: Change the website field so that it only contains the domain*/
    
    /* Essentially, we are applying a regex pattern that is matching the domain of a given url and updating the field with this match*/
    websiteFieldToDomainQuery := `
        UPDATE "MY_TABLE"
        SET website = substring(website from '(?:.*://)?(?:www\.)?([^/?]*)')
    `;
    executeSQL(websiteFieldToDomainQuery);
    /* Essentially, we are applying a regex pattern that is matching the domain of a given url and updating the field with this match*/

    

    /*  __________________Part 2 | Query________________________ */

    /* Q: Count how many spots contain the same domain */
    
    /* Grouping by website to find the count and getting the total count of those occurences */
    multipleSpotsCountQuery :=`
        select COUNT(*) OVER () AS TotalRecordCount
        from "MY_TABLE"
        group by website
        having count(website)>1
        limit 1;
    `;
    executeSQL(multipleSpotsCountQuery); 



    /*  __________________Part 3 | Query________________________ */

    /* Q: Return spots which have a domain with a count greater than 1 */
    
    /* As we want the full record of spots, we select all from those that when grouped by website count is greater than 1  */
    spotsWithMultipleCountQuery :=`
    SELECT * FROM "MY_TABLE"
    WHERE website IN (SELECT website
        FROM "MY_TABLE"
        GROUP BY website
        HAVING COUNT(website) > 1 
    );
    `;
    executeSQL(spotsWithMultipleCountQuery); 



    /*  __________________Part 4 | Query________________________ */

    /* Q: Make a PL/SQL function for point 1. */
    
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
        fmt.Println("Error executing sql query: ",error);
    }else{
        fmt.Println("Successful SQL query execution!");
    }
}
