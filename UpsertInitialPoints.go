package main

import (
	"fmt"
	// "strconv"
	// "strings"
	// "io"
	"log"
	"time"
	// "os"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	// "net/http"
	// "reflect"
	// "testing"
	// "golang.org/x/tools/go/analysis/passes/defers"

	// "github.com/Chouette2100/srapi/v2"
	"github.com/Chouette2100/srdblib/v2"
)

const sqlUpsertPoints = "INSERT INTO points ( ts, user_id, eventid,  point, `rank`, gap, pstatus) " +
	"    VALUES (?, ? , ?, 0, 1, 0, '=') " +
	" ON DUPLICATE KEY UPDATE " +
	"     point = point "

// イベント参加者をDBに登録する、重複したデータは何もしない
func UpsertInitialPoints(
	eu *srdblib.Eventuser,
	tstarttime time.Time,
) (
	err error,
) {

	var result sql.Result
	result, err = srdblib.Dbmap.Exec(
		sqlUpsertPoints,
		tstarttime,
		eu.Userno,
		eu.Eventid,
	)
	if err != nil {
		err = fmt.Errorf("Dbmap.Exec(): %w", err)
		log.Printf("    開始時ポイントの登録に失敗しました。 eventid=%s userno=%d result=%+v err=%s\n",
			eu.Eventid, eu.Userno, result, err)
		return
	}

	return
}
