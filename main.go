package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/judeosbert/bus_tracker_bot/telegram"
)

func main() {
	r := gin.Default()
	r.POST("/", func(ctx *gin.Context) {
		body, _ := io.ReadAll(ctx.Request.Body)
		fmt.Sprintf("incoming body %s",string(body))

		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
		telegram.HandleIncomingMessage(ctx)
	})
	r.Run()
}
