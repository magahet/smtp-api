package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var server, user, pass string
var port int

func main() {
	flag.IntVar(&port, "port", 587, "SMTP Server")
	flag.StringVar(&user, "user", "", "user")
	flag.StringVar(&pass, "pass", "", "pass")

	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatal("You must specify an smtp server")
	}
	server = os.Args[1]

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	router.Use(cors.Default())

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

	sender := NewSender(server, port)
	msg := &Message{req.From, req.Subject, req.Text, req.To}
	if err := sender.Send(msg, user, pass); err != nil {
		wrappedError := errors.Wrap(err, "Could not send email")
		c.JSON(http.StatusInternalServerError,
			gin.H{"status": "error", "message": wrappedError.Error()})
		return
	}

	log.Println("Message sent")

	c.JSON(http.StatusOK,
		gin.H{"status": "success", "data": nil})
}
