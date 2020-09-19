package model

import (
	"20dojo-online/pkg/db"
	"database/sql"
	"strconv"
	"strings"
)

//UserCollectionItem テーブルのデータ
type UserCollectionItem struct {
	UserID           string
	CollectionItemID string
}

//HasItem IsNewがtrueかfalseか判別し取得
func HasItem(tx *sql.Tx, IDs []string, userID string) ([]string, error) {
	rows, err := tx.Query("SELECT * FROM user_collection_item WHERE user_id = ?;", userID)
	if err != nil {
		return nil, err
	}
	i := UserCollectionItem{}
	var AllItems []UserCollectionItem
	for rows.Next() {
		err = rows.Scan(&i.UserID, &i.CollectionItemID)
		if err != nil {
			return nil, err
		}
		AllItems = append(AllItems, i)
	}

	UserCollectionMap := make(map[string]UserCollectionItem, len(AllItems))

	for _, v := range AllItems {
		UserCollectionMap[v.CollectionItemID] = v
	}

	var NewItemIDs []string
	for _, id := range IDs {
		flag := true
		if _, ok := UserCollectionMap[id]; ok {
			flag = false
		}
		if flag {
			NewItemIDs = append(NewItemIDs, id)
		}
	}

	return NewItemIDs, nil
}

//UpdatetUserCollectionItems バルクインサート
func UpdatetUserCollectionItems(tx *sql.Tx, IDs []string, userID string) error {
	var inserData strings.Builder
	inserData.WriteString("INSERT INTO user_collection_item(user_id, collection_item_id) VALUES ")

	for i, id := range IDs {
		if i != 0 {
			inserData.WriteString(", ")
		}
		inserData.WriteString("(")
		inserData.WriteString(strconv.Quote(userID))
		inserData.WriteString(", ")
		inserData.WriteString(strconv.Quote(id))
		inserData.WriteString(")")
	}

	stmt, err := tx.Prepare(inserData.String())
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	return err
}

//GetUserCollections ユーザーコレクションアイテムのＩＤを取得
func GetUserCollections(userID string) ([]string, error) {
	rows, err := db.Conn.Query("SELECT collection_item_id FROM user_collection_item WHERE user_id = ? ORDER BY collection_item_id;", userID)
	if err != nil {
		return nil, err
	}
	i := UserCollectionItem{}
	var IDs []string
	for rows.Next() {
		err = rows.Scan(&i.CollectionItemID)
		if err != nil {
			return nil, err
		}
		IDs = append(IDs, i.CollectionItemID)
	}
	return IDs, nil
}
