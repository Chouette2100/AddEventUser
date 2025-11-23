// Copyright © 2052 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	// "os"

	"net/http"
	// "reflect"
	// "testing"
	// "golang.org/x/tools/go/analysis/passes/defers"

	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

type EventuserList []srdblib.Eventuser

func GetEventuserList(
	client *http.Client,
	event srdblib.Event,
	ib int,
	ie int,
	isGetList bool,
	isGetPoint bool,
) (
	eul *EventuserList,
	NoOfRooms int, // これは取得したルーム数ではなくイベントに参加しているルーム数
	lenOfPointList int,
	err error,
) {

	// DBからイベント情報を取得する
	eventinfo := srdblib.Event{}
	// var intf interface{}
	intf, erl := srdblib.Dbmap.Get(&eventinfo, event.Eventid)
	if erl != nil {
		err = fmt.Errorf("srdblib.Dbmap.Get(): %w", err)
		return
	}
	if intf == nil {
		err = fmt.Errorf("イベント情報が見つかりません。eventUrlKey=%s", event.Eventid)
		return
	}
	eventinfo = *intf.(*srdblib.Event)
	// イベントIDを取得する
	ieventid := eventinfo.Ieventid

	if event.Starttime.After(time.Now()) {
		dt := time.Since(event.Starttime).Hours()
		iet := ie + (int(dt)-1)*4 // 開始まで時間があるときはその時間に応じて取得上限を下げる
		if iet < 7 {
			iet = 7 // 最低でも7件は取得する
		}
		if iet < ie {
			if iet < ib {
				err = fmt.Errorf("取得上限が取得下限より小さくなりました。ib=%d ie=%d", ib, ie)
				return
			}
			ie = iet
			log.Printf("取得上限を調整しました。ie=%d\n", ie)
		}
	}

	// イベント参加者リストを取得する
	eua := strings.Split(event.Eventid, "?block_id=")
	if len(eua) == 2 {
		// ブロックイベントのとき
		if event.Starttime.After(time.Now()) {
			err = fmt.Errorf("ブロックイベントがまだ開始されていません。eventid=%s starttime=%s",
				event.Eventid, event.Starttime.Format("2006-01-02 15:04:05"))
			return
		}
		var blockid int
		blockid, _ = strconv.Atoi(eua[1])
		eventid := eua[0]
		log.Printf("eventid=%s  ieventid=%d blockid=%d\n",
			eventid, ieventid, blockid)

		ebr, erl := srapi.GetEventBlockRanking(client, ieventid, blockid, ib, ie)
		if erl != nil {
			err = fmt.Errorf("srapi.GetEventBlockRanking(): %w", erl)
			return
		}

		// イベント参加者数を取得する
		NoOfRooms = ebr.Total_entries
		eventinfo.Noentry = NoOfRooms
		srdblib.Dbmap.Update(&eventinfo)

		eua := make([]srdblib.Eventuser, len(ebr.Block_ranking_list))
		for i, br := range ebr.Block_ranking_list {
			eua[i].Eventid = event.Eventid
			eua[i].Userno, _ = strconv.Atoi(br.Room_id)
			eua[i].Istarget = "Y"
			eua[i].Iscntrbpoints = "N"
			eua[i].Graph = "Y"
			eua[i].Color = "red"
			eua[i].Point = br.Point
			eua[i].Vld = 0
		}
		eul = (*EventuserList)(&eua)
	} else {
		// ブロックイベントでないとき（通常のイベント、レベルイベント）
		eqr, erl := srapi.GetEventQuestRoomsByApi(client, event.Eventid, ib, ie)
		if erl != nil {
			err = fmt.Errorf("srapi.GetEventQuestRoomsByApi(): %w", erl)
			return
		}
		// イベント参加者数を取得する
		NoOfRooms = eqr.TotalEntries
		eventinfo.Noentry = NoOfRooms
		srdblib.Dbmap.Update(&eventinfo)

		eua := make([]srdblib.Eventuser, len(eqr.EventQuestLevelRanges[0].Rooms))
		for i, qr := range eqr.EventQuestLevelRanges[0].Rooms {
			eua[i].Eventid = event.Eventid
			eua[i].Userno = qr.RoomID
			eua[i].Istarget = "Y"
			eua[i].Iscntrbpoints = "N"
			eua[i].Graph = "Y"
			eua[i].Color = "red"
			eua[i].Point = qr.Point
			// eua[i].Vld = qr.QuestLevel
			eua[i].Vld = 0
		}
		eul = (*EventuserList)(&eua)
	}

	return
}
