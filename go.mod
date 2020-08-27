module penpal-api

go 1.13

require (
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/pkg/errors v0.9.1
	go.mongodb.org/mongo-driver v1.4.0
	smtp v0.0.0-00010101000000-000000000000
)

replace smtp => ./smtp
