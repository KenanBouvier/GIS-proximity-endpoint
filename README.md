# Task 2 Details


# Routing
There is one /proximity route that takes in the expected parameters.
Then the ProximityController handles both cases of circle and square.
For Circle: Get distance between both points. As center position to boundary/radius is always constant
For Square: Create box boundary from center position with radius to determine intersection with each spot location

# Circle & Square Handlers 
Overview of handler operations:
    For Circle: Get distance between both points. As center position to boundary/radius is always constant
    For Square: Create box boundary from center position with radius to determine intersection with each spot location

The time complexity for circle handler is constant.
Although the time complexity for the square handler is also constant, it has a very noticable factor of about 4x slower.
The inherent difficulty with squares instead of circles is that the valid area is different depending on the angle we are directing towards. Unlike a circle in which at every angle the valid distance is constant.

To make our square handler more efficient, we can assume a circle from the square. A circle that contains the whole square. Else we are removing valid areas that should be considered.
Therefore our new circle will have radius of the furthest distance from the center any border of the square. This is the corner of the square.
This is also sqrt(2)*halfLengthOfSquare which in our case = the inputted radius parameter.

Given the nature of this specific endpoint, the majority of our spots in our db will not fall within the supplied radius. Therefore an overwhelming number of operations will be dealt with using this outer circle and by consequence, increase efficiency. 

One extra thing I included was that if the data was within our larger circle(squarehandler) then we can also skip our squarehandler operations if the smaller circle within the square contained our spot. This way we can instantly know. This leaves an even smaller area of locations that will result in our slower square handling.


# Performance
## Time to complete request - (Before Efficient Algorithm)
### First 5 requests for circle, Next 5 requests for square
![image Inefficient](./images/inefficientAlgorithm.png)

## Time to complete request - (With Efficient Algorithm)
### All 5 requests for square
![image Efficient](./images/efficientAlgorithm.png)

We can see that the inefficiences involved with square lookups have now become non-existent

# [SORTING LOGIC] 
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

