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

const sqlUpsertEventuser = `
INSERT INTO eventuser ( eventid, userno, istarget, iscntrbpoints, graph, color, point, vld, status)
    VALUES (?, ? , 'Y', 'N', ?, 'red', 0, 1, 0)
ON DUPLICATE KEY UPDATE
    point = point ;
`

// イベント参加者をDBに登録する、重複したデータは何もしない
func UpsertEventuser(
	eu *srdblib.Eventuser,
	tnow time.Time,
) (
	err error,
) {

	var result sql.Result
	result, err = srdblib.Dbmap.Exec(
		sqlUpsertEventuser,
		eu.Eventid,
		eu.Userno,
		eu.Graph,
	)
	if err != nil {
		err = fmt.Errorf("Dbmap.Exec(): %w", err)
		log.Printf("    イベント参加者の登録に失敗しました。 eventid=%s userno=%d result=%+v err=%s\n",
			eu.Eventid, eu.Userno, result, err)
		return
	}

	return
}
