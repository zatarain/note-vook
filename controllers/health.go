package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheck(context *gin.Context) {
	context.String(http.StatusOK, "OK, go!")
}
