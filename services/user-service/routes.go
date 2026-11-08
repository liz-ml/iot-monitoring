package main

import (
    "github.com/gin-gonic/gin"

    "github.com/liz-ml/iot-monitoring/services/user-service/handlers"
)

func SetupRouter() *gin.Engine {
    router := gin.Default()

    router.GET("/health", handlers.HealthCheck)

    return router
}