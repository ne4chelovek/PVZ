package model

import (
	"time"
)

type Reception struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    string    `json:"pvzId"`
	Status   string    `json:"status"` // in_progress, close
}

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"` // электроника, одежда, обувь
	ReceptionID string    `json:"receptionId"`
}

type PVZWithReceptions struct {
	PVZ        *PVZ                     `json:"pvz"`
	Receptions []*ReceptionWithProducts `json:"receptions"`
}

type ReceptionWithProducts struct {
	Reception *Reception `json:"reception"`
	Products  []*Product `json:"products"`
}
