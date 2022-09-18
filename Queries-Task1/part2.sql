SELECT COUNT(*) OVER () AS TotalRecordCount
FROM "MY_TABLE"
GROUP BY website
HAVING count(website)>1
LIMIT 1;

-- Grouping by website then using OVER clause to get the total count and then retrieving the top totalRecordCount
