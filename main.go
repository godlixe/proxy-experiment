package main

import (
	"log/slog"
	"net/http"
	"proxy-experiment/internal"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DataResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Response struct {
	Message   string `json:"message"`
	IPAddress string `json:"ip_address"`
	Username  string `json:"username"`
}

const PROXY_REDIS_KEY_PREFIX = "proxy-experiment"
const PROXY_REDIS_COUNT_KEY_PREFIX = "proxy-experiment-count"

func main() {
	db := internal.InitDB()
	server := gin.Default()

	slog.Info("Server running.")

	server.GET("/", func(ctx *gin.Context) {

		response := Response{
			Message:   "proxy experiment, user data:",
			IPAddress: ctx.RemoteIP(),
			Username:  ctx.GetHeader("x-username"),
		}
		ctx.JSON(http.StatusOK, response)
	})

	server.GET("/hit", func(ctx *gin.Context) {
		username := ctx.GetHeader("x-username")
		if username == "" {
			slog.Error("Empty username.")
			ctx.AbortWithStatusJSON(http.StatusBadRequest, "Error")
			return
		}

		ipAddress := ctx.RemoteIP()

		var user internal.User

		tx := db.Where("username = ?", username).First(&user)
		if tx.Error != nil && tx.Error != gorm.ErrRecordNotFound {
			slog.Error(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Error")
			return
		}

		// create user if user doesn't exist
		if user.Username == "" {
			user := internal.User{
				Username: username,
				Count:    0,
			}

			// Add user to the database
			tx := db.Create(user)
			if tx.Error != nil {
				slog.Error(tx.Error.Error())
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Error")
				return
			}
		}

		// Increment hit count
		tx = db.Model(&internal.User{}).Where("username = ?", username).Update("count", gorm.Expr("count + ?", 1))
		if tx.Error != nil {
			slog.Error(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Error")
			return
		}

		accessLog := internal.AccessLog{
			Username:  username,
			IPAddress: ipAddress,
		}

		// Insert to access logs
		tx = db.Create(&accessLog)
		if tx.Error != nil {
			slog.Error(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Error")
		}

		response := Response{
			Message:   "Hit success.",
			IPAddress: ipAddress,
			Username:  username,
		}

		ctx.JSON(http.StatusOK, response)
	})

	server.GET("/list", func(ctx *gin.Context) {
		var users []internal.User
		tx := db.Find(&users)
		if tx.Error != nil {
			slog.Error(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Error")
			return
		}

		response := DataResponse{
			Message: "hit success.",
			Data:    users,
		}

		ctx.JSON(http.StatusOK, response)
	})

	server.GET("/list/:username", func(ctx *gin.Context) {
		username := ctx.Param("username")

		var accessLogs []internal.AccessLog
		tx := db.Where("username = ?", username).Find(&accessLogs)
		if tx.Error != nil {
			slog.Error(tx.Error.Error())
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, "Error")
		}

		response := DataResponse{
			Message: "list detail success.",
			Data:    accessLogs,
		}

		ctx.JSON(http.StatusOK, response)
	})

	server.Run(":8084")
}
