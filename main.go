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
    Distance float64 `json:"distance"'`
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
    // query to retrieve the spots data
    sqlQuery := `
    SELECT id,website,description,ST_AsText(coordinates),name,rating from "MY_TABLE"
    ;
    `
    rows,err := db.Query(sqlQuery);

    if err !=nil{
        log.Fatal(err);
    }
    defer rows.Close();

    // Our supplied parameters
    suppliedPoint := fmt.Sprintf("POINT(%f %f)",suppliedParams.Longitude,suppliedParams.Latitude);

    /*  Control flow Management 


    [SORTING LOGIC] 
    To manage the final returned object of this request, initially we are going to have two objects: within50 and outside50.
    They represent spots that are within and outside 50 metres from our supplied point respectively (considering radius param).

    Sorting each of these objects as such:
    within50 -> sorted by rating
    outside50 -> sorted by distance
    To do the sort by distance, as distance is not a field in our data, I have created an array of type spotDistance with this added field. 
    That way I am able to sort with ease before finally iterating through and removing that field. 

    Once we have both of these sorted we can now append them together as so:

    finalResult = append(within50,outside50...);
    This will get us the final form we want to finally return through json
    */

    var outside50 []spotDistance;
    var inside50 []spot;

    for rows.Next(){
        var rating sql.NullFloat64;
        var id,website,description,coordinates,name sql.NullString;

        if err := rows.Scan(&id,&name,&website,&coordinates,&description,&rating);err != nil{
            log.Fatal(err);
        }
        //Now initialize our struct object to then manage later
        var spotInstance spot;
        spotInstance.Id = id.String;
        spotInstance.Coordinates = coordinates.String;
        spotInstance.Name=name.String;
        spotInstance.Website = website.String;
        spotInstance.Rating=rating.Float64;

        /* We must manage both circle and square boundaries. 
        For Circle: Get distance between both points;
        For Square: Use envelope 
        */

        // Arguments to this query will be the two different points: the supplied input location and the current row location
        getDistanceQuery := fmt.Sprintf(`SELECT ST_Distance(
            '%s'::geography,
            '%s'::geography
        );`,coordinates.String,suppliedPoint);

        // fmt.Println(getDistanceQuery);

        distanceResult,disterr := db.Query(getDistanceQuery);
        if disterr !=nil{
            log.Fatal(err);
        }
        defer distanceResult.Close(); 
        for distanceResult.Next(){
            var dist float64;
            distanceResult.Scan(&dist);

            if suppliedParams.Type=="circle"{
                CircleHandler(&inside50,&outside50,dist,rating,suppliedParams,id,coordinates,name,website);
            }else{
                SquareHandler(&inside50,&outside50,suppliedParams,suppliedPoint,dist,rating,id,coordinates,name,website);
            }
        }
    }
    // Now we have completed our checks through all the spots and assigned in correct objects
    // we must now do the sorts mentioned 

    sort.Slice(outside50,func(i,j int)bool{
        return outside50[i].Distance < outside50[j].Distance;
    }) 
    sort.Slice(inside50,func(i,j int)bool{
        return inside50[i].Rating > inside50[j].Rating;
    })
    var outside50Filtered []spot;

    for _,singleSpot := range outside50{ // here we are reforming our outside50 but 
        var filteredSpot spot;
        filteredSpot.Id = singleSpot.Id;
        filteredSpot.Coordinates = singleSpot.Coordinates;
        filteredSpot.Rating = singleSpot.Rating;
        filteredSpot.Name = singleSpot.Name;
        filteredSpot.Website = singleSpot.Website;

        outside50Filtered = append(outside50Filtered,filteredSpot);
    }
    // so now we have our properly sorted inside50 and properly sorted outside50
    // we can then combine in result and send json object

    var result []spot;
    result = append(inside50,outside50Filtered...);

    c.IndentedJSON(http.StatusOK,result);
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
    /*________________ DATA SETUP ________________ */
    if setup{
        path := filepath.Join("spots.sql");
        data, ioError := ioutil.ReadFile(path)
        if ioError != nil {fmt.Println("Error retrieving .sql file: ", ioError);}
        sql := string(data)
        _, error := db.Exec(sql);
        if error != nil {fmt.Println("Error executing sql: ",error)
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
