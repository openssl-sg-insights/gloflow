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
	"time"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
type GF_login_attempt struct {
	V_str                string             `bson:"v_str"` // schema_version
	Id                   primitive.ObjectID `bson:"_id,omitempty"`
	Id_str               gf_core.GF_ID      `bson:"id_str"`
	Deleted_bool         bool               `bson:"deleted_bool"`
	Creation_unix_time_f float64            `bson:"creation_unix_time_f"`

	User_type_str        string       `bson:"user_type_str"` // "regular"|"admin"
	User_name_str        GF_user_name `bson:"user_name_str"`
	
	Pass_confirmed_bool  bool `bson:"pass_confirmed_bool"`
	Email_confirmed_bool bool `bson:"email_confirmed_bool"`
	MFA_confirmed_bool   bool `bson:"mfa_confirmed_bool"`
}

//---------------------------------------------------
func login_attempt__get_if_valid(p_user_name_str GF_user_name,
	p_ctx         context.Context,
	p_runtime_sys *gf_core.Runtime_sys) (*GF_login_attempt, *gf_core.GF_error) {

	login_attempt_max_age_seconds_f := 5*60.0
	
	var login_attempt *GF_login_attempt
	login_attempt, gf_err := db__login_attempt__get_by_username(p_user_name_str,
		p_ctx,
		p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	if login_attempt == nil {
		return nil, nil
	}


	// get the age of the login_attempt
	current_unix_time_f := float64(time.Now().UnixNano())/1000000000.0
	age_f               := current_unix_time_f - login_attempt.Creation_unix_time_f


	// login_attempt has expired
	if age_f > login_attempt_max_age_seconds_f {

		// mark it as deleted
		expired_bool := true
		update_op := &GF_login_attempt__update_op{Deleted_bool: &expired_bool}
		gf_err = db__login_attempt__update(&login_attempt.Id_str,
			update_op,
			p_ctx,
			p_runtime_sys)
		if gf_err != nil {
			return nil, gf_err
		}

		return nil, nil
	} else {
		return login_attempt, nil
	}

	return nil, nil
}