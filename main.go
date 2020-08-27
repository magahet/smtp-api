package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var server, user, pass string
var servicePort, smtpPort int

func handleGetEvents(c *gin.Context) {
	var loadedEvents, err = GetAllEvents()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": loadedEvents})
}

func handleGetEvent(c *gin.Context) {
	var event Event
	if err := c.BindUri(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	var loadedEvent, err = GetEventByID(event.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, loadedEvent)
}

func handleCreateEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	id, err := Create(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	if err := sendEmail(&event); err != nil {
		log.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func handleUpdateEvent(c *gin.Context) {
	var event Event
	if err := c.ShouldBindJSON(&event); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	savedEvent, err := Update(&event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": savedEvent})
}

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()

	// Middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors())

	v1 := r.Group("/v1")
	{
		events := v1.Group("/events")
		{
			events.GET("/events/:id", handleGetEvent)
			events.GET("/events/", handleGetEvents)
			events.PUT("/events/", handleCreateEvent)
			events.POST("/events/", handleUpdateEvent)
		}
	}

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
}
