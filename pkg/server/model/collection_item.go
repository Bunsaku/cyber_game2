package model

import (
	"20dojo-online/pkg/db"
	"database/sql"
)

//CollectionItem テーブルのデータ
type CollectionItem struct {
	ID     string
	Name   string
	Rarity int
}

//GetCollectionItem コレクションアイテムテーブルのデータの取得
func GetCollectionItem(tx *sql.Tx, IDs []string) ([]CollectionItem, error) {
	rows, err := tx.Query("SELECT * FROM collection_item ;")
	if err != nil {
		return nil, err
	}
	i := CollectionItem{}
	var items []CollectionItem
	for rows.Next() {
		err = rows.Scan(&i.ID, &i.Name, &i.Rarity)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	collectionMap := make(map[string]CollectionItem, len(items))

	for _, v := range items {
		collectionMap[v.ID] = v
	}

	var results []CollectionItem
	for _, id := range IDs {
		if v, ok := collectionMap[id]; ok {
			results = append(results, v)
		}
	}
	return results, nil

}

//GetAllCollections すべてのコレクションアイテムデータ取得
func GetAllCollections() ([]CollectionItem, error) {
	rows, err := db.Conn.Query("SELECT * FROM collection_item ORDER BY id;")
	if err != nil {
		return nil, err
	}
	i := CollectionItem{}
	var items []CollectionItem
	for rows.Next() {
		err = rows.Scan(&i.ID, &i.Name, &i.Rarity)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil

}
