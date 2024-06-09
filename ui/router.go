package ui

import (
	"lyrics-quiz/db"

	"github.com/gin-gonic/gin"
)

func DBContext() gin.HandlerFunc {
	dbConn := db.ConnectDB()
	return func(c *gin.Context) {
		c.Set("db", dbConn)
	}
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(DBContext())

	// ping
	r.GET("/ping/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}
