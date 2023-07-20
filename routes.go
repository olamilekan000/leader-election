package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddRoutes(router *gin.Engine) {

	router.GET("/", Pong)
	router.GET("/points", Points)
}

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"code":    200,
		"data":    map[string]string{},
	})
}

func Points(c *gin.Context) {
	rVal, err := rcl.Get(context.Background(), rKey)
	if err != nil {

		fmt.Println("err:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"code":    200,
		"data": map[string]string{
			"count": *rVal,
		},
	})
}
