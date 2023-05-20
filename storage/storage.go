package storage

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/temelpa/timetravel/entity"
)

const databaseFile = "sqlite-database.db"

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	log.Println("Initializing storage")

	if _, err := os.Stat(databaseFile); err == nil {
		log.Println("Database file already exists, skipping creation.")
	} else {
		err = createDatabase()
		if err != nil {
			return nil, err
		}
		log.Println(databaseFile, "created")
	}

	db, err := sql.Open("sqlite3", "./"+databaseFile)
	if err != nil {
		return nil, err
	}

	return &Storage{db: db}, nil
}

func createDatabase() error {
	os.Remove(databaseFile) // Delete the file to avoid duplicated records.

	file, err := os.Create(databaseFile) // Create SQLite file
	if err != nil {
		return err
	}
	file.Close()

	db, err := sql.Open("sqlite3", "./"+databaseFile) // Open the created SQLite file
	if err != nil {
		return err
	}
	defer db.Close()

	err = createTable(db)
	if err != nil {
		return err
	}

	return nil
}

func createTable(db *sql.DB) error {
	createRecordsTableSQL := `CREATE TABLE records (
		"id" integer NOT NULL PRIMARY KEY,
		"data" TEXT,
		"created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		"deleted_at" TIMESTAMP	
	);` // SQL Statement for Create Table

	log.Println("Creating records table...")
	_, err := db.Exec(createRecordsTableSQL) // Execute SQL Statement
	if err != nil {
		return err
	}
	log.Println("Records table created")

	return nil
}

func (s *Storage) InsertRecord(id int, data string) (int, error) {
	log.Println("Inserting record...")
	insertRecordSQL := `INSERT INTO records (id, data) VALUES (?, ?)`

	statement, err := s.db.Prepare(insertRecordSQL)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	_, err = statement.Exec(id, data)
	if err != nil {
		return 0, err
	}

	return 0, nil
}

func (s *Storage) GetRecordByID(id int) (*entity.Record, error) {
	log.Println("Getting record...")
	getRecordSQL := `SELECT id, data, created_at, deleted_at FROM records WHERE id = ?`

	statement, err := s.db.Prepare(getRecordSQL)
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	row := statement.QueryRow(id)
	record := &entity.Record{}
	var data string
	var nullableString sql.NullString
	err = row.Scan(&record.ID, &data, &record.CreatedAt, &nullableString)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	log.Println(data)
	err = json.Unmarshal([]byte(data), &record.Data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return record, nil
}

func (s *Storage) UpdateRecord(id int, data string) (int, error) {
	log.Println("Updating record...")
	updateRecordSQL := `UPDATE records SET data = ? WHERE id = ?`

	statement, err := s.db.Prepare(updateRecordSQL)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	_, err = statement.Exec(data, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) DeleteRecord(id int) error {
	log.Println("Deleting record...")
	deleteRecordSQL := `UPDATE records SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`

	statement, err := s.db.Prepare(deleteRecordSQL)
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return err
	}

	return nil
}