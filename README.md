# thinlace

A go program to read data from MySQL & output Excel XLSX.

Please make sure env is set already
```
DATABASE_URL=mysqlUN:mysqlPWD@tcp(mysqlHOST:3306)/mysqlDB
QUERY=SELECT *, DATE_FORMAT(created_at, '%W %M %Y') FROM dummy
HEADER=ID,NAME,AGE,GENDER,CREATED_AT,FORMATTED_CREATED_AT
XLSX_FILENAME=Book5.xlsx
````