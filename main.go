package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bluemon0919/radiko-linker/api"

	"github.com/gin-gonic/gin"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	api := api.NewWebAPI()
	flag.Parse()
	router := gin.Default()
	v1 := router.Group("")
	{
		// ヘルプを表示する
		v1.GET("/", func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, "static/help.html")
		})
		// リンクを生成して返すページを表示する
		v1.GET("/make-link", func(c *gin.Context) {
			http.ServeFile(c.Writer, c.Request, "static/make-link.html")
		})
		// リンクを受け取ってradikoにジャンプさせるAPI
		v1.GET("/jump/:station/:dayOfWeek/:startTime", api.JumpAPI)
	}
	router.Run(":8080")

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
