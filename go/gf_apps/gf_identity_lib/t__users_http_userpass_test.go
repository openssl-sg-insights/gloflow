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
	"fmt"
	"testing"
	"context"
	"strings"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/parnurzeal/gorequest"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------
func Test__users_http_userpass(p_test *testing.T) {

	fmt.Println(" TEST__IDENTITY_USERS_HTTP_USERPASS >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	test_port_int := 2000
	ctx           := context.Background()
	runtime_sys   := T__init()
	http_agent    := gorequest.New()


	

	test__user_name_str := "ivan_t"
	test__user_pass_str := "pass_lksjds;lkdj"
	test__email_str     := "ivan_t@gloflow.com"



	// CLEANUP
	coll_name_str := "gf_users"
	gf_core.Mongo__delete(bson.M{}, coll_name_str, 
		map[string]interface{}{
			"caller_err_msg_str": "failed to cleanup test user DB",
		},
		ctx, runtime_sys)

	db__user__add_to_invite_list(test__email_str,
		ctx,
		runtime_sys)
	
	//---------------------------------
	// TEST_USER_CREATE_HTTP

	fmt.Println("====================================")
	fmt.Println("test user CREATE USERPASS")
	fmt.Println("user_name_str", test__user_name_str)
	fmt.Println("pass_str",      test__user_pass_str)
	fmt.Println("email_str",     test__email_str)

	url_str := fmt.Sprintf("http://localhost:%d/v1/identity/userpass/create", test_port_int)
	data_map := map[string]string{
		"user_name_str": test__user_name_str,
		"pass_str":      test__user_pass_str,
		"email_str":     test__email_str,
	}
	data_bytes_lst, _ := json.Marshal(data_map)
	_, body_str, errs := http_agent.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	spew.Dump(body_str)

	if (len(errs) > 0) {
		fmt.Println(errs)
		p_test.FailNow()
	}

	body_map := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
        p_test.FailNow()
    }

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user create http request failed")

	user_exists_bool         := body_map["data"].(map[string]interface{})["user_exists_bool"].(bool)
	user_in_invite_list_bool := body_map["data"].(map[string]interface{})["user_in_invite_list_bool"].(bool)

	if (user_exists_bool) {
		fmt.Println("supplied user already exists and cant be created")
		p_test.FailNow()
	}
	if (!user_in_invite_list_bool) {
		fmt.Println("supplied user is not in the invite list")
		p_test.FailNow()
	}

	//---------------------------------
	// TEST_USER_LOGIN

	fmt.Println("====================================")
	fmt.Println("test user LOGIN USERPASS")

	url_str = fmt.Sprintf("http://localhost:%d/v1/identity/userpass/login", test_port_int)
	data_map = map[string]string{
		"user_name_str": test__user_name_str,
		"pass_str":      test__user_pass_str,
	}
	data_bytes_lst, _ = json.Marshal(data_map)
	resp, body_str, errs := http_agent.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	if (len(errs) > 0) {
		fmt.Println(errs)
		p_test.FailNow()
	}

	// check if the login response sets a cookie for all future auth requests
	auth_cookie_present_bool := false
	for k, v := range resp.Header {
		if (k == "Set-Cookie") {
			for _, vv := range v {
				o := strings.Split(vv, "=")[0]
				if o == "gf_sess_data" {
					auth_cookie_present_bool = true
				}
			}
		}
	}
	assert.True(p_test, auth_cookie_present_bool,
		"login response does not contain the expected 'gf_sess_data' cookie")

	body_map = map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
        p_test.FailNow()
    }

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user login http request failed")

	user_exists_bool = body_map["data"].(map[string]interface{})["user_exists_bool"].(bool)
	pass_valid_bool := body_map["data"].(map[string]interface{})["pass_valid_bool"].(bool)
	user_id_str     := body_map["data"].(map[string]interface{})["user_id_str"].(string)

	assert.True(p_test, user_id_str != "", "user_id not set in the response")

	fmt.Println("user login response:")
	fmt.Println("user_exists_bool", user_exists_bool)
	fmt.Println("pass_valid_bool",  pass_valid_bool)
	fmt.Println("user_id_str",      user_id_str)

	//---------------------------------
	// TEST_USER_UPDATE

	fmt.Println("====================================")
	fmt.Println("test user UPDATE")
	fmt.Println("user inputs:")
	fmt.Println("user_id_str", user_id_str)

	url_str = fmt.Sprintf("http://localhost:%d/v1/identity/update", test_port_int)
	data_map = map[string]string{
		"user_name_str":   "ivan_t_new",
		"email_str":       "ivan_t_new@gloflow.com",
		"description_str": "some new description",
	}
	data_bytes_lst, _ = json.Marshal(data_map)
	_, body_str, errs = http_agent.Post(url_str).
		Send(string(data_bytes_lst)).
		End()

	body_map = map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
		p_test.FailNow()
	}

	spew.Dump(body_map)

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user updating http request failed")

	//---------------------------------
	// TEST_USER_GET_ME

	fmt.Println("====================================")
	fmt.Println("test user GET ME")
	
	url_str = fmt.Sprintf("http://localhost:%d/v1/identity/me", test_port_int)
	data_bytes_lst, _ = json.Marshal(data_map)
	_, body_str, errs = http_agent.Get(url_str).
		End()

	body_map = map[string]interface{}{}
	if err := json.Unmarshal([]byte(body_str), &body_map); err != nil {
		fmt.Println(err)
        p_test.FailNow()
    }

	assert.True(p_test, body_map["status"].(string) != "ERROR", "user get me http request failed")

	user_name_str         := body_map["data"].(map[string]interface{})["user_name_str"].(string)
	email_str             := body_map["data"].(map[string]interface{})["email_str"].(string)
	description_str       := body_map["data"].(map[string]interface{})["description_str"].(string)
	profile_image_url_str := body_map["data"].(map[string]interface{})["profile_image_url_str"].(string)
	banner_image_url_str  := body_map["data"].(map[string]interface{})["banner_image_url_str"].(string)

	fmt.Println("====================================")
	fmt.Println("user login response:")
	fmt.Println("user_name_str",         user_name_str)
	fmt.Println("email_str",             email_str)
	fmt.Println("description_str",       description_str)
	fmt.Println("profile_image_url_str", profile_image_url_str)
	fmt.Println("banner_image_url_str",  banner_image_url_str)

	//---------------------------------
}