/*
GloFlow application and media management/publishing platform
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

package gf_images_utils

import (
	// "fmt"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	// "github.com/globalsign/mgo/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//---------------------------------------------------
func DB__put_image(p_image *Gf_image,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_images_db.DB__put_image()")
	
	ctx := context.Background()

	// UPSERT
	query  := bson.M{"t": "img", "id_str": p_image.Id_str,}
	gf_err := gf_core.Mongo__upsert(query,
		p_image,
		map[string]interface{}{"image_id_str": p_image.Id_str,},
		p_runtime_sys.Mongo_coll,
		ctx, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}

	/*// spec          - a dict specifying elements which must be present for a document to be updated
	// upsert = True - insert doc if it doesnt exist, else just update
	_, err := p_runtime_sys.Mongo_coll.Upsert(bson.M{"t": "img", "id_str": p_image.Id_str,}, p_image)
	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to update/upsert gf_image in a mongodb",
			"mongodb_update_error",
			map[string]interface{}{"image_id_str": p_image.Id_str,},
			err, "gf_images_utils", p_runtime_sys)
		return gf_err
	}*/

	return nil
}

//---------------------------------------------------
func DB__get_image(p_image_id_str Gf_image_id,
	p_runtime_sys *gf_core.Runtime_sys) (*Gf_image, *gf_core.Gf_error) {
	p_runtime_sys.Log_fun("FUN_ENTER", "gf_image_db.DB__get_image()")
	


	ctx := context.Background()
	var image Gf_image

	q             := bson.M{"t": "img", "id_str": p_image_id_str}
	coll_name_str := p_runtime_sys.Mongo_coll.Name()
	err           := p_runtime_sys.Mongo_db.Collection(coll_name_str).FindOne(ctx, q).Decode(&image)
	if err != nil {

		// FIX!! - a record not being found in the DB is possible valid state. it should be considered
		//         if this should not return an error but instead just a "nil" value for the record.
		if err == mongo.ErrNoDocuments {
			gf_err := gf_core.Mongo__handle_error("image does not exist in mongodb",
				"mongodb_not_found_error",
				map[string]interface{}{"image_id_str": p_image_id_str,},
				err, "gf_images_utils", p_runtime_sys)
			return nil, gf_err
		}
		
		gf_err := gf_core.Mongo__handle_error("failed to get image from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"image_id_str": p_image_id_str,},
			err, "gf_images_utils", p_runtime_sys)
		return nil, gf_err
	}


	/*var image Gf_image
	err := p_runtime_sys.Mongo_coll.Find(bson.M{"t": "img", "id_str": p_image_id_str}).One(&image)

	// FIX!! - a record not being found in the DB is possible valid state. it should be considered
	//         if this should not return an error but instead just a "nil" value for the record.
	if fmt.Sprint(err) == "not found" {
		gf_err := gf_core.Mongo__handle_error("image does not exist in mongodb",
			"mongodb_not_found_error",
			map[string]interface{}{"image_id_str": p_image_id_str,},
			err, "gf_images_utils", p_runtime_sys)
		return nil, gf_err
	}

	if err != nil {
		gf_err := gf_core.Mongo__handle_error("failed to get image from mongodb",
			"mongodb_find_error",
			map[string]interface{}{"image_id_str": p_image_id_str,},
			err, "gf_images_utils", p_runtime_sys)
		return nil, gf_err
	}*/
	
	return &image, nil
}