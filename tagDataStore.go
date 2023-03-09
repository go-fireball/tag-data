package main

import (
	"encoding/json"
)

func (t *PostgresDataStore) GetTagDataByRange(seriesId int64, start int64, end int64) (*Series, error) {
	rows, err := t.Database.Query(`SELECT "Dt", "Value" FROM public."SeriesData" WHERE "SeriesId" = $1 AND "Dt" >= $2 AND "Dt" <= $3 order by "Dt"`, seriesId, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var series = Series{}
	series.SeriesId = seriesId

	for rows.Next() {
		var dt int64
		var val []byte
		if err := rows.Scan(&dt, &val); err != nil {
			return &series, err
		}
		if err != nil {
			return nil, err
		}

		var storedVal any
		err = json.Unmarshal(val, &storedVal)

		series.Dt = append(series.Dt, &dt)
		series.Val = append(series.Val, &storedVal)
	}
	if err = rows.Err(); err != nil {
		return &series, err
	}
	return &series, nil
}

func (t *PostgresDataStore) CreateTagData(records []*DataIngestById) error {
	tx, err := t.Database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO public."SeriesData" ("SeriesId", "Dt", "Value") VALUES ($1, $2, $3)  ON CONFLICT ("SeriesId", "Dt") DO UPDATE SET "Value" = excluded."Value"`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		valBytes, err := json.Marshal(record.Val)
		if err != nil {
			return err
		}
		val := string(valBytes)

		//result, err := tx.Exec(`INSERT INTO public."SeriesData" ("SeriesId", "Dt", "Value") VALUES ($1, $2, $3)
		//									ON CONFLICT ("SeriesId", "Dt") DO UPDATE SET "Value" = excluded."Value"`, 	record.SeriesId, record.Dt, val)

		result, err := stmt.Exec(record.SeriesId, record.Dt, val)

		if err != nil {
			return err
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			panic("Can't Update or Insert the row")
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
