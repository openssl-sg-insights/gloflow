/*
GloFlow application and media management/publishing platform
Copyright (C) 2022 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_web3_lib

import (
	"testing"
	"context"
	"github.com/parnurzeal/gorequest"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
)

//---------------------------------------------------
func TestAddresses(pTest *testing.T) {


	runtime, _, err := gf_eth_core.TgetRuntime()
	if err != nil {
		pTest.FailNow()
	}

	// testWeb3MonitorServiceInt  := 2000
	testIdentityServicePortInt := 2001
	HTTPagent := gorequest.New()
	ctx       := context.Background()

	// CREATE_AND_LOGIN_NEW_USER
	gf_identity_lib.TestCreateAndLoginNewUser(pTest,
		HTTPagent,
		testIdentityServicePortInt,
		ctx,
		runtime.RuntimeSys)	




}