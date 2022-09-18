package main

import (
	"database/sql"
)


// Passing data by reference

func CircleHandler(inside50 *[]spot, outside50 *[]spotDistance , dist float64,rating sql.NullFloat64, radius float64, id,coordinates,name,website sql.NullString) {

    // we don't consider the spot as it is outside our considering area
    if dist > radius{
        return ;
    }
    //Now must be within supplied radius

    // within 50 metres
    if dist<=50 {
        var spotInstance spot;
        spotInstance.Id = id.String;
        spotInstance.Coordinates = coordinates.String;
        spotInstance.Name=name.String;
        spotInstance.Website = website.String;
        spotInstance.Rating=rating.Float64;

        *inside50 = append(*inside50,spotInstance);
    }else{ // outside50 as dist <= radius
        var spotInstance spotDistance;
        spotInstance.Id = id.String;
        spotInstance.Coordinates = coordinates.String;
        spotInstance.Name=name.String;
        spotInstance.Website = website.String;
        spotInstance.Rating=rating.Float64;
        spotInstance.Distance=dist;

        *outside50 = append(*outside50,spotInstance);
    } 
}
