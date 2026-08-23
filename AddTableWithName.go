// Copyright © 2025 chouette.21.00@gmail.com
// Released under the MIT license
// https://opensource.org/licenses/mit-license.php
package main

import (
	"github.com/Chouette2100/srdblib/v3"
)

func AddTableWithName() {
	Dbmap.AddTableWithName(srdblib.User{}, "user").SetKeys(false, "Userno")
	Dbmap.AddTableWithName(srdblib.Wuser{}, "wuser").SetKeys(false, "Userno")
	Dbmap.AddTableWithName(srdblib.Userhistory{}, "userhistory").SetKeys(false, "Userno", "Ts")
	Dbmap.AddTableWithName(srdblib.Event{}, "event").SetKeys(false, "Eventid")
	Dbmap.AddTableWithName(srdblib.Wevent{}, "wevent").SetKeys(false, "Eventid")
	Dbmap.AddTableWithName(srdblib.Eventuser{}, "eventuser").SetKeys(false, "Eventid", "Userno")
	Dbmap.AddTableWithName(srdblib.Weventuser{}, "weventuser").SetKeys(false, "Eventid", "Userno")

	Dbmap.AddTableWithName(srdblib.GiftScore{}, "giftscore").SetKeys(false, "Giftid", "Ts", "Userno")
	Dbmap.AddTableWithName(srdblib.ViewerGiftScore{}, "viewergiftscore").SetKeys(false, "Giftid", "Ts", "Viewerid")
	Dbmap.AddTableWithName(srdblib.Viewer{}, "viewer").SetKeys(false, "Viewerid")
	Dbmap.AddTableWithName(srdblib.ViewerHistory{}, "viewerhistory").SetKeys(false, "Viewerid", "Ts")
	// srdblib.Dbmap.AddTableWithName(ShowroomCGIlib.Contribution{}, "contribution").SetKeys(false, "Ieventid", "Roomid", "Viewerid")

	Dbmap.AddTableWithName(srdblib.Campaign{}, "campaign").SetKeys(false, "Campaignid")
	Dbmap.AddTableWithName(srdblib.GiftRanking{}, "giftranking").SetKeys(false, "Campaignid", "Grid")
	Dbmap.AddTableWithName(srdblib.Accesslog{}, "accesslog").SetKeys(false, "Ts", "Eventid")

	Dbmap.AddTableWithName(srdblib.Todo{}, "todo").SetKeys(false, "ID")
}
