package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var env = GetEnvironmentConfig()

func CreateTag(context *gin.Context) {

	partitionKey := context.Request.Header.Get("PartitionKey")
	var tagView View

	if err := context.Bind(&tagView); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tagKey := &Key{
		TagName:      tagView.TagName,
		PartitionKey: partitionKey,
	}

	existingTag, err := env.PgDb.GetTagByTagKey(tagKey)

	if err != nil {
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if existingTag != nil {
		context.JSON(http.StatusConflict, gin.H{"error": "Tag Already Exists", "Tag": existingTag})
		return
	}

	result, err := env.PgDb.Create(Tag{
		TagName:      tagView.TagName,
		PartitionKey: partitionKey,
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	context.JSON(http.StatusOK, result)
}

func GetTags(context *gin.Context) {
	partitionKey := context.Request.Header.Get("PartitionKey")
	if len(partitionKey) != 36 {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Expecting PartitionKey in the Header",
		})
		return
	}
	tags, err := env.PgDb.GetTagsByPartitionKey(partitionKey)

	if err != nil {
		context.JSON(http.StatusInternalServerError, "Internal Server Error")
	}
	context.JSON(http.StatusOK, tags)
}
