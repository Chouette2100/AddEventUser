// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/Chouette2100/srdblib/v2"
)

// 条件にあうイベントを抽出し、ルームを登録する
func CollectAndAddEventUsers(
	client *http.Client,
	ie int,
	dtago int,
	dtfromnow int,
) (
	err error,
) {

	// 現在時刻を取得
	tnow := time.Now().Truncate(time.Second)

	// 現在開催中のイベントを抽出する
	sqlst := `
	SELECT * FROM event
	WHERE starttime BETWEEN DATE_SUB(NOW(), INTERVAL ? HOUR) AND DATE_ADD(NOW(), INTERVAL ? HOUR)
	ORDER BY starttime ASC, eventid ASC
	`

	type EventUrlKey struct {
		Eventid string `db:"eventid"`
	}

	eventinflist := []srdblib.Event{}
	_, erl := srdblib.Dbmap.Select(&eventinflist, sqlst, dtago, dtfromnow)
	if erl != nil {
		err = fmt.Errorf("Dbmap.Exec(): %w", erl)
		return
	}
	if len(eventinflist) == 0 {
		log.Printf("現在開催中のイベントはありません。\n")
		return
	}

	// 抽出したイベントそれぞれの参加者をDBに登録する
	for _, eil := range eventinflist {
		eventUrlKey := eil.Eventid
		if eventUrlKey != "weekend_short_008?block_id=85001" {
			log.Printf("非対象イベント %s %s [%s]\n", eil.Starttime.Format("2006-01-02 15:04"), eventUrlKey, eil.Event_name)
			continue
		}
		log.Printf("開催中開催予定 %s %s [%s]\n", eil.Starttime.Format("2006-01-02 15:04"), eventUrlKey, eil.Event_name)
		err = addNewUser(
			client,
			eil,
			ie,
		)
		if err != nil {
			log.Printf("addNewUser(): %s", err.Error())
		}
	}

	log.Printf("CollectAndAddEventUsers(): 処理が完了しました。処理時刻=%s\n", tnow.Format("2006-01-02 15:04:05"))

	return
}
