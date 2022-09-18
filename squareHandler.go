package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"
)

func SquareHandler(inside50 *[]spot,outside50 *[]spotDistance, suppliedParams proximity,suppliedPoint string,dist float64,rating sql.NullFloat64, id,coordinates,name,website sql.NullString) {

    /*
    The STProjectQueries are defined as such:

        First param is our point. 
        Second param is distance
        Third param is bearing(from north clockwise)

    We can calculate our values through some pythagoras: 

        distance=sqrt(2)*radius
        angle = 225 or 45  (degrees)

    */

    var cornerDistance float64 = math.Sqrt(2)*suppliedParams.Radius;
    var bottomLeftangle float64 = 225;
    var topRightangle float64 = 45;

    bottomLeftSTProjectQuery := fmt.Sprintf(`
    SELECT ST_AsText(ST_Project('%s'::geography,%f,radians(%f)));`,
    suppliedPoint,cornerDistance,bottomLeftangle);

    topRightSTProjectQuery := fmt.Sprintf(`
    SELECT ST_AsText(ST_Project('%s'::geography,%f,radians(%f)));`,
    suppliedPoint,cornerDistance,topRightangle);

    var bottomLeft,topRight sql.NullString;

    bottomLeftResult,projectError := db.Query(bottomLeftSTProjectQuery);
    if projectError !=nil{
        log.Fatal(err);
    }
    defer bottomLeftResult.Close(); 

    for bottomLeftResult.Next(){
        bottomLeftResult.Scan(&bottomLeft); 
    }

    toprightResult,projectError := db.Query(topRightSTProjectQuery);
    if projectError !=nil{
        log.Fatal(err);
    }
    defer toprightResult.Close(); 

    for toprightResult.Next(){
        toprightResult.Scan(&topRight); 
    }

    /* 
    The boundingBox query requires commas between the long and lat.
    Whereas our POINT(x y) does not have this.
    To solve this: I can split the string by delimiter ' ' and get our correct string form:
    formattedPoint := splitPoint[0]+','+splitPoint[1];
    */

    var bottomLeftSplitPoint []string = strings.Split(bottomLeft.String, " ");
    var topRightSplitPoint []string = strings.Split(topRight.String, " ");
    var coordinatesSplitPoint []string = strings.Split(coordinates.String, " ");

    var bottomLeftFormatted string = bottomLeftSplitPoint[0]+","+bottomLeftSplitPoint[1];
    var topRightFormatted string = topRightSplitPoint[0]+","+topRightSplitPoint[1];
    var coordinatesFormatted string = coordinatesSplitPoint[0]+","+coordinatesSplitPoint[1];

    // coordinates formatted is the row's location 
    boundingBoxQuery := fmt.Sprintf(`
    SELECT ST_%s && 
    ST_MakeBox2D(ST_%s, ST_%s) 
    AS overlaps;
    `,coordinatesFormatted,bottomLeftFormatted,topRightFormatted);

    // fmt.Println(boundingBoxQuery);

    var withinBoundary sql.NullBool;

    boundingResult,boundingError := db.Query(boundingBoxQuery);
    if boundingError != nil{
        log.Fatal(boundingError);
    }
    defer boundingResult.Close();
    for boundingResult.Next(){
        boundingResult.Scan(&withinBoundary);
    }
    
    // inside50 and also within supplied boundary
    if dist<=50 && withinBoundary.Bool{
        var spotInstance spot;
        spotInstance.Id = id.String;
        spotInstance.Coordinates = coordinates.String;
        spotInstance.Name=name.String;
        spotInstance.Website = website.String;
        spotInstance.Rating=rating.Float64;

        *inside50 = append(*inside50,spotInstance);
    }else if withinBoundary.Bool { // outside50
        var spotInstance spotDistance;
        spotInstance.Id = id.String;
        spotInstance.Coordinates = coordinates.String;
        spotInstance.Name=name.String;
        spotInstance.Website = website.String;
        spotInstance.Rating=rating.Float64;
        spotInstance.Distance=dist;

        *outside50 = append(*outside50,spotInstance);
    }else{
        // we don't consider the spot as it is outside our considering area
    }

}
