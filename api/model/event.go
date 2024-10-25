package model

import (
	"gorm.io/gorm"
	"time"
)

type EventBasic struct {
	gorm.Model
	Name      string        `json:"name"`
	StartDate time.Time     `json:"start_date"`
	EndDate   *time.Time    `json:"end_date"`
	Capacity  uint          `json:"capacity"`
	Location  LocationEvent `json:"location" gorm:"embedded"`
}

type LocationEvent struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
}

func (e EventBasic) TableName() string {
	return "event"
}

type CreateEventRequest struct {
	Name      string        `json:"name"`
	StartDate time.Time     `json:"start_date"`
	EndDate   *time.Time    `json:"end_date"`
	Capacity  uint          `json:"capacity"`
	Location  LocationEvent `json:"location" `
}
