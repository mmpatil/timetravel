package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// GET /records/{id}
// GetRecord retrieves the record.
func (a *API) GetRecords(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	record, err := a.records.GetRecord(
		ctx,
		int(idNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, record, http.StatusOK)
	logError(err)
}

// GET /records/{id}
// GetRecordV2 retrieves the record.
func (a *API) GetRecordsV2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}
	record, err := a.recordsV2.GetAllRecordsByID(ctx, int(idNumber))
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, record, http.StatusOK)
	logError(err)
}

// GET /record/{id}
// GetRecord retrieves the latest record.
func (a *API) GetLastestRecordV2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	record, err := a.recordsV2.GetLastestRecordByID(
		ctx,
		int(idNumber),
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}
	err = writeJSON(w, record, http.StatusOK)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}
}

// GET /records/{id}/{start}/{end}
// GetRecordsBetweenTimestamp retrieves the record.
func (a *API) GetRecordsBetweenTimestampV2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	startTimeString := mux.Vars(r)["start"]
	endTimeString := mux.Vars(r)["end"]

	idNumber, err := strconv.ParseInt(id, 10, 32)

	if err != nil || idNumber <= 0 {
		err := writeError(w, "invalid id; id must be a positive number", http.StatusBadRequest)
		logError(err)
		return
	}

	var startTime time.Time
	if startTimeString != "" {
		var err error
		startTime, err = time.Parse(time.RFC3339, startTimeString)
		if err != nil {
			err := writeError(w, "invalid time format", http.StatusBadRequest)
			logError(err)
			return
		}
	}

	var endTime time.Time
	if startTimeString != "" {
		var err error
		endTime, err = time.Parse(time.RFC3339, endTimeString)
		if err != nil {
			err := writeError(w, "invalid time format", http.StatusBadRequest)
			logError(err)
			return
		}
	}
	log.Printf("start time: %v", startTime)
	log.Printf("end time: %v", endTime)
	record, err := a.recordsV2.GetRecordsByIDBetweenTimestamp(
		ctx,
		int(idNumber),
		startTime,
		endTime,
	)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}

	err = writeJSON(w, record, http.StatusOK)
	if err != nil {
		err := writeError(w, fmt.Sprintf("record of id %v does not exist", idNumber), http.StatusBadRequest)
		logError(err)
		return
	}
}
