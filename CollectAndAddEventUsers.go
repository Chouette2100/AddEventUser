// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"net/http"

	"github.com/Chouette2100/srdblib/v2"
)

// parseTimeInterval は "120m", "2h", "3d", "-30m", "-3d" のような文字列を解析し、
// MySQLのINTERVAL句で使用する数値と単位を返す
func parseTimeInterval(s string) (value int, unit string, err error) {
	// 符号付き数字と単位(m/h/d)を抽出する正規表現
	re := regexp.MustCompile(`^(-?\d+)([mhd])$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 3 {
		err = fmt.Errorf("無効な時間形式: %s (例: 120m, -30m, 2h, -3d)", s)
		return
	}

	value, err = strconv.Atoi(matches[1])
	if err != nil {
		return
	}

	switch matches[2] {
	case "m":
		unit = "MINUTE"
	case "h":
		unit = "HOUR"
	case "d":
		unit = "DAY"
	}

	return
}

// 条件にあうイベントを抽出し、ルームを登録する
func CollectAndAddEventUsers(
	client *http.Client,
	ie int,
	dtago string,
	dtfromnow string,
) (
	err error,
) {

	// 現在時刻を取得
	tnow := time.Now().Truncate(time.Second)

	// 時間指定をパースする
	agoValue, agoUnit, err := parseTimeInterval(dtago)
	if err != nil {
		err = fmt.Errorf("dtago解析エラー: %w", err)
		return
	}
	fromnowValue, fromnowUnit, err := parseTimeInterval(dtfromnow)
	if err != nil {
		err = fmt.Errorf("dtfromnow解析エラー: %w", err)
		return
	}

	// SQL関数と値を決定（可読性のため、負の値の場合は関数を切り替える）
	var agoFunc string
	var agoAbsValue int
	if agoValue < 0 {
		agoFunc = "DATE_ADD"
		agoAbsValue = -agoValue
	} else {
		agoFunc = "DATE_SUB"
		agoAbsValue = agoValue
	}

	var fromnowFunc string
	var fromnowAbsValue int
	if fromnowValue < 0 {
		fromnowFunc = "DATE_SUB"
		fromnowAbsValue = -fromnowValue
	} else {
		fromnowFunc = "DATE_ADD"
		fromnowAbsValue = fromnowValue
	}

	// 現在開催中のイベントを抽出する
	// SQLインジェクションの可能性がない部分なので、単位は直接文字列生成
	sqlst := fmt.Sprintf(`
	SELECT * FROM event
	WHERE starttime BETWEEN %s(NOW(), INTERVAL %d %s) AND %s(NOW(), INTERVAL %d %s)
	  AND endtime > NOW()
	ORDER BY starttime ASC, eventid ASC
	`, agoFunc, agoAbsValue, agoUnit, fromnowFunc, fromnowAbsValue, fromnowUnit)

	log.Printf("生成されたSQL: %s\n", sqlst)

	type EventUrlKey struct {
		Eventid string `db:"eventid"`
	}

	eventinflist := []srdblib.Event{}
	_, erl := srdblib.Dbmap.Select(&eventinflist, sqlst)
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
		// if eventUrlKey != "weekend_short_008?block_id=85001" {
		// 	log.Printf("非対象イベント %s %s [%s]\n", eil.Starttime.Format("2006-01-02 15:04"), eventUrlKey, eil.Event_name)
		// 	continue
		// }
		log.Printf("開催中開催予定 %s %s [%s]\n", eil.Starttime.Format("2006-01-02 15:04"), eventUrlKey, eil.Event_name)
		err = addNewUser(
			client,
			eil,
			ie,
		)
		if err != nil {
			log.Printf("addNewUser(): %s", err.Error())
			err = nil
		}
	}

	log.Printf("CollectAndAddEventUsers(): 処理が完了しました。処理時刻=%s\n", tnow.Format("2006-01-02 15:04:05"))

	return
}
