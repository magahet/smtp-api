package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

var Server string
var Port int = 25

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You must specify an smtp server")
	}
	Server = os.Args[1]
	if len(os.Args) == 3 {
		var err error
		Port, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("SMTP port not valid")
		}
	}
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()

	router.POST("/message", sendMessage)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	router.Run()
	// router.Run(":3000") for a hard coded port
}

type jsonError struct {
	err error
}

type EmailRequest struct {
	From    string   `form:"from" json:"from" binding:"required"`
	Subject string   `form:"subject" json:"subject" binding:"required"`
	Text    string   `form:"text" json:"text" binding:"required"`
	To      []string `form:"to" json:"to" binding:"required"`
}

func (err jsonError) Error() string {
	return errors.Wrap(err, "request json missing fields or invalid").Error()
}

func sendMessage(c *gin.Context) {
	var req EmailRequest
	if err := c.ShouldBind(&req); err != nil {
		// wrappedError := errors.Wrap(err, "Could not parse request or missing fields")
		c.JSON(http.StatusBadRequest,
			gin.H{"status": "fail", "data": err})
		return
	}

	sender := NewSender(Server, Port)
	msg := &Message{req.From, req.Subject, req.Text, req.To}
	if err := sender.Send(msg); err != nil {
		wrappedError := errors.Wrap(err, "Could not send email")
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": "error", "message": wrappedError.Error()})
		return
	}

	log.Println("Message sent")

	c.JSON(http.StatusOK,
		gin.H{"status": "success", "data": nil})
}
