// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/go-gorp/gorp"

	"github.com/Chouette2100/srdblib/v2"
)

func TestGetEventuserList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		client     *http.Client
		event      srdblib.Event
		ib         int
		ie         int
		isGetList  bool
		isGetPoint bool
		want       *EventuserList
		want1      int
		want2      int
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "listenerupupup_showroom260",
			client:     &http.Client{},
			event:      srdblib.Event{Eventid: "popteen_akb48g_model"},
			ib:         1,
			ie:         5,
			isGetList:  true,
			isGetPoint: false,
			want:       &EventuserList{},
			want1:      0,
			want2:      0,
			wantErr:    false,
		},
		/*
			{
				name:        "hanakin_happy_night_007?block_id=85901",
				client:      &http.Client{},
				eventUrlKey: "hanakin_happy_night_007?block_id=85901",
				ib:          1,
				ie:          5,
				isGetList:   true,
				isGetPoint:  false,
				want:        &main.EventuserList{},
				want1:       0,
				want2:       0,
				wantErr:     false,
			},
		*/
	}
	// --------------------------------

	// DB接続
	var err error
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, gotErr := GetEventuserList(tt.client, tt.event, tt.ib, tt.ie, tt.isGetList, tt.isGetPoint)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetEventuserList() failed: %s", gotErr.Error())
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetEventuserList() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetEventuserList() = %v, want %v", got, tt.want)
			}
			if true {
				t.Errorf("GetEventuserList() = %v, want1 %v", got1, tt.want1)
			}
			if true {
				t.Errorf("GetEventuserList() = %v, want2 %v", got2, tt.want2)
			}
		})
	}
}
