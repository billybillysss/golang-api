package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// selectSql generates a SELECT SQL query based on the given table name, and optionally primary keys
func selectSql(table interface{}, tableName string, id ...int) string {

	var idStr []string
	var condition string
	colNames, pkColName := getColumns(table) // Get the column names and the primary key column name

	for _, value := range id {
		idStr = append(idStr, strconv.Itoa(value))
	}

	if len(id) > 0 {
		condition = " WHERE " + pkColName + " IN (" + strings.Join(idStr, ",") + ")"
	} else {
		condition = ""
	}

	cols_string := strings.Join(colNames, ", ") // Join the column names with commas

	return fmt.Sprintf("SELECT %s FROM %s", cols_string, tableName) + condition // Return the generated SQL query
}

// deleteSql generates a DELETE SQL query based on the given table name and primary key
func deleteSql(table interface{}, tableName string, id int) string {
	_, pkColName := getColumns(table) // Get the primary key column name
	idStr := strconv.Itoa(id)         // Convert the primary key value to a string

	return fmt.Sprintf("DELETE FROM %s WHERE %s = %s", tableName, pkColName, idStr) // Return the generated SQL query
}

// updateOrInsertSql generates an SQL query based on the given method and JSON string
func updateOrInsertSql(tableName, jsonStr, method string) string {

	// Skip columns are columns that should be ignored when doing an update or insert
	updateOrInsertSkipColumns := strings.Split(skipColumns, ",") // Split the string by comma
	for i, col := range updateOrInsertSkipColumns {
		updateOrInsertSkipColumns[i] = strings.TrimSpace(col) // Trim leading and trailing whitespace from each part
	}

	res := make(map[string]any) // The map to store the JSON data
	var columns []string        // The columns to be used in the SQL query
	var member Member           // The member struct to validate the JSON data

	// Schema validation
	json.Unmarshal([]byte(jsonStr), &member) // Unmarshal the JSON data into the member struct

	data, _ := json.Marshal(member) // Marshal the member struct back into JSON data
	json.Unmarshal(data, &res)      // Unmarshal the JSON data into the res map

	primaryKeys := getPK(member) // Get the primary keys of the table

	if method == "update" {

		var condition []string // The condition to be used in the SQL query

		for key, value := range res {

			// If the column is a primary key, add it to the condition
			if inColumns(key, primaryKeys) {

				// If the value is a number, add it without quotes
				if isNumeric(value) {
					condition = append(condition, fmt.Sprintf("%s = %v", key, value))

					// If the value is not a number, add it with quotes
				} else {
					condition = append(condition, fmt.Sprintf("%s = '%v'", key, value))
				}

				// If the column is in the skip columns, skip it
			} else if inColumns(key, updateOrInsertSkipColumns) {
				continue

				// If the column is not a primary key or in the skip columns, add it to the columns to be updated
			} else {

				// If the value is a number, add it without quotes
				if isNumeric(value) {
					columns = append(columns, fmt.Sprintf("%s = %v", key, value))

					// If the value is not a number, add it with quotes
				} else {
					columns = append(columns, fmt.Sprintf("%s = '%v'", key, value))
				}
			}
		}

		tableStr := "UPDATE " + tableName + " SET " // The beginning of the SQL query

		updateStr := strings.Join(columns, ", ") // The part of the SQL query that sets the columns to be updated

		conditionStr := " WHERE " + strings.Join(condition, " and ") // The part of the SQL query that sets the condition

		return tableStr + updateStr + conditionStr // Return the full SQL query

		// If the method is "insert"
	} else if method == "insert" {

		var values []string // The values to be used in the SQL query

		for key, value := range res {

			// If the column is in the skip columns or a primary key, skip it
			if inColumns(key, updateOrInsertSkipColumns) || inColumns(key, primaryKeys) {
				continue

				// If the column is not in the skip columns or a primary key, add it to the columns and values
			} else {

				// If the value is a number, add it without quotes
				if isNumeric(value) {
					columns = append(columns, key)
					values = append(values, fmt.Sprintf("%v", value))

					// If the value is not a number, add it with quotes
				} else {
					columns = append(columns, key)
					values = append(values, fmt.Sprintf("'%v'", value))
				}
			}
		}

		tableStr := "INSERT INTO " + tableName + " " // The beginning of the SQL query

		insertStr := fmt.Sprintf("(%s) VALUES (%s)", strings.Join(columns, ", "), strings.Join(values, ", ")) // The part of the SQL query that sets the columns and values

		return tableStr + insertStr // Return the full SQL query

		// If the method is not "update" or "insert", return an empty string
	} else {
		return ""
	}
}

// isNumeric checks if a given value is numeric
func isNumeric(value any) bool {

	switch reflect.TypeOf(value).Kind() {

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}

}

// getPK returns the primary key column names of a given struct
func getPK(table interface{}) []string {
	/*
		This function takes a struct as input and returns a slice of strings
		containing the column names of the primary key.
	*/

	reflectType := reflect.TypeOf(table)
	var pks []string
	nFields := reflectType.NumField()

	for i := 0; i < nFields; i++ {

		// Lookup the "pk" tag in the field's struct tags
		pk, isFound := reflectType.Field(i).Tag.Lookup("pk")
		if !isFound {
			// if the tag is not found, move on to the next field
			continue
		}
		pks = append(pks, pk)
	}
	return pks
}

// inColumns checks if a given column name is present in a slice of strings
func inColumns(column string, columns []string) bool {
	for _, col := range columns {
		if column == col {
			return true
		}
	}
	return false
}

// getColumns returns a slice of column names and the primary key column name
// of a given struct. The function uses the "db" struct tag to get the column
// names and the "pk" struct tag to get the primary key column name.
func getColumns(table interface{}) ([]string, string) {
	v := reflect.TypeOf(table)

	colNames := make([]string, 0, v.NumField()-1) // make slice to hold column names

	var pkColName string // primary key column name
	for i := 0; i < v.NumField(); i++ {
		// Lookup the "db" tag in the field's struct tags
		colName, columnExist := v.Field(i).Tag.Lookup("db")
		// Lookup the "pk" tag in the field's struct tags
		pk, pkExist := v.Field(i).Tag.Lookup("pk")

		// If the "db" tag is not found and the "pk" tag is not found, move on to the next field
		if !columnExist && !pkExist {
			continue
			// If the "db" tag is not found but the "pk" tag is found, set the primary key column name
		} else if !columnExist && pkExist {
			pkColName = pk
			// If the "db" tag is found and the "pk" tag is not found, add the column name to the slice
		} else if columnExist && !pkExist {
			colNames = append(colNames, colName)
			// If the "db" tag is found and the "pk" tag is found, add the column name to the slice
			// and set the primary key column name
		} else {
			colNames = append(colNames, colName)
			pkColName = pk
		}
	}
	return colNames, pkColName
}
