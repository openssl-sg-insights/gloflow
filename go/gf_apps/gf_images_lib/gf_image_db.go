/*
GloFlow media management/publishing system
Copyright (C) 2019 Ivan Trajkovic

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

package gf_images_lib

import (
	"fmt"
	"time"
	"math/rand"
	"github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_images_lib/gf_images_utils"
)
//---------------------------------------------------
func DB__get_random_imgs_range(p_imgs_num_to_get_int int, //5
	p_max_random_cursor_position_int int, //2000
	p_flow_name_str                  string,
	p_runtime_sys                    *gf_core.Runtime_sys) ([]*gf_images_utils.Gf_image,*gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER","gf_image_db.DB__get_random_imgs_range()")

	rand.Seed(time.Now().Unix())
	random_cursor_position_int := rand.Intn(p_max_random_cursor_position_int) //new Random().nextInt(p_max_random_cursor_position_int)
	p_runtime_sys.Log_fun("INFO","random_cursor_position_int - "+fmt.Sprint(random_cursor_position_int))

	var imgs_lst []*gf_images_utils.Gf_image
	err := p_runtime_sys.Mongodb_coll.Find(bson.M{
			"t"                   :"img",
			"creation_unix_time_f":bson.M{"$exists":true,},
			"flows_names_lst":     bson.M{"$in":[]string{p_flow_name_str},},
			//---------------------
			//IMPORTANT!! - this is the new member that indicates which page url (if not directly uploaded) the
			//              image came from. only use these images, since only they can be properly credited
			//              to the source site
			"origin_page_url_str" :bson.M{"$exists":true,},
			//---------------------
		}).
		Skip(random_cursor_position_int).
		Limit(p_imgs_num_to_get_int).
		All(&imgs_lst)
		
	if err != nil {
		gf_err := gf_core.Error__create("failed to get random img range from the DB",
			"mongodb_find_error",
			&map[string]interface{}{
				"imgs_num_to_get_int":           p_imgs_num_to_get_int,
				"max_random_cursor_position_int":p_max_random_cursor_position_int,
				"flow_name_str":                 p_flow_name_str,
			},
			err, "gf_images_lib", p_runtime_sys)
		return nil, gf_err
	}

	return imgs_lst, nil
}
//---------------------------------------------------
func DB__image_exists(p_image_id_str string, p_runtime_sys *gf_core.Runtime_sys) (bool, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_image_db.DB__image_exists()")

	c,err := p_runtime_sys.Mongodb_coll.Find(bson.M{"t":"img","id_str":p_image_id_str}).Count()
	if err != nil {
		gf_err := gf_core.Error__create("failed to check if image exists in the DB",
			"mongodb_find_error",
			&map[string]interface{}{"image_id_str":p_image_id_str,},
			err, "gf_images_lib", p_runtime_sys)
		return false, gf_err
	}

	if c > 0 {
		return true, nil
	} else {
		return false, nil
	}
}