package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"exampleApi/entities/user"
)

func Handlers(r *gin.Engine, db *sql.DB) {

	InitMiddlewares(r, db)

	r.POST("/users/", user.CreateUser)
	// r.GET("/task/", server.getAllTasksHandler)
	// r.DELETE("/task/", server.deleteAllTasksHandler)
	// r.GET("/task/:id", server.getTaskHandler)
	// r.DELETE("/task/:id", server.deleteTaskHandler)
	// r.GET("/tag/:tag", server.tagHandler)
	// r.GET("/due/:year/:month/:day", server.dueHandler)
}
