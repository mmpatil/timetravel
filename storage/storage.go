package storage

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/temelpa/timetravel/entity"
)

const databaseFile = "sqlite-database.db"
const timestampFormat = "2006-01-02T15:04:05Z"

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
		"id" integer NOT NULL,
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

func (s *Storage) GetRecordsByID(id int) ([]*entity.Record, error) {
	log.Println("Getting record...")
	getRecordSQL := `SELECT id, data, created_at, deleted_at FROM records WHERE id = ?`

	statement, err := s.db.Prepare(getRecordSQL)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	var records []*entity.Record
	rows, err := statement.Query(id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		record := &entity.Record{}
		var data string
		var nullableTime sql.NullTime
		var timestampStr string
		err = rows.Scan(&record.ID, &data, &timestampStr, &nullableTime)
		log.Println(err)
		if err != nil {
			return nil, err
		}
		log.Println(timestampStr)
		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		log.Println(err)
		if err != nil {
			return nil, err
		}
		log.Println(timestamp)
		log.Println(data)
		err = json.Unmarshal([]byte(data), &record.Data)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		record.CreatedAt = timestamp
		records = append(records, record)
	}
	return records, nil
}

func (s *Storage) GetLastestRecordByID(id int) (*entity.Record, error) {
	log.Println("Getting latest record...")
	getRecordSQL := `SELECT id, data, created_at, deleted_at FROM records WHERE id = ? ORDER BY created_at DESC LIMIT 1`

	statement, err := s.db.Prepare(getRecordSQL)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	row := statement.QueryRow(id)
	record := &entity.Record{}
	var data string
	var nullableTime sql.NullTime
	var timestampStr string
	err = row.Scan(&record.ID, &data, &timestampStr, &nullableTime)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	log.Println(timestampStr)
	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	log.Println(err)
	if err != nil {
		return nil, err
	}
	log.Println(timestamp)
	log.Println(data)
	err = json.Unmarshal([]byte(data), &record.Data)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	record.CreatedAt = timestamp

	return record, nil
}

func (s *Storage) GetRecordsByIDBetweenTimestamp(id int, startTime, endTime time.Time) ([]*entity.Record, error) {
	log.Println("Getting record...")
	getRecordSQL := `SELECT id, data, created_at, deleted_at FROM records WHERE id = ? AND created_at BETWEEN ? AND ? ORDER BY created_at DESC`

	statement, err := s.db.Prepare(getRecordSQL)
	if err != nil {
		return nil, err
	}
	defer statement.Close()
	var records []*entity.Record
	rows, err := statement.Query(id, startTime, endTime)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		record := &entity.Record{}
		var data string
		var nullableTime sql.NullTime
		var timestampStr string
		err = rows.Scan(&record.ID, &data, &timestampStr, &nullableTime)
		log.Println(err)
		if err != nil {
			return nil, err
		}
		log.Println(timestampStr)
		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		log.Println(err)
		if err != nil {
			return nil, err
		}
		log.Println(timestamp)
		log.Println(data)
		err = json.Unmarshal([]byte(data), &record.Data)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		record.CreatedAt = timestamp
		records = append(records, record)
	}
	return records, nil
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
