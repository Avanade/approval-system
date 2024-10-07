package database

import (
	"context"
	"database/sql"
	"fmt"
	"main/config"

	_ "github.com/microsoft/go-mssqldb"
)

type Database struct {
	connString string
}

func NewDatabase(config config.ConfigManager) Database {
	fmt.Println("ConnectDb New")
	connectionString := config.GetDatabaseConnectionString()
	return Database{
		connString: connectionString,
	}
}

func (d *Database) Connect() (*sql.DB, error) {
	conn, err := sql.Open("sqlserver", d.connString)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (d *Database) Query(query string, args ...any) (*sql.Rows, error) {
	con, err := d.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()

	ctx := context.Background()
	rows, err := con.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *Database) QueryRow(query string, args ...any) (*sql.Row, error) {
	con, err := d.Connect()
	if err != nil {
		return nil, err
	}
	defer con.Close()

	ctx := context.Background()
	row := con.QueryRowContext(ctx, query, args...)
	if row == nil {
		err = fmt.Errorf("QueryRowContext returned nil")
	}
	return row, nil
}

func (d *Database) Execute(query string, args ...any) error {
	con, err := d.Connect()
	if err != nil {
		return err
	}
	defer con.Close()

	ctx := context.Background()
	_, err = con.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

// helper function to convert rows.Scan() to map[string]interface{}
func (d *Database) RowsToMap(rows *sql.Rows) ([]map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}
		err := rows.Scan(pointers...)
		if err != nil {
			return nil, err
		}
		result := make(map[string]interface{})
		for i, val := range values {
			result[columns[i]] = val
		}
		results = append(results, result)
	}

	return results, nil
}
