package main

import (
    "fmt"
	"database/sql"
)


// Passing data by reference

func CircleHandler(inside50 *[]spot, outside50 *[]spotDistance , dist float64,rating sql.NullFloat64, suppliedParams proximity, id,coordinates,name,website sql.NullString) {

    // inside50 and also within supplied radius
    if dist<=50 && dist<=suppliedParams.Radius {
        fmt.Println("Within50");
        var spotInstance spot;
        spotInstance.Id = id.String;
        spotInstance.Coordinates = coordinates.String;
        spotInstance.Name=name.String;
        spotInstance.Website = website.String;
        spotInstance.Rating=rating.Float64;

        *inside50 = append(*inside50,spotInstance);
    }else if dist<=suppliedParams.Radius { // outside50
        var spotInstance spotDistance;
        spotInstance.Id = id.String;
        spotInstance.Coordinates = coordinates.String;
        spotInstance.Name=name.String;
        spotInstance.Website = website.String;
        spotInstance.Rating=rating.Float64;
        spotInstance.Distance=dist;

        *outside50 = append(*outside50,spotInstance);
    }else{
        // we don't consider the spot as it is outside our considering area;
    } 
}
