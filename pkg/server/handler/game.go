package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

//HandleGameFinish is 'ユーザーのスコアとコインを更新'
func HandleGameFinish() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// リクエストBodyから更新後情報を取得
		var requestBody gameFinishRequest
		json.NewDecoder(request.Body).Decode(&requestBody)

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

		//userテーブルのハイスコアとコインの更新処理

		addedScore := requestBody.Score
		addedCoin := addedScore
		err = model.UpdateUserScoreCoin(user, addedScore, addedCoin)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		response.Success(writer, &gameFinishResponse{Coin: addedCoin})

	}
}

type gameFinishRequest struct {
	Score int32 `json:"score"`
}
type gameFinishResponse struct {
	Coin int32 `json:"coin"`
}
