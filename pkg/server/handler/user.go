package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"

	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

// HandleUserCreate ユーザ情報作成処理
func HandleUserCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// リクエストBodyから更新後情報を取得
		var requestBody userCreateRequest
		err := json.NewDecoder(request.Body).Decode(&requestBody)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// UUIDでユーザIDを生成する
		userID, err := uuid.NewRandom()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// UUIDで認証トークンを生成する
		authToken, err := uuid.NewRandom()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// データベースにユーザデータを登録する
		err = model.InsertUser(&model.User{
			ID:        userID.String(),
			AuthToken: authToken.String(),
			Name:      requestBody.Name,
			HighScore: 0,
		})
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// 生成した認証トークンを返却
		response.Success(writer, &userCreateResponse{Token: authToken.String()})
	}
}

type userCreateRequest struct {
	Name string `json:"name"`
}

type userCreateResponse struct {
	Token string `json:"token"`
}
