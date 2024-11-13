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
	ID                     int       `db:"id"`
	Description            string    `db:"description"`
	Driver_id              int       `db:"driver_id"`
	Truck_id               int       `db:"truck_id"`
	Scheduled_date         time.Time `db:"scheduled_date"`
	Created                time.Time `db:"created_at"`
	Status                 JobStatus `db:"status"`
	Distance               float64   `db:"distance"`
	Package_size           float64   `db:"package_size"`
	Scheduled_arrival_time time.Time `db:"scheduled_arrival_time"`
	Client_name            string    `db:"client_name"`
	Start_location         string    `db:"start_location"`
	Destination_location   string    `db:"destination_location"`
	Package_weight         float64   `db:"package_weight"`
}

type Driver struct {
	ID                  int       `db:"id"`
	Name                string    `db:"name"`
	License_number      string    `db:"license_number"`
	Phone_number        string    `db:"phone_number"`
	Created             time.Time `db:"created_at"`
	Status              Status    `db:"status"`
	Average_consumption float64   `db:"average_consumption"`
	Km_traveled         float64   `db:"km_traveled"`
	Active              bool      `db:"active"`
}

type Truck struct {
	ID                  int       `db:"id"`
	Model               string    `db:"model"`
	License_plate       string    `db:"license_plate"`
	Created             time.Time `db:"created_at"`
	Status              Status    `db:"status"`
	Km_traveled         float64   `db:"km_traveled"`
	Average_consumption float64   `db:"average_consumption"`
	Active              bool      `db:"active"`
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

		if columnName == "id" || columnName == "created_at" {
			continue
		}

		if columnName == "" {
			columnName = field.Name
		}

		setClauses = append(setClauses, fmt.Sprintf("%s = ?", columnName))
		values = append(values, val.Field(i).Interface())
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE id = ?",
		tableName,
		strings.Join(setClauses, ", "),
	)

	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	values = append(values, id)

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}
