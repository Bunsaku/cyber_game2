package handler

import (
	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/db"
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

//HandleGachaDraw ガチャに関する処理
func HandleGachaDraw() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// リクエストBodyから更新後情報を取得
		var requestBody gachaDrawRequest
		json.NewDecoder(request.Body).Decode(&requestBody)

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println(errors.New("userID is empty"))
			response.BadRequest(writer, "Internal Server Error")
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

		times := requestBody.Times
		if times <= 0 {
			log.Println(errors.New("illegal value"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		//減らすコインの数を計算
		decreaseCoin := constant.GachaCoinConsumption * times
		if user.Coin >= decreaseCoin {
			user.Coin = user.Coin - decreaseCoin
		} else {
			log.Println(errors.New("not enough coin"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		var selectedCollectionItems []model.CollectionItem
		var NewItemIDs []string
		//トランザクション処理
		err = db.Transact(db.Conn, func(tx *sql.Tx) error {
			//コインを減らす処理
			err := model.DcreaseUserCoin(tx, user)
			if err != nil {
				return err
			}
			//ガチャの出現頻度データを取得
			items, err := model.GetAllGachaProbability(tx)
			if err != nil {
				return err
			}
			//抽選する
			selectedCollectionItemIDs := selecteItems(times, items)
			//抽選したIDをもとにコレクションアイテムデータを取得
			selectedCollectionItems, err = model.GetCollectionItem(tx, selectedCollectionItemIDs)
			if err != nil {
				return err
			}
			//新しいアイテムかどうかを判別
			NewItemIDs, err = model.HasItem(tx, selectedCollectionItemIDs, userID)
			if err != nil {
				return err
			}

			//NewItemIDsが空ではないときに、ユーザーコレクションアイテムに挿入
			if len(NewItemIDs) != 0 {
				err = model.UpdatetUserCollectionItems(tx, NewItemIDs, userID)
				if err != nil {
					fmt.Println(err)
					return err
				}
			}
			return err
		})

		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		//responseの形に整形
		var gachaResponse GachaDrawResponse

		for _, v := range selectedCollectionItems {
			collectionItems := CollectionItems{
				CollectionID: v.ID,
				Name:         v.Name,
				Rarity:       v.Rarity,
				IsNew:        false,
			}
			gachaResponse.Results = append(gachaResponse.Results, &collectionItems)
		}

		gachaResponseMap := make(map[string]*CollectionItems, len(gachaResponse.Results))

		for _, v := range gachaResponse.Results {
			gachaResponseMap[v.CollectionID] = v
		}

		for _, id := range NewItemIDs {
			if v, ok := gachaResponseMap[id]; ok {
				v.IsNew = true
			}
		}

		response.Success(writer, &gachaResponse)
	}

}

type gachaDrawRequest struct {
	Times int32 `json:"times"`
}

//CollectionItems データ
type CollectionItems struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int    `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

//GachaDrawResponse is 'GacahDrawの返り値'
type GachaDrawResponse struct {
	Results []*CollectionItems `json:"results"`
}

func selecteItems(times int32, items []model.GachaProbability) []string {
	var sumRatio int
	for _, item := range items {
		sumRatio += item.Ratio
	}
	var selectedCollectionItems []string
	var selectedCollectionItem model.GachaProbability
	//抽選
	for i := 0; i < int(times); i++ {
		randNum := rand.Intn(sumRatio)
		for _, item := range items {
			if randNum < item.Ratio {
				selectedCollectionItem = item
				break
			}
			randNum -= item.Ratio
		}
		selectedCollectionItems = append(selectedCollectionItems, selectedCollectionItem.CollectionItemID)
	}
	return selectedCollectionItems
}
