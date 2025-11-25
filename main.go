// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-gorp/gorp"
	"golang.org/x/term"

	"github.com/Chouette2100/exsrapi/v2"
	"github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

/*
100101 2025-11-23 最初のバージョン
100201 2025-11-23 時間指定を時間単位（d/h/m）で指定できるように変更する
100301 2025-11-24 GetEventQuestRoomsByApi()をGetEventRoomsByApi()に変更する
100302 2025-11-25 eventuserテーブルに登録後、userテーブルに存在しないuseridを新規登録する機能を追加する
*/

const Version = "100302"

// イベントの参加者を調べ、一定数以下ならDB(eventuser)に登録する
func main() {

	// ログファイルの作成
	logfile, err := exsrapi.CreateLogfile(Version, exsrapi.Version, srapi.Version, srdblib.Version)
	if err != nil {
		log.Printf("ログファイルの作成に失敗しました。%v\n", err)
		return
	}
	defer logfile.Close()

	// フォアグラウンド（端末に接続されているか）を判定
	isForeground := term.IsTerminal(int(os.Stdout.Fd()))
	if isForeground {
		// フォアグラウンドならログファイル + コンソール
		log.SetOutput(io.MultiWriter(logfile, os.Stdout))
	} else {
		// バックグラウンドならログファイルのみ
		log.SetOutput(logfile)
	}

	log.SetFlags(log.Lmicroseconds)
	log.Printf("Version=%s Start\n", Version)
	// --------------------------------

	// DB接続
	var dbconfig *srdblib.DBConfig
	dbconfig, err = srdblib.OpenDb("DBConfig.yml")
	if err != nil {
		log.Printf("Database error. err = %v\n", err)
		return
	}
	if dbconfig.UseSSH {
		defer srdblib.Dialer.Close()
	}
	defer srdblib.Db.Close()
	srdblib.Db.SetMaxOpenConns(8)
	srdblib.Db.SetMaxIdleConns(12)

	srdblib.Db.SetConnMaxLifetime(time.Minute * 5)
	srdblib.Db.SetConnMaxIdleTime(time.Minute * 5)

	defer srdblib.Db.Close()
	log.Printf("%+v\n", dbconfig)

	dial := gorp.MySQLDialect{Engine: "InnoDB", Encoding: "utf8mb4"}
	srdblib.Dbmap = &gorp.DbMap{Db: srdblib.Db,
		Dialect:         dial,
		ExpandSliceArgs: true, //スライス引数展開オプションを有効化する
	}
	AddTableWithName()
	// --------------------------------

	// userテーブルの更新判定の閾値、ApiRoomProfile()の実行頻度を設定する
	fileenv := "Env.yml"
	err = exsrapi.LoadConfig(fileenv, &srdblib.Env)
	if err != nil {
		err = fmt.Errorf("exsrapi.Loadconfig(): %w", err)
		log.Printf("%s\n", err.Error())
		return
	}
	log.Printf("Env.yml: %+v\n", srdblib.Env)
	// --------------------------------

	/// 環境変数から設定値を取得する
	snorooms := os.Getenv("SR_ADD_EVENTUSER_NOROOMS")
	norooms, _ := strconv.Atoi(snorooms)
	dtago := os.Getenv("SR_ADD_EVENTUSER_DTAGO")
	dtfromnow := os.Getenv("SR_ADD_EVENTUSER_DTFROMNOW")

	log.Printf("設定値: SR_ADD_EVENTUSER_NOROOMS=%d SR_ADD_EVENTUSER_DTAGO=%s SR_ADD_EVENTUSER_DTFROMNOW=%s\n",
		norooms, dtago, dtfromnow)

	/* 起動時パラメータの取得
	var eventid string
	var ibreg, iereg int
	// 起動時パラメータからeventid, ibreg, ieregを取得する。
	if len(os.Args) < 4 {
		log.Printf("Usage: srAddEvent eventid ibreg iereg\n")
		return
	}
	eventid = os.Args[1]
	ibreg, _ = strconv.Atoi(os.Args[2])
	iereg, _ = strconv.Atoi(os.Args[3])
	log.Printf(" eventid =[%s], ibreg=%d, iereg=%d\n", eventid, ibreg, iereg)
	client := http.DefaultClient
	addNewUser(client, eventid, 40)
	*/

	// 現在開催中のイベントの参加者をDBに登録する
	err = CollectAndAddEventUsers(
		http.DefaultClient,
		norooms,
		dtago,
		dtfromnow,
	)
	if err != nil {
		log.Printf("CollectAndAddEventUsers error: %v\n", err)
		return
	}

	// 追加されたuseridがuserテーブルに存在しないときは新規登録する
	err = SearchAndAddNewUsers(3) // 3時間以内に始まるイベントも対象とする
	if err != nil {
		log.Printf("SearchAndAddNewUsers error: %v\n", err)
		return
	}
	// --------------------------------
	log.Printf("End\n")
}
