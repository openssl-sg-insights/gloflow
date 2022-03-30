/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

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

package gf_home_lib

import (
	"net/http"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GFserviceInfo struct {

	// AUTH_LOGIN_URL - url of the login page to which the system should
	//                  redirect users after certain operations
	AuthLoginURLstr string
}

//-------------------------------------------------
func InitService(pTemplatesPathsMap map[string]string,
	pServiceInfo *GFserviceInfo,
	pHTTPmux     *http.ServeMux,
	pRuntimeSys  *gf_core.Runtime_sys) *gf_core.GF_error {

	//------------------------
	// HANDLERS
	gfErr := initHandlers(pTemplatesPathsMap,
		pServiceInfo.AuthLoginURLstr,
		pHTTPmux,
		pRuntimeSys)
	if gfErr != nil {
		return gfErr
	}

	//------------------------

	return nil
}