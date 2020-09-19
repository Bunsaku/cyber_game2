package model

import (
	"database/sql"
)

//GachaProbability is gacha_probabilityのデータ
type GachaProbability struct {
	CollectionItemID string
	Ratio            int
}

//GetAllGachaProbability  ガチャ確立取得"
func GetAllGachaProbability(tx *sql.Tx) ([]GachaProbability, error) {
	rows, err := tx.Query("SELECT * FROM gacha_probability;")
	if err != nil {
		return nil, err
	}
	i := GachaProbability{}
	var items []GachaProbability
	for rows.Next() {
		err = rows.Scan(&i.CollectionItemID, &i.Ratio)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
