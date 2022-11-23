package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"database/sql"

	"github.com/xuri/excelize/v2"

	_ "github.com/go-sql-driver/mysql"
)

const columnString = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	envDbURL        string
	envQuery        string
	envHeader       string
	envXLSXFileName string
	db              *sql.DB
)

func init() {
	envDbURL = GetEnvValue("DATABASE_URL")
	envQuery = GetEnvValue("QUERY")
	envHeader = GetEnvValue("HEADER")
	envXLSXFileName = GetEnvValue("XLSX_FILENAME")
}

func getData() [][]string {
	rows, err := db.Query(envQuery)
	defer rows.Close()

	checkError("Error creating the query", err)

	lines := make([][]string, 0)

	header := strings.Split(envHeader, ",")
	lines = append(lines, header)

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		checkError("Error getting columns from table", err)
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// now let's loop through the table lines and append them to the slice declared above
	for rows.Next() {
		// read the row on the table
		// each column value will be stored in the slice
		err = rows.Scan(scanArgs...)

		checkError("Error scanning rows from table", err)

		var value string
		var line []string

		for _, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
				line = append(line, value)
			}
		}

		lines = append(lines, line)
	}

	checkError("Error scanning rows from table", rows.Err())

	return lines

}

func main() {
	log.Println("THINLACE STARTED.")
	var err error

	db, err = sql.Open("mysql", envDbURL)
	defer db.Close()

	if err != nil {
		panic(err)
	}

	f := excelize.NewFile()
	index := f.NewSheet("Sheet1")
	for i, row := range getData() {
		for j, col := range row {
			column := ""
			x := j / 26
			if x == 0 {
				column = string(columnString[j])
			} else {
				y := j % 26
				column = string(columnString[x-1]) + string(columnString[y])
			}
			f.SetCellValue("Sheet1", column+strconv.Itoa(i+1), col)
		}
	}

	f.SetActiveSheet(index)
	if err := f.SaveAs(envXLSXFileName); err != nil {
		fmt.Println(err)
	}

	log.Println(envXLSXFileName, "generated successfully.")
	log.Println("THINLACE COMPLETED.")
}

func GetEnvValue(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Sorry!!! ENV %s is not set", key)
	}
	return val
}

// helper function to handle errors
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
