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

package gf_admin_lib

import (
	// "fmt"
	"net/http"
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func init_handlers__users(pHTTPmux *http.ServeMux,
	pServiceInfo         *GF_service_info,
	pIdentityServiceInfo *gf_identity_lib.GF_service_info,
	pLocalHub            *sentry.Hub,
	pRuntimeSys          *gf_core.Runtime_sys) *gf_core.GF_error {

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/admin/users/get_all",
		"/v1/admin/users/get_all_invite_list",
		"/v1/admin/users/add_to_invite_list",
	}
	metrics := gf_rpc_lib.Metrics__create_for_handlers("gf_admin", handlers_endpoints_lst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpcHandlerRuntime := &gf_rpc_lib.GF_rpc_handler_runtime {
		Mux:                pHTTPmux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         pLocalHub,
		Auth_login_url_str: "/v1/admin/login_ui",
	}

	//---------------------
	// GET_ALL
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/get_all",
		func(pCtx context.Context, p_resp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				_, userIDstr, _, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, p_resp, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}


				gfErr = gf_identity_lib.AdminIs(userIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				//---------------------

				usersLst, gfErr := gf_identity_lib.AdminPipelineGetAllUsers(pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					"users_lst": usersLst,
				}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// GET_ALL_INVITE_LIST
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/get_all_invite_list",
		func(pCtx context.Context, p_resp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				_, userIDstr, _, gf_err := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, p_resp, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}


				gf_err = gf_identity_lib.AdminIs(userIDstr, pCtx, pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------

				invite_list_lst, gf_err := gf_identity_lib.Admin__pipeline__get_all_invite_list(pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"invite_list_lst": invite_list_lst,
				}
				return output_map, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// ADD_TO_INVITE_LIST
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/add_to_invite_list",
		func(pCtx context.Context, p_resp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				input_map, userIDstr, _, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, p_resp, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var email_str string
				if val_str, ok := input_map["email_str"]; ok {
					email_str = val_str.(string)
				}

				gfErr = gf_identity_lib.AdminIs(userIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input := &gf_identity_lib.GF_admin__input_add_to_invite_list{
					UserIDstr: userIDstr,
					EmailStr:  email_str,
				}

				//---------------------

				gfErr = gf_identity_lib.AdminPipelineUserAddToInviteList(input,
					pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				output_map := map[string]interface{}{
					
				}
				return output_map, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	// REMOVE_FROM_INVITE_LIST
	// AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(true, "/v1/admin/users/remove_from_invite_list",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gf_core.GF_error) {

			if pReq.Method == "POST" {

				//---------------------
				// INPUT
				
				inputMap, userIDstr, _, gfErr := gf_identity_core.HTTPgetUserStdInput(pCtx, pReq, pResp, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				var emailStr string
				if valStr, ok := inputMap["email_str"]; ok {
					emailStr = valStr.(string)
				}

				gfErr = gf_identity_lib.AdminIs(userIDstr, pCtx, pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				input := &gf_identity_lib.GFadminRemoveFromInviteListInput{
					UserIDstr: userIDstr,
					EmailStr:  emailStr,
				}

				//---------------------

				gfErr = gf_identity_lib.AdminPipelineUserRemoveFromInviteList(input,
					pCtx,
					pIdentityServiceInfo,
					pRuntimeSys)
				if gfErr != nil {
					return nil, gfErr
				}

				outputMap := map[string]interface{}{
					
				}
				return outputMap, nil
			}
			return nil, nil
		},
		rpcHandlerRuntime,
		pRuntimeSys)

	//---------------------
	
	return nil
}