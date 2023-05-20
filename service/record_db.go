package service

import (
	"context"
	"encoding/json"

	"github.com/temelpa/timetravel/entity"
	"github.com/temelpa/timetravel/storage"
)

type DatabaseService struct {
	storage *storage.Storage
}

func NewDatabaseService(storage *storage.Storage) DatabaseService {
	return DatabaseService{storage: storage}
}

func (s *DatabaseService) GetRecord(ctx context.Context, id int) (entity.Record, error) {
	record, err := s.storage.GetRecordByID(id)
	if err != nil {
		return entity.Record{}, err
	}
	if record.ID == 0 {
		return entity.Record{}, ErrRecordDoesNotExist
	}

	newRecord := record.Copy() // copy is necessary so modifations to the record don't change the stored record
	return newRecord, nil
}

func (s *DatabaseService) CreateRecord(ctx context.Context, record entity.Record) error {
	id := record.ID
	if id <= 0 {
		return ErrRecordIDInvalid
	}

	dataBytes, err := json.Marshal(record.Data)
	if err != nil {
		return err
	}

	existingRecord, err := s.storage.InsertRecord(id, string(dataBytes))
	if err != nil {
		return err
	}
	if existingRecord != 0 {
		return ErrRecordAlreadyExists
	}
	return nil
}

func (s *DatabaseService) UpdateRecord(ctx context.Context, id int, updates map[string]string) (entity.Record, error) {
	updateBytes, err := json.Marshal(updates)
	if err != nil {
		return entity.Record{}, err
	}

	entry, err := s.storage.UpdateRecord(id, string(updateBytes))
	if err != nil {
		return entity.Record{}, err
	}
	if entry == 0 {
		return entity.Record{}, ErrRecordDoesNotExist
	}

	updatedEntry, err := s.storage.GetRecordByID(id)
	if err != nil {
		return entity.Record{}, err
	}

	return updatedEntry.Copy(), nil
}
