package main

import (
	"database/sql"
)

type (
	Tag struct {
		PartitionKey string      `json:"partitionKey"`
		TagName      string      `json:"tagName"`
		MetaData     []*MetaData `json:"metaData"`
		SeriesId     int64       `json:"seriesId"`
	}

	MetaData struct {
		Key   sql.NullString `json:"key"`
		Value sql.NullString `json:"value"`
	}

	View struct {
		TagName  string      `json:"tagName"`
		MetaData []*MetaData `json:"metaData"`
	}

	Key struct {
		PartitionKey string `json:"partitionKey"`
		TagName      string `json:"tagName"`
	}

	Series struct {
		SeriesId int64    `json:"seriesId"`
		Dt       []*int64 `json:"dt"`
		Val      []*any   `json:"Val"`
	}

	DataIngestByTagName struct {
		TagName string `json:"tagName"`
		Dt      int64  `json:"dt"`
		Val     any    `json:"Val"`
	}

	DataIngestById struct {
		SeriesId int64 `json:"seriesId"`
		Dt       int64 `json:"dt"`
		Val      any   `json:"Val"`
	}

	DataRequest struct {
		TagName string `json:"tagName"`
		Start   int64  `json:"start"`
		End     int64  `json:"end"`
	}
)
