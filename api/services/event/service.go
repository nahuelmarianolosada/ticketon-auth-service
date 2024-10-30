package event

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"ticketon-auth-service/api/model"
	evtRepo "ticketon-auth-service/api/repository/event"
	"time"
)

func CreateEvent(ctx context.Context, bodyReq model.CreateEventRequest, userID any) (*model.EventBasic, error) {
	evtToCreate := model.EventBasic{
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
		Name:      bodyReq.Name,
		StartDate: bodyReq.StartDate,
		EndDate:   bodyReq.EndDate,
		Capacity:  bodyReq.Capacity,
		Location: model.LocationEvent{
			Latitude:     bodyReq.Location.Latitude,
			Longitude:    bodyReq.Location.Longitude,
			LocationName: bodyReq.Location.LocationName,
		},
		UserID: uint(userID.(int)),
	}

	record := evtRepo.DB.Create(evtToCreate)
	if record.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(record.Error, &mysqlErr) {
			// Check if it's a duplicate entry error (error code 1062)
			if mysqlErr.Number == 1062 {
				return nil, model.ApiError{Message: "event already exists", Err: record.Error}
			}
		}
		return nil, model.ApiError{Message: record.Error.Error()}
	}

	if evtCreated, ok := record.Statement.Model.(*model.EventBasic); ok {
		evtToCreate.ID = evtCreated.ID
	}

	return &evtToCreate, nil
}
