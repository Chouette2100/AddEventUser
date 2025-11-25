package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Chouette2100/srdblib/v2"
)

// SearchAndAddNewUsers は、指定された時間内に開始するイベントの参加者を調べ、
// userテーブルに存在しないユーザーを新規登録する
func SearchAndAddNewUsers(
	hoursAhead int, // 何時間以内に開始するイベントを対象とするか
) (
	err error,
) {
	log.Printf("SearchAndAddNewUsers: hoursAhead=%d Start\n", hoursAhead)
	// --------------------------------

	// userテーブルに存在しないルームを抽出する
	type NewUser struct {
		// gorpでSELECT結果をマッピングするため公開(Export)が必要
		// フィールド名はテーブル列 userno に合わせ、タグで明示
		Userno int `db:"userno"`
	}
	var newUsers []NewUser

	sqlst := `
	SELECT DISTINCT eu.userno
	FROM eventuser eu
	LEFT JOIN user u ON eu.userno = u.userno
	JOIN event e ON eu.eventid = e.eventid
	WHERE u.userno IS NULL
		AND NOW() BETWEEN DATE_SUB(e.starttime, INTERVAL ? HOUR) AND e.endtime
	`

	_, err = srdblib.Dbmap.Select(&newUsers, sqlst, hoursAhead)
	if err != nil {
		err = fmt.Errorf("新規ユーザー抽出エラー: %w", err)
		return
	}

	log.Printf("新規登録対象ユーザー数=%d\n", len(newUsers))

	// userテーブルに新規登録する
	user := srdblib.User{}
	for _, nu := range newUsers {
		user.Userno = nu.Userno
		_, err = srdblib.UpinsUser(http.DefaultClient, time.Now().Truncate(time.Second), &user)
		if err != nil {
			err = fmt.Errorf("ユーザー新規登録エラー: userno=%d: %w", nu.Userno, err)
			return
		}
		log.Printf("ユーザー新規登録: userno=%d\n", nu.Userno)
	}

	// --------------------------------
	log.Printf("SearchAndAddNewUsers End\n")
	return
}
