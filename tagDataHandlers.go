package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateTagData(context *gin.Context) {
	partitionKey := context.Request.Header.Get("PartitionKey")
	if len(partitionKey) != 36 {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Expecting PartitionKey in the Header",
		})
		return
	}
	var dataIngest []DataIngestByTagName
	if err := context.Bind(&dataIngest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tagMap = make(map[string]*Tag)

	for index := range dataIngest {
		tagName := dataIngest[index].TagName
		if _, ok := tagMap[tagName]; !ok {
			tag, err := env.PgDb.GetTagByTagKey(&Key{
				TagName:      dataIngest[index].TagName,
				PartitionKey: partitionKey,
			})
			if err != nil {
				context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if tag == nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": "Tag Not Found"})
				return
			}
			tagMap[tagName] = tag
		}
	}

	var allDataIngestById []*DataIngestById

	for index := range dataIngest {
		tagName := dataIngest[index].TagName
		var dataIngestById = &DataIngestById{
			SeriesId: tagMap[tagName].SeriesId,
			Dt:       dataIngest[index].Dt,
			Val:      dataIngest[index].Val,
		}
		allDataIngestById = append(allDataIngestById, dataIngestById)
	}
	err := env.PgDb.CreateTagData(allDataIngestById)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func GetTagData(context *gin.Context) {
	partitionKey := context.Request.Header.Get("PartitionKey")
	if len(partitionKey) != 36 {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Expecting PartitionKey in the Header",
		})
		return
	}

	var dataReqeust DataRequest

	if err := context.Bind(&dataReqeust); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tag, err := env.PgDb.GetTagByTagKey(&Key{
		TagName:      dataReqeust.TagName,
		PartitionKey: partitionKey,
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tag == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Tag Not Found"})
	}

	series, err := env.PgDb.GetTagDataByRange(tag.SeriesId, dataReqeust.Start, dataReqeust.End)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, series)

}
