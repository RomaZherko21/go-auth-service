package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"

	"exampleApi/entities/auth"
)

func Handlers(r *gin.Engine, db *sql.DB, redisDb *redis.Client) {

	InitMiddlewares(r, db, redisDb)

	r.POST("/signUp", auth.SignUp)
	r.POST("/signIn", auth.SignIn)
	r.POST("/signOut", authMiddleware, auth.SignOut)

	r.GET("/private", authMiddleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "This is a private route"})
	})

	// r.GET("/task/", server.getAllTasksHandler)
	// r.DELETE("/task/", server.deleteAllTasksHandler)
	// r.GET("/task/:id", server.getTaskHandler)
	// r.DELETE("/task/:id", server.deleteTaskHandler)
	// r.GET("/tag/:tag", server.tagHandler)
	// r.GET("/due/:year/:month/:day", server.dueHandler)
}
