package main

import (
    "database/sql"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "path/filepath"
    "sort"
    "github.com/gin-gonic/gin"
    _ "github.com/lib/pq"
)

/*_________________ DB setup start ______________________ */
const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "password"
    dbname   = "spotlas"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
"password=%s dbname=%s sslmode=disable",
host, port, user, password, dbname);
var db, err = sql.Open("postgres", psqlInfo);

/*_________________ DB setup end ______________________ */


// In meters
var RangeToSortRating float64 = 100;

func main(){

    // Set true/false for spots.sql and postgis setup yes/no respectively
    dataSetup(true);    


    /*  __________________TASK 2 | Endpoint________________________*/

    router:=gin.Default(); 
    router.GET("/proximity",proximityRoute);
    router.Run("localhost:8080");
}

type proximity struct{
    Longitude float64 `json:"longitude"`
    Latitude float64 `json:"latitude"`
    Radius float64 `json:"radius"`
    Type string `json:"type"`
}

type spot struct{
    Id string `json:"id"`
    Coordinates string `json:"coordinates"`
    Name string `json:"name"`
    Website string `json:"website"`
    Rating float64 `json:"rating"`
}

// Same as spot type but with extra distance field.
type spotDistance struct{
    Id string `json:"id"`
    Coordinates string `json:"coordinates"`
    Name string `json:"name"`
    Website string `json:"website"`
    Rating float64 `json:"rating"`
    Distance float64 `json:"distance"`
}

func proximityRoute(c *gin.Context ){
    var suppliedParams proximity; 

    // Storing input params from body to access throughout our program
    if err := c.BindJSON(&suppliedParams);err!=nil{
        c.IndentedJSON(http.StatusNotFound,gin.H{"error":"Error with args. Set: longitude, latitude, radius, type"})
        return; 
    }
    if suppliedParams.Type != "circle" && suppliedParams.Type != "square"{
        c.IndentedJSON(http.StatusNotFound,gin.H{"error":"Not a valid type. circle or square"})
        return;
    }
    proximityController(suppliedParams,c);
}

func proximityController(suppliedParams proximity, c *gin.Context){
    sqlQuery := `
    SELECT id,website,description,ST_AsText(coordinates),name,rating from "MY_TABLE";
    `;
    rows,err := db.Query(sqlQuery);
    if err != nil { log.Fatal(err); }
    defer rows.Close();

    // Our supplied parameters
    suppliedPoint := fmt.Sprintf("POINT(%f %f)",suppliedParams.Longitude,suppliedParams.Latitude);


    var outside []spotDistance;
    var inside []spot;

    for rows.Next(){
        var rating sql.NullFloat64;
        var id,website,description,coordinates,name sql.NullString;

        if err := rows.Scan(&id,&name,&website,&coordinates,&description,&rating);err != nil{
            log.Fatal(err);
        }

        /* 
            We must manage both circle and square boundaries. 
            For Circle: Get distance between both points. As center position to boundary/radius is always constant
            For Square: Create box boundary from center position with radius to determine intersection with each spot location
        */

        // Arguments: supplied input location (constant) and the spot location (from db record)
        getDistanceQuery := fmt.Sprintf(`SELECT ST_Distance(
            '%s'::geography,
            '%s'::geography
        );`,coordinates.String,suppliedPoint);

        distanceResult,disterr := db.Query(getDistanceQuery);
        if disterr !=nil{
            log.Fatal(err);
        }
        defer distanceResult.Close(); 
        for distanceResult.Next(){
            var dist float64;
            distanceResult.Scan(&dist);

            if suppliedParams.Type=="circle"{
                CircleHandler(&inside,&outside,dist,rating,suppliedParams.Radius,id,coordinates,name,website);
            }else{
                SquareHandler(&inside,&outside,suppliedParams.Radius,suppliedPoint,dist,rating,id,coordinates,name,website);
            }
        }
    }
    // Now we have completed our checks through all the spots and assigned in correct objects
    // we must now do the sorts mentioned in readme 

    sort.Slice(outside,func(i,j int)bool{
        return outside[i].Distance < outside[j].Distance;
    }) 
    sort.Slice(inside,func(i,j int)bool{
        return inside[i].Rating > inside[j].Rating;
    })
    var outsideFiltered []spot;

    for _,singleSpot := range outside{ // here we are reforming our outside50 but 
        var filteredSpot spot;
        filteredSpot.Id = singleSpot.Id;
        filteredSpot.Coordinates = singleSpot.Coordinates;
        filteredSpot.Rating = singleSpot.Rating;
        filteredSpot.Name = singleSpot.Name;
        filteredSpot.Website = singleSpot.Website;

        outsideFiltered = append(outsideFiltered,filteredSpot);
    }
    // so now we have our properly sorted inside50 and properly sorted outside50
    // we can then combine in result and send json object

    var result []spot;
    result = append(inside,outsideFiltered...);

    c.IndentedJSON(http.StatusOK,result);
}

func dataSetup(setup bool){
    /*________________ DATA SETUP ________________ */
    if !setup{return}
    path := filepath.Join("./Queries/spots.sql");
    data, _ := ioutil.ReadFile(path)
    sql := string(data)
    db.Exec(sql);
    db.Exec(`CREATE EXTENSION postgis;`);
    fmt.Println("Table and data queried + postgis extension set!")
}
