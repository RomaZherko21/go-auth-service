package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"exampleApi/entities/auth"
)

func Handlers(r *gin.Engine, db *sql.DB, redisDb *redis.Client) {

	InitMiddlewares(r, db, redisDb)

	r.POST("/auth/signIn", auth.SignIn)
	r.POST("/auth/signUp", auth.SignUp)
	r.DELETE("/auth/signOut", authMiddleware, auth.SignOut)
	r.DELETE("/auth/signOutAll", authMiddleware, auth.SignOutFromAllDevices)
	r.POST("/auth/refresh", auth.Refresh)

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
