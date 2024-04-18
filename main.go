package main

import (
	"github.com/gin-gonic/gin"
	"github.com/judeosbert/bus_tracker_bot/telegram"
)

func main(){
	r := gin.Default()
	r.POST("/",func(ctx *gin.Context) {
		telegram.HandleIncomingMessage(ctx)
	})
	r.Run()
}
