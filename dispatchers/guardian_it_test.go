// +build integration

/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package dispatchers

import (
	"reflect"
	"testing"
	"time"

	"github.com/cgrates/cgrates/utils"
)

var sTestsDspGrd = []func(t *testing.T){
	testDspGrdPing,
}

//Test start here
func TestDspGuardianSTMySQL(t *testing.T) {
	testDsp(t, sTestsDspGrd, "TestDspGuardianS", "all", "all2", "attributes", "dispatchers", "tutorial", "oldtutorial", "dispatchers")
}

func TestDspGuardianSMongo(t *testing.T) {
	testDsp(t, sTestsDspGrd, "TestDspGuardianS", "all", "all2", "attributes_mongo", "dispatchers_mongo", "tutorial", "oldtutorial", "dispatchers")
}

func testDspGrdPing(t *testing.T) {
	var reply string
	if err := allEngine.RCP.Call(utils.GuardianSv1Ping, new(utils.CGREvent), &reply); err != nil {
		t.Error(err)
	} else if reply != utils.Pong {
		t.Errorf("Received: %s", reply)
	}
	if err := dispEngine.RCP.Call(utils.GuardianSv1Ping, &CGREvWithApiKey{
		CGREvent: utils.CGREvent{
			Tenant: "cgrates.org",
		},
		DispatcherResource: DispatcherResource{
			APIKey: "grd12345",
		},
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.Pong {
		t.Errorf("Received: %s", reply)
	}
}

func testDspGrdLock(t *testing.T) {
	// lock
	args := utils.AttrRemoteLock{
		ReferenceID: "",
		LockIDs:     []string{"lock1"},
		Timeout:     500 * time.Millisecond,
	}
	var reply string
	if err := dispEngine.RCP.Call(utils.GuardianSv1RemoteLock, &AttrRemoteLockWithApiKey{
		AttrRemoteLock: args,
		TenantArg: utils.TenantArg{
			Tenant: "cgrates.org",
		},
		DispatcherResource: DispatcherResource{
			APIKey: "grd12345",
		},
	}, &reply); err != nil {
		t.Error(err)
	} else if reply != utils.Pong {
		t.Errorf("Received: %s", reply)
	}

	var unlockReply []string
	if err := dispEngine.RCP.Call(utils.GuardianSv1RemoteUnlock, &AttrRemoteUnlockWithApiKey{
		RefID: reply,
		TenantArg: utils.TenantArg{
			Tenant: "cgrates.org",
		},
		DispatcherResource: DispatcherResource{
			APIKey: "grd12345",
		},
	}, &unlockReply); err != nil {
		t.Error(err)
	} else if reply != utils.Pong {
		t.Errorf("Received: %s", reply)
	}
	if !reflect.DeepEqual(args.LockIDs, unlockReply) {
		t.Errorf("Expected: %s , received: %s", utils.ToJSON(args.LockIDs), utils.ToJSON(unlockReply))
	}
}