package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"ticketon-auth-service/api/model"
	evtRepo "ticketon-auth-service/api/repository/event"
	"time"
)

func GetEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ApiError{Message: "id is required"})
		return
	}

	evtFound, err := evtRepo.DB.First(eventID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.AbortWithStatusJSON(http.StatusNotFound, model.ApiError{Message: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, evtFound)
}

func CreateEvent(c *gin.Context) {
	var evtReq model.CreateEventRequest

	if err := c.ShouldBindJSON(&evtReq); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ApiError{Message: err.Error()})
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ApiError{Message: "user_id missing in jwt"})
		return
	}

	evtToCreate := model.EventBasic{
		Model: gorm.Model{
			CreatedAt: time.Now(),
		},
		Name:      evtReq.Name,
		StartDate: evtReq.StartDate,
		EndDate:   evtReq.EndDate,
		Capacity:  evtReq.Capacity,
		Location: model.LocationEvent{
			Latitude:     evtReq.Location.Latitude,
			Longitude:    evtReq.Location.Longitude,
			LocationName: evtReq.Location.LocationName,
		},
		UserID: uint(userID.(int)),
	}

	record := evtRepo.DB.Create(evtToCreate)
	if record.Error != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(record.Error, &mysqlErr) {
			// Check if it's a duplicate entry error (error code 1062)
			if mysqlErr.Number == 1062 {
				c.AbortWithStatusJSON(http.StatusConflict, model.ApiError{Message: "event already exists"})
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: record.Error.Error()})
		return
	}

	if evtCreated, ok := record.Statement.Model.(*model.EventBasic); ok {
		evtToCreate.ID = evtCreated.ID
	}

	c.JSON(http.StatusCreated, evtToCreate)
}

func UpdateEvent(c *gin.Context) {
	// Get the user ID from the URL path
	evtID := c.Param("id")

	// Check if the user exists
	if _, err := evtRepo.DB.First(evtID); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, model.ApiError{Message: "Event not found. " + err.Error()})
		return
	}

	// Bind the request body to the CreateUserRequest struct
	var updatedEvtData model.CreateEventRequest
	if err := c.ShouldBindJSON(&updatedEvtData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ApiError{Message: err.Error()})
		return
	}

	// Update fields that are allowed to be updated
	existingEvt := model.EventBasic{
		Name:      updatedEvtData.Name,
		StartDate: updatedEvtData.StartDate,
		EndDate:   updatedEvtData.EndDate,
		Capacity:  updatedEvtData.Capacity,
		Location: model.LocationEvent{
			Latitude:     updatedEvtData.Location.Latitude,
			Longitude:    updatedEvtData.Location.Longitude,
			LocationName: updatedEvtData.Location.LocationName,
		},
	}

	// Save the updated user to the database
	if err := evtRepo.DB.Update(existingEvt).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}

	// Return the updated user data in the response
	c.JSON(http.StatusOK, existingEvt)
}

func DeleteEvent(c *gin.Context) {
	eventID := c.Param("id")
	if eventID == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, model.ApiError{Message: "id is required"})
		return
	}

	evtFound, err := evtRepo.DB.First(eventID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.AbortWithStatusJSON(http.StatusNotFound, model.ApiError{Message: err.Error()})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: err.Error()})
		return
	}

	gormResp := evtRepo.DB.Delete(*evtFound)
	if gormResp.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, model.ApiError{Message: gormResp.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
