package models

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

type JobStatus string

const (
	StatusActive    JobStatus = "active"
	StatusPending   JobStatus = "pending"
	StatusCompleted JobStatus = "completed"
	StatusIssue     JobStatus = "issue"
	StatusCanceled  JobStatus = "canceled"
)

type Status string

const (
	Occupied  Status = "ocupied"
	Pending   Status = "pending"
	Available Status = "аvailable"
	Issue     Status = "issue"
	Archived  Status = "аrchived"
)

type Job struct {
	ID           int       `db:"id"`
	Description  string    `db:"description"`
	StartDate    time.Time `db:"start_date"`
	EndDate      time.Time `db:"end_date"`
	Driver_id    int       `db:"driver_id"`
	Truck_id     int       `db:"truck_id"`
	Status       JobStatus `db:"status"`
	TruckStartKm float64   `db:"starting_km"`
	TruckEndKm   float64   `db:"ending_km"`
	Fuel_spent   float64   `db:"fuel_spent"`
	Expenses     float64   `db:"expenses"`
	Revenue      float64   `db:"revenue"`
	CreatedAt    time.Time `db:"created_at"`
}

type Driver struct {
	ID             int       `db:"id"`
	Name           string    `db:"name"`
	License_number string    `db:"license_number"`
	Phone_number   string    `db:"phone_number"`
	Created        time.Time `db:"created_at"`
	Status         Status    `db:"status"`
}

type Truck struct {
	ID            int       `db:"id"`
	Model         string    `db:"model"`
	License_plate string    `db:"license_plate"`
	Created       time.Time `db:"created_at"`
	Status        Status    `db:"status"`
}

type DbModel struct {
	DB *sql.DB
}

// Insert function that accepts any table name and struct data
func (m *DbModel) Insert(table_name string, data interface{}) (int, error) {
	// Use reflection to inspect the struct fields
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	var columns []string
	var values []interface{}

	// Iterate through the struct fields
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		columnName := field.Tag.Get("db") // Use struct tags for DB column names

		// Skip ID, Created and Status fields
		if columnName == "id" || columnName == "created_at" || columnName == "status" {
			continue
		}

		if columnName == "" {
			columnName = field.Name
		}
		columns = append(columns, columnName)
		values = append(values, val.Field(i).Interface())
	}

	// Generate the query dynamically
	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		table_name,
		strings.Join(columns, ", "),
		strings.Join(generatePlaceholders(len(columns)), ", "),
	)

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(insertedID), nil
}

func generatePlaceholders(n int) []string {
	placeholders := make([]string, n)
	for i := range placeholders {
		placeholders[i] = "?"
	}
	return placeholders
}

func (m *DbModel) Get(table_name string, result interface{}, params map[string]interface{}) error {
	slice := reflect.ValueOf(result).Elem()
	elemType := slice.Type().Elem().Elem()

	var columns []string

	sampleInstance := reflect.New(elemType).Elem()

	for i := 0; i < sampleInstance.NumField(); i++ {
		field := elemType.Field(i)
		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = field.Name
		}
		columns = append(columns, columnName)
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s",
		strings.Join(columns, ", "),
		table_name,
	)

	var whereClauses []string
	var args []interface{}

	for param, value := range params {
		whereClauses = append(whereClauses, fmt.Sprintf("%s = ?", param))
		args = append(args, value)
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	log.Printf("Executing query: %s with params: %v", query, args)

	rows, err := stmt.Query(args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		newInstance := reflect.New(elemType).Elem()
		scan_dest := make([]interface{}, len(columns))

		for i := range scan_dest {
			scan_dest[i] = newInstance.Field(i).Addr().Interface()
		}

		if err := rows.Scan(scan_dest...); err != nil {
			return err
		}

		slice.Set(reflect.Append(slice, newInstance.Addr()))
	}

	if slice.Len() == 0 {
		return nil
	}

	return nil
}

func (m *DbModel) Update(tableName string, id int, data interface{}) error {
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	var setClauses []string
	var values []interface{}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		columnName := field.Tag.Get("db")

		// Skip fields with no 'db' tag or invalid fields
		if columnName == "" {
			continue
		}

		// Skip fields with zero values (not explicitly set)
		fieldValue := val.Field(i).Interface()
		if isZero(fieldValue) {
			continue
		}

		// Add column name and value to the query
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", columnName))
		values = append(values, fieldValue)
	}

	// If no fields are set, return an error
	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Construct the query
	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = ?",
		tableName,
		strings.Join(setClauses, ", "),
	)

	// Add the ID to the end of the values slice
	values = append(values, id)

	// Prepare and execute the query
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}
func isZero(value interface{}) bool {
	return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
}
