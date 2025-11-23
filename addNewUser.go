// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"fmt"
	// "strconv"
	// "strings"
	// "io"
	"log"
	"time"
	// "os"

	"net/http"
	// "reflect"
	// "testing"
	// "golang.org/x/tools/go/analysis/passes/defers"

	// "github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

// イベント参加者をDBに登録する
func addNewUser(
	client *http.Client,
	event srdblib.Event,
	ie int,
) (
	err error,
) {
	var eul *EventuserList
	NoOfRooms := 0
	eul, NoOfRooms, _, err = GetEventuserList(
		client,
		event,
		1,
		ie,
		true,
		false,
	)
	if err != nil {
		err = fmt.Errorf("GetEventuserList(): %w", err)
		return
	}

	// ルーム数が多すぎるときは次の処理に移る
	if NoOfRooms > ie {
		log.Printf("    イベント参加者数が多すぎます。len(*eul)=%d NoOfRooms=%d ie=%d",
			len(*eul), NoOfRooms, ie)
		return
	} else {
		log.Printf("    イベント参加者数が許容範囲内です。len(*eul)=%d NoOfRooms=%d ie=%d",
			len(*eul), NoOfRooms, ie)
	}

	tnow := time.Now().Truncate(time.Second)
	tstarttime := event.Starttime.Truncate(time.Second)

	log.Printf("イベント参加者リストの取得に成功しました。len(eul)=%d NoOfRooms=%d\n",
		len(*eul), NoOfRooms)
	// イベント参加者リストと開始時ポイントをDBに登録する
	for _, eu := range *eul {
		err = UpsertEventuser(&eu, tnow)
		if err != nil {
			err = fmt.Errorf("UpsertEventuser(): %w", err)
			return
		}
		err = UpsertPoints(&eu, tstarttime)
		if err != nil {
			err = fmt.Errorf("UpsertPoints(): %w", err)
			return
		}
	}
	log.Printf("イベント参加者リストのDB登録に成功しました。count=%d\n", len(*eul))

	return

}
