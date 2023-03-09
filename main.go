package main

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func main() {
	router := gin.Default()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.GET("/", func(c *gin.Context) {
		c.JSON(
			http.StatusOK,
			gin.H{
				"code":     http.StatusOK,
				"datetime": time.Now(),
			},
		)
	})

	router.POST("create-tag", CreateTag)
	router.GET("get-all-tags", GetTags)

	router.POST("ingest-data", CreateTagData)
	router.POST("get-tag-data", GetTagData)

	router.Run("localhost:9090").Error()
}
