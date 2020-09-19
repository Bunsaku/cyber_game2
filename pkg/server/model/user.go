package model

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/db"
)

// User userテーブルデータ
type User struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int32
}

//RankingList データ
type RankingList struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int    `json:"rank"`
	Score    int32  `json:"score"`
}

//RankingListResponse is 'GetRankingListの返り値'
type RankingListResponse struct {
	Ranks []RankingList `json:"ranks"`
}

// InsertUser データベースをレコードを登録する
func InsertUser(record *User) error {
	// userテーブルへのレコードの登録を行うSQLを入力する
	stmt, err := db.Conn.Prepare("INSERT INTO user (id, auth_token, name, high_score) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.ID, record.AuthToken, record.Name, record.HighScore)
	return err
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func SelectUserByAuthToken(authToken string) (*User, error) {
	// auth_tokenを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT * FROM user WHERE auth_token = ?", authToken)
	return convertToUser(row)
}

// SelectUserByPrimaryKey 主キーを条件にレコードを取得する
func SelectUserByPrimaryKey(userID string) (*User, error) {
	// idを条件にSELECTを行うSQLを第1引数に入力する
	row := db.Conn.QueryRow("SELECT * FROM user WHERE id = ?", userID)
	return convertToUser(row)
}

// UpdateUserByPrimaryKey 主キーを条件にレコードを更新する
func UpdateUserByPrimaryKey(record *User) error {
	// idを条件に指定した値でnameカラムの値を更新するSQLを入力する
	stmt, err := db.Conn.Prepare("UPDATE user SET name = ? WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.Name, record.ID)
	return err
}

//UpdateUserScoreCoin is 'userのスコアとコインを更新する'
func UpdateUserScoreCoin(user *User, addedScore int32) error {
	//ハイスコアよりaddedScoreが大きければ、値を更新しそうでなければ更新しない、コインはどちらの場合でもインクリメントする
	if user.HighScore < addedScore {
		stmt, err := db.Conn.Prepare("UPDATE user SET high_score = ? WHERE id = ? ")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(addedScore, user.ID)
		return err
	}
	return nil
}

//GetRankingList is 'ranking listを取得する関数'
func GetRankingList() (*RankingListResponse, error) {
	//データベースからランキングの取得
	rows, err := db.Conn.Query("SELECT id,name,high_score FROM user ORDER BY high_score DESC LIMIT 5 ")
	if err != nil {
		return nil, err
	}
	var RankingListResponses RankingListResponse
	r := RankingList{}
	i := 1
	for rows.Next() {
		err = rows.Scan(&r.UserID, &r.UserName, &r.Score)
		if err != nil {
			return nil, err
		}
		r.Rank = i
		RankingListResponses.Ranks = append(RankingListResponses.Ranks, r)
		i++
	}
	return &RankingListResponses, nil
}

// convertToUser rowデータをUserデータへ変換する
func convertToUser(row *sql.Row) (*User, error) {
	user := User{}
	err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &user, nil
}
