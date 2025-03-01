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
	"time"
	"bytes"
	"strings"
	"strconv"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_identity_lib/gf_identity_core"
)

//------------------------------------------------
type GF_user_auth_mfa__input_confirm struct {
	UserNameStr           gf_identity_core.GFuserName `validate:"required,min=3,max=50"`
	Extern_htop_value_str string    `validate:"required,min=10,max=200"`
	Secret_key_base32_str string    `validate:"required,min=8,max=200"`
}

//------------------------------------------------
func mfaPipelineConfirm(pInput *GF_user_auth_mfa__input_confirm,
	pCtx        context.Context,
	pRuntimeSys *gf_core.RuntimeSys) (bool, *gf_core.GFerror) {




	htopValueStr, gfErr := TOTPgenerateValue(pInput.Secret_key_base32_str,
		pRuntimeSys)
	if gfErr != nil {
		return false, gfErr
	}




	if pInput.Extern_htop_value_str == htopValueStr {

		

		// get a preexisting login_attempt if one exists and hasnt expired for this user.
		// if it has then a new one will have to be created.
		var loginAttempt *GF_login_attempt
		loginAttempt, gfErr = loginAttemptGetIfValid(pInput.UserNameStr,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}

		// there is no login_attempt for this user thats active, or it has expired.
		// do nothing, forcing the user to restart the login process.
		if loginAttempt == nil {
			return false, nil
		}

		//------------------------
		// UPDATE_LOGIN_ATTEMPT
		// if password is valid then update the login_attempt 
		// to indicate that the password has been confirmed
		mfaConfirmBool := true
		updateOp := &GF_login_attempt__update_op{MFA_confirmed_bool: &mfaConfirmBool}
		gfErr = db__login_attempt__update(&loginAttempt.Id_str,
			updateOp,
			pCtx,
			pRuntimeSys)
		if gfErr != nil {
			return false, gfErr
		}

		//------------------------
		return true, nil
	} else {
		return false, nil
	}
	return false, nil
}

//---------------------------------------------------
// TOTP - https://datatracker.ietf.org/doc/html/rfc6238

func TOTPgenerateValue(pSecretKeyStr string,
	pRuntimeSys *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// index number of a 30s time period since start of Unix time.
	// TOTP is specified to use UTC time, so timezones dont have to be 
	// accounted for between multiple parties.
	intervalInt := time.Now().UTC().Unix() / 30

	fmt.Println("interval", intervalInt)
	tokenStr, gfErr  := HOTPgenerateValue(pSecretKeyStr,
		intervalInt,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return tokenStr, nil
}

//---------------------------------------------------
// HOTP_GENERATE_TOKEN
// HOTP - https://datatracker.ietf.org/doc/html/rfc4226

// secret_key - symmetric key

func HOTPgenerateValue(p_secret_key_base32_str string,
	p_time_interval_int int64,
	p_runtime_sys       *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	// expected length, can be 6-10 long, google-authenticator uses 6
	hotp_token_length_int := 6

	// convert secret_key to base32.
	// one way to represent Base32 numbers in a human-readable way is by using 
	// a standard 32-character set, such as the twenty-two upper-case letters A–V and the digits 0-9.
	// its done in order to allow for this limited alphabet.
	// first switch all characters to upper case just in case.
	key_uppercase_base32_str := strings.ToUpper(p_secret_key_base32_str)
	key_bytes_lst, err       := base32.StdEncoding.DecodeString(key_uppercase_base32_str)
	if err != nil {
		gf_err := gf_core.ErrorCreate("failed to base32 encode secret key for HOTP token generation",
			"generic_error",
			map[string]interface{}{
				"time_interval_int": p_time_interval_int,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return "", gf_err
	}

	// BIG_ENDIAN - byte order where the highest-order bit is stored first (lower memory address),
	//              and least last (higher address).
	// COUNTER - time_interval is the counter in the HOTP standard, defined to be 8 bytes.
	time_interval_bytes_lst := make([]byte, 8)
	binary.BigEndian.PutUint64(time_interval_bytes_lst, uint64(p_time_interval_int))

	// sign secret_key using HMAC-SHA1
	// SHA1 - 160-bit string (20 bytes), 40 digits in hex
	// HMAC - HMAC is a specific type of message authentication code involving a cryptographic
	//        hash function and a secret cryptographic key. As with any MAC, it may be used to
	//        simultaneously verify both the data integrity and authenticity of a message.
	//        trades off the need for a complex public key infrastructure by delegating the 
	//        key exchange to the communicating parties, who are responsible for establishing 
	//        and using a trusted channel to agree on the key prior to communication.
	// HMAC-SHA-1(K,C) - this is a hash that is a function of the key and counter (interval)
	hash := hmac.New(sha1.New, key_bytes_lst)
	hash.Write(time_interval_bytes_lst)
	hash_bytes_lst := hash.Sum(nil)




	//-----------------------
	// TRUNCATE - using a subset of the hash
	// use last half-byte to choose the index in the hash to start from.
	// this number can be a maximum of 15 (0xf), whereas sha1 hash is 20 bytes,
	// which leaves exactly 4 bytes which is needed.
	offset_int := (hash_bytes_lst[len(hash_bytes_lst)-1] & 0xf)



	


	var header_int uint32
	
	// generate 4 byte (32 bit) chunk from hash, starting at offset
	chunk_4_bytes_lst := hash_bytes_lst[offset_int : offset_int + 4]
	reader            := bytes.NewReader(chunk_4_bytes_lst)
	err = binary.Read(reader, binary.BigEndian, &header_int)
	if err != nil {
		gf_err := gf_core.ErrorCreate("failed to base32 encode secret key for HOTP token generation",
			"io_reader_error",
			map[string]interface{}{
				"time_interval_int": p_time_interval_int,
			},
			err, "gf_identity_lib", p_runtime_sys)
		return "", gf_err
	}

	//-----------------------

	/*code := int64(chunk_4_bytes_lst[0])<<24 |
		int64(chunk_4_bytes_lst[1])<<16 |
		int64(chunk_4_bytes_lst[2])<<8 |
		int64(chunk_4_bytes_lst[3])*/

	
	// 1_000_000 - is 10^d - "d" represents a "d" number of least significant digits,
	//             which are selected out here for usage.
	//             (ignoring the most significant bits)
	htop_value := (int(header_int) & 0x7fffffff) % 1_000_000

	htop_value_str := strconv.Itoa(int(htop_value))

	htop_value_padded_str := padWith0s(htop_value_str, hotp_token_length_int)
	return htop_value_padded_str, nil
}

//---------------------------------------------------
// PAD_WITH_0S - append/pad string with 0's until its of target_length

func padWith0s(p_str string, p_target_length_int int) string {
	if len(p_str) >= p_target_length_int {
		return p_str
	}
	for i := (p_target_length_int - len(p_str)); i > 0; i-- {
		p_str = "0" + p_str
	}
	return p_str
}