package main

import (
	_ "github.com/lib/pq"
)

func (t *PostgresDataStore) GetTagsByPartitionKey(partitionKey string) ([]*Tag, error) {
	rows, err := t.Database.Query(`SELECT T."PartitionKey",
												T."TagName",
												T."SeriesId",
												T0."MetaDataKey",
												T0."MetaDataValue"
											FROM "Tag" AS T
											LEFT JOIN "TagMetaData" AS T0 ON (T."SeriesId" = T0."SeriesId")
											WHERE T."PartitionKey" = $1
											ORDER BY T."SeriesId"`, partitionKey)
	if err != nil {
		return nil, err
	}
	var tagMap = make(map[string]*Tag)

	defer rows.Close()

	for rows.Next() {
		var tag Tag
		var tagMetaData MetaData

		if err := rows.Scan(&tag.PartitionKey, &tag.TagName, &tag.SeriesId, &tagMetaData.Key, &tagMetaData.Value); err != nil {
			return nil, err
		}
		if val, ok := tagMap[tag.TagName]; ok {
			val.MetaData = append(val.MetaData, &tagMetaData)
		} else {
			tag.MetaData = []*MetaData{&tagMetaData}
			tagMap[tag.TagName] = &tag
		}
	}

	values := make([]*Tag, 0, len(tagMap))

	for _, v := range tagMap {
		values = append(values, v)
	}

	if err = rows.Err(); err != nil {
		return values, err
	}

	return values, nil
}

func (t *PostgresDataStore) Create(tag Tag) (*Tag, error) {
	// Fetch the Partition key
	_, err := t.Database.Exec(`INSERT INTO PUBLIC."Tag"("PartitionKey","TagName") VALUES ($1, $2)`, tag.PartitionKey, tag.TagName)
	if err != nil {
		return nil, err
	}
	returnTag, err := t.GetTagByTagKey(&Key{
		TagName:      tag.TagName,
		PartitionKey: tag.PartitionKey,
	})
	if err != nil {
		return nil, err
	} else {
		return returnTag, err
	}

}

func (t *PostgresDataStore) GetTagByTagKey(tagKey *Key) (*Tag, error) {
	rows, err := t.Database.Query(`SELECT T."PartitionKey",
									T."TagName",
									T."SeriesId",
									T0."MetaDataKey",
									T0."MetaDataValue"
								FROM "Tag" AS T
								LEFT JOIN "TagMetaData" AS T0 ON (T."SeriesId" = T0."SeriesId")
								WHERE T."PartitionKey" = $1 AND T."TagName" = $2
								`, tagKey.PartitionKey, tagKey.TagName)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var tag *Tag
	for rows.Next() {
		tag = &Tag{}
		var tagMetaData MetaData

		if err := rows.Scan(&tag.PartitionKey, &tag.TagName, &tag.SeriesId, &tagMetaData.Key, &tagMetaData.Value); err != nil {
			return nil, err
		}
		if tagMetaData.Key.Valid {
			if tag.MetaData == nil {
				tag.MetaData = []*MetaData{&tagMetaData}
			} else {
				tag.MetaData = append(tag.MetaData, &tagMetaData)
			}
		}
	}
	if err = rows.Err(); err != nil {
		return tag, err
	}
	return tag, nil
}
