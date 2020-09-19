package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

//HandleCollectionList 'ランキングの取得'
func HandleCollectionList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println(errors.New("userID is empty"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// ユーザデータの取得処理と存在チェックを実装
		user, err := model.SelectUserByPrimaryKey(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if user == nil {
			log.Println(errors.New("user not found"))
			response.BadRequest(writer, fmt.Sprintf("user not found. userID=%s", userID))
			return
		}

		allCollection, err := model.GetAllCollections()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		userCollection, err := model.GetUserCollections(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		var collectionListResponse CollectionListResponse
		collectionListResponse.Collections = make([]*CollectionList, 0, len(allCollection))

		collectionListResponseMap := make(map[string]*CollectionList, len(allCollection))
		for _, v := range allCollection {
			collectionList := CollectionList{
				CollectionID: v.ID,
				Name:         v.Name,
				Rarity:       v.Rarity,
				HasItem:      false,
			}
			collectionListResponseMap[v.ID] = &collectionList
			collectionListResponse.Collections = append(collectionListResponse.Collections, &collectionList)
		}
		for _, id := range userCollection {
			if v, ok := collectionListResponseMap[id]; ok {
				v.HasItem = true
			}
		}

		response.Success(writer, &collectionListResponse)

	}

}

//CollectionList データ
type CollectionList struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

//CollectionListResponse 返り値
type CollectionListResponse struct {
	Collections []*CollectionList `json:"collections"`
}
