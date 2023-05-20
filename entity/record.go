package entity

import "encoding/json"

type Record struct {
	ID        int               `json:"id"`
	Data      map[string]string `json:"data"`
	CreatedAt string            `json:"created_at"`
	DeletedAt string            `json:"deleted_at"`
}

func (d *Record) Copy() Record {
	values := d.Data

	newMap := map[string]string{}
	for key, value := range values {
		newMap[key] = value
	}

	return Record{
		ID:   d.ID,
		Data: newMap,
	}
}

func (d *Record) ToString() string {
	return string(d.ToJSON())
}

func (d *Record) ToJSON() []byte {
	data, _ := json.Marshal(d)
	return data
}
