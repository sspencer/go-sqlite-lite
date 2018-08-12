package main

import (
	"fmt"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
)

func main() {
	err := main2()
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func main2() error {
	conn, err := sqlite3.Open("mydatabase.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer conn.Close()

	err = conn.Exec(`CREATE TABLE student(name STRING, age INTEGER)`)
	if err != nil {
		return fmt.Errorf("failed to create students table: %v", err)
	}

	err = withTx(conn, func() error {
		return insertStudents1(conn)
	})
	if err != nil {
		return fmt.Errorf("failed to insert students: %v", err)
	}

	err = queryStudents1(conn)
	if err != nil {
		return fmt.Errorf("failed to query students: %v", err)
	}

	return nil
}

// withTx performs a function wrapped with a Begin and performs a Rollback if
// an error occurs, or Commit otherwise
func withTx(conn *sqlite3.Conn, f func() error) error {
	if err := conn.Begin(); err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Perform work inside the transaction
	err := f()
	if err != nil {
		conn.Rollback()
		return err
	}
	if err = conn.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}
	return nil
}

func insertStudents1(conn *sqlite3.Conn) error {
	// Create a prepared statement
	stmt, err := conn.Prepare(`INSERT INTO student VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare to insert to students table: %v", err)
	}
	defer stmt.Close()

	// This is how you can Bind arguments, Step the statement, then Reset it
	if err = stmt.Bind("Bill", 18); err != nil {
		return fmt.Errorf("failed to bind arguments: %v", err)
	}
	if _, err = stmt.Step(); err != nil {
		return fmt.Errorf("failed to step student insert: %v", err)
	}
	if err = stmt.Reset(); err != nil {
		return fmt.Errorf("failed to reset student insert: %v", err)
	}

	// Even more convenient, Exec will call Bind, Step as many times as needed
	// and always Reset the statement
	if err = stmt.Exec("Tom", 18); err != nil {
		return fmt.Errorf("failed to insert student: %v", err)
	}
	if err = stmt.Exec("John", 19); err != nil {
		return fmt.Errorf("failed to insert student: %v", err)
	}
	if err = stmt.Exec("Bob", 18); err != nil {
		return fmt.Errorf("failed to insert student: %v", err)
	}

	return nil
}

func queryStudents1(conn *sqlite3.Conn) error {
	// Prepare can prepare a statement and optionally also bind arguments
	stmt, err := conn.Prepare(`SELECT * FROM student WHERE age = ?`, 18)
	if err != nil {
		return fmt.Errorf("failed to select from students table: %v", err)
	}
	defer stmt.Close()

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return fmt.Errorf("step failed while querying students: %v", err)
		}
		if !hasRow {
			break
		}

		// Use Scan to access column data from a row
		var name string
		var age int
		err = stmt.Scan(&name, &age)
		if err != nil {
			return fmt.Errorf("scan failed while querying students: %v", err)
		}

		fmt.Println("name:", name, "age:", age)
	}

	return nil
}