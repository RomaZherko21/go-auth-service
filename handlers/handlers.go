package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"exampleApi/entities/user"
)

func Handlers(r *gin.Engine, db *sql.DB) {

	InitMiddlewares(r, db)

	r.POST("/signIn", user.SignIn)
	r.POST("/signUp", user.SignUp)

	r.POST("/users", user.CreateUser)

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
