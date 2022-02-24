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

///<reference path="../../../d/jquery.d.ts" />



import * as gf_identity_eth      from "./gf_identity_eth";
import * as gf_identity_userpass from "./gf_identity_userpass";
import * as gf_identity_http     from "./gf_identity_http";
import * as gf_utils             from "./../../../gf_core/ts/gf_utils";

declare const window: any;
declare var Web3;

//-------------------------------------------------
export async function init_with_http(p_urls_map) {
    
    const http_api_map = get_http_api(p_urls_map);
    init(http_api_map);
}

//-------------------------------------------------
export async function init(p_http_api_map) {

    $("#identity #login").on("click", async function(p_e) {
        
        const opened_bool = $("#identity").has("#auth_pick_dialog").length > 0;

        if (!opened_bool) {
            
            const method_str = await auth_method_pick();
            switch (method_str) {
                //--------------------------
                // ETH_METAMASK
                case "eth_metamask":
                    await gf_identity_eth.user_auth_pipeline(p_http_api_map);
                    break;
                
                //--------------------------
                // USER_AND_PASS
                case "userpass":
                    await gf_identity_userpass.user_auth_pipeline(p_http_api_map);
                    break;

                //--------------------------
            }
        }
        else {
            $("#identity").find("#auth_pick_dialog").remove();
        }
    });
}

//-------------------------------------------------
async function auth_method_pick() {
    const p = new Promise(function(p_resolve_fun, p_reject_fun) {

        const container_identity = $("#identity");
        const container = $(`
        <div id="auth_pick_dialog">

            <div id="wallet_pick_dialog">
                <div id="metamask">
                    <div id="icon">
                        <img src="https://gloflow.com/images/static/assets/gf_metamask_icon.svg"></img>
                    </div>
                    <div id="descr">metamask browser wallet</div>
                </div>
            </div>

            <div id="user_and_pass_pick_dialog">
                <div id="label">user and password</div>
            </div>
        </div>`);

        container_identity.append(container);

        // close_dialog
        gf_utils.click_outside(container_identity, ()=>{
            $(container).remove();
        });

        $(container).find("#metamask").on('click', (p_e)=>{
            p_e.stopPropagation();

            // user picked to auth with wallet so remove initial p
            // auth method pick dialog.
            // $(container).remove();

            p_resolve_fun("eth_metamask");
        })

        $(container).find("#user_and_pass_pick_dialog").on('click', (p_e)=>{
            p_e.stopPropagation();

            // user picked to auth with username/pass so 
            // remove initial auth method pick dialog
            $(container).remove();

            p_resolve_fun("userpass");
        });
    });
    return p;
}

//-------------------------------------------------
export function get_standard_http_urls() {
    const login_url_str = '/v1/identity/users/login';
    const urls_map = {
        "login": login_url_str
    };
    return urls_map;
}

//-------------------------------------------------
export function get_admin_http_urls() {
    const login_url_str = '/v1/admin/login';
    const urls_map = {
        "login": login_url_str
    };
    return urls_map;
}

//-------------------------------------------------
export function get_http_api(p_urls_map) {
    const http_api_map = {

        // ETH
        "eth": {
            "user_preflight_fun": async (p_user_address_eth_str)=>{
                const output_map = await gf_identity_http.user_preflight(null, p_user_address_eth_str);
                return output_map;
            },
            "user_login_fun": async (p_user_address_eth_str, p_auth_signature_str)=>{
                
                const output_map = await gf_identity_http.user_eth_login(p_user_address_eth_str,
                    p_auth_signature_str);
                return output_map;
            },
            "user_create_fun": async (p_user_address_eth_str, p_auth_signature_str)=>{
                const output_map = await gf_identity_http.user_eth_create(p_user_address_eth_str, p_auth_signature_str);
                return output_map;
            }
        },

        // USERPASS
        "userpass": {
            "user_login_fun": async (p_user_name_str, p_pass_str, p_email_str)=>{
                const url_str    = p_urls_map["login"];
                const output_map = await gf_identity_http.user_userpass_login(p_user_name_str,
                    p_pass_str,
                    p_email_str,
                    url_str);
                return output_map;
            },
            "user_create_fun": async (p_user_name_str)=>{
                const output_map = {};
                return output_map;
            }
        },

        // MFA
        "mfa": {
            "user_mfa_confirm": async (p_user_name_str, p_mfa_val_str)=>{
                const output_map = await gf_identity_http.user_mfa_confirm(p_user_name_str,
                    p_mfa_val_str);
                return output_map;
            }
        }
    };
    return http_api_map;
}