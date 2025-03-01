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

package gf_identity_lib

import (
	// "fmt"
	"net/http"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_session"
	// "github.com/davecgh/go-spew/spew"
)

//------------------------------------------------
func initHandlersEth(p_http_mux *http.ServeMux,
	p_service_info *GFserviceInfo,
	p_runtime_sys  *gf_core.RuntimeSys) *gf_core.GFerror {

	//---------------------
	// METRICS
	handlers_endpoints_lst := []string{
		"/v1/identity/eth/preflight",
		"/v1/identity/eth/login",
		"/v1/identity/eth/create",
	}
	metricsGroupNameStr := "eth"
	metrics := gf_rpc_lib.MetricsCreateForHandlers(metricsGroupNameStr, p_service_info.Name_str, handlers_endpoints_lst)

	//---------------------
	// RPC_HANDLER_RUNTIME
	rpc_handler_runtime := &gf_rpc_lib.GFrpcHandlerRuntime {
		Mux:                p_http_mux,
		Metrics:            metrics,
		Store_run_bool:     true,
		Sentry_hub:         nil,
		Auth_login_url_str: "/landing/main",
	}

	//---------------------
	// USERS_PREFLIGHT
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/eth/preflight",
		func(pCtx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				_, _, user_address_eth_str, gf_err := gf_identity_core.HTTPgetUserStdInput(pCtx, p_req, p_resp, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input :=&GF_user_auth_eth__input_preflight{
					User_address_eth_str: user_address_eth_str,
				}

				//---------------------

				output, gf_err := users_auth_eth__pipeline__preflight(input, pCtx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"user_exists_bool": output.User_exists_bool,
					"nonce_val_str":    output.Nonce_val_str,
				}
				return output_map, nil
			}

			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// USERS_LOGIN
	// NO_AUTH
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/eth/login",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, _, user_address_eth_str, gf_err := gf_identity_core.HTTPgetUserStdInput(p_ctx, p_req, p_resp, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}
				auth_signature_str := gf_identity_core.GF_auth_signature(input_map["auth_signature_str"].(string))

				input :=&GF_user_auth_eth__input_login{
					User_address_eth_str: user_address_eth_str,
					Auth_signature_str:   auth_signature_str,
				}

				//---------------------
				// LOGIN
				output, gf_err := users_auth_eth__pipeline__login(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				//---------------------
				// SET_SESSION_ID - sets gf_sid cookie on all future requests
				sessionDataStr        := string(output.JWT_token_val)
				sessionTTLhoursInt, _ := gf_identity_core.GetSessionTTL()
				gf_session.SetOnReq(sessionDataStr, p_resp, sessionTTLhoursInt)

				//---------------------

				output_map := map[string]interface{}{
					"auth_signature_valid_bool": output.Auth_signature_valid_bool,
					"nonce_exists_bool":         output.Nonce_exists_bool,
					"user_id_str":               output.User_id_str,
				}
				return output_map, nil
			}


			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	// USERS_CREATE
	// NO_AUTH - unauthenticated users are able to create new users
	gf_rpc_lib.CreateHandlerHTTPwithAuth(false, "/v1/identity/eth/create",
		func(p_ctx context.Context, p_resp http.ResponseWriter, p_req *http.Request) (map[string]interface{}, *gf_core.GFerror) {

			if p_req.Method == "POST" {

				//---------------------
				// INPUT
				input_map, gf_err := gf_core.HTTPgetInput(p_req, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				input :=&GF_user_auth_eth__input_create{
					UserTypeStr:          "standard",
					User_address_eth_str: gf_identity_core.GF_user_address_eth(input_map["user_address_eth_str"].(string)),
					Auth_signature_str:   gf_identity_core.GF_auth_signature(input_map["auth_signature_str"].(string)),
				}
				
				//---------------------
				output, gf_err := users_auth_eth__pipeline__create(input, p_ctx, p_runtime_sys)
				if gf_err != nil {
					return nil, gf_err
				}

				output_map := map[string]interface{}{
					"auth_signature_valid_bool": output.Auth_signature_valid_bool,
					"nonce_exists_bool":         output.Nonce_exists_bool,
				}
				return output_map, nil
			}

			return nil, nil
		},
		rpc_handler_runtime,
		p_runtime_sys)

	//---------------------
	return nil
}