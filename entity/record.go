package entity

import (
	"encoding/json"
	"time"
)

type Record struct {
	ID        int               `json:"id"`
	Data      map[string]string `json:"data"`
	CreatedAt time.Time         `json:"created_at"`
	DeletedAt time.Time         `json:"deleted_at"`
}

func (d *Record) Copy() Record {
	values := d.Data

	newMap := map[string]string{}
	for key, value := range values {
		newMap[key] = value
	}

	return Record{
		ID:        d.ID,
		Data:      newMap,
		CreatedAt: d.CreatedAt,
		DeletedAt: d.DeletedAt,
	}
}

func (d *Record) ToString() string {
	return string(d.ToJSON())
}

func (d *Record) ToJSON() []byte {
	data, _ := json.Marshal(d)
	return data
}
