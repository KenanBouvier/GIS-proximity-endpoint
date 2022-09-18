SELECT * FROM "MY_TABLE"
WHERE website IN (SELECT website
    FROM "MY_TABLE"
    GROUP BY website
    HAVING COUNT(website) > 1 
)
and length(website)>0
order by website asc ;

-- As we want the full record of the spots based on the website count. We can surround the GROUP BY with a select.

-- As some website fields are empty or null we must consider these.
-- We can check if the length of the data in website field >0  
-- As for the null fields, it is already ignored in COUNT
