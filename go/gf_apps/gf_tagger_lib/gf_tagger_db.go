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

package gf_tagger_lib

import (
	"fmt"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/gf_apps/gf_publisher_lib/gf_publisher_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
// VAR
//---------------------------------------------------
func db__get_objects_with_tag_count(p_tag_str string,
	p_object_type_str string,
	p_runtime_sys     *gf_core.RuntimeSys) (int64, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_tagger_db.db__get_objects_with_tag_count()")

	switch p_object_type_str {
		case "post":

			ctx := context.Background()
			count_int, err := p_runtime_sys.Mongo_coll.CountDocuments(ctx, bson.M{
				"t":        "post",
				"tags_lst": bson.M{"$in": []string{p_tag_str,}},
			})

			if err != nil {
				gf_err := gf_core.MongoHandleError(fmt.Sprintf("failed to count of posts with tag - %s in DB", p_tag_str),
					"mongodb_find_error",
					map[string]interface{}{
						"tag_str":         p_tag_str,
						"object_type_str": p_object_type_str,
					},
					err, "gf_tagger_lib", p_runtime_sys)
				return 0, gf_err
			}
			return count_int, nil
	}
	return 0, nil
}

//---------------------------------------------------
// POSTS
//---------------------------------------------------
func db__get_post_notes(p_post_title_str string,
	p_runtime_sys *gf_core.RuntimeSys) ([]*GF_note, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_tagger_db.db__get_post_notes()")

	post, gf_err := gf_publisher_core.DB__get_post(p_post_title_str, p_runtime_sys)
	if gf_err != nil {
		return nil, gf_err
	}

	post_notes_lst := post.Notes_lst
	notes_lst      := []*GF_note{}
	for _,s := range post_notes_lst {

		note := &GF_note{
			User_id_str:           s.User_id_str,
			Body_str:              s.Body_str,
			Target_obj_id_str:     post.Title_str,
			Target_obj_type_str:   "post",
			Creation_datetime_str: s.Creation_datetime_str,
		}
		notes_lst = append(notes_lst,note)
	}
	p_runtime_sys.LogFun("INFO", "got # notes - "+fmt.Sprint(len(notes_lst)))
	return notes_lst, nil
}

//---------------------------------------------------
func db__add_post_note(p_note *GF_note,
	p_post_title_str string,
	p_runtime_sys    *gf_core.RuntimeSys) *gf_core.GFerror {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_tagger_db.db__add_post_note()")

	//--------------------
	post_note := &gf_publisher_core.Gf_post_note{
		User_id_str:           p_note.User_id_str,
		Body_str:              p_note.Body_str,
		Creation_datetime_str: p_note.Creation_datetime_str,
	}

	//--------------------
	
	ctx := context.Background()
	_, err := p_runtime_sys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":         "post",
			"title_str": p_post_title_str,
		}, 
		bson.M{"$push": bson.M{"notes_lst": post_note},
	})
	
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to update a gf_post in a mongodb with a new note in DB",
			"mongodb_update_error",
			map[string]interface{}{
				"post_title_str": p_post_title_str,
				"note":           p_note,
			},
			err, "gf_tagger_lib", p_runtime_sys)
		return gf_err
	}
	return nil
}

//---------------------------------------------------
func db__get_posts_with_tag(p_tag_str string,
	p_page_index_int int,
	p_page_size_int  int,
	p_runtime_sys    *gf_core.RuntimeSys) ([]*gf_publisher_core.Gf_post, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_tagger_db.db__get_posts_with_tag()")
	p_runtime_sys.LogFun("INFO",      fmt.Sprintf("p_tag_str - %s", p_tag_str))

	// FIX!! - potentially DOESNT SCALE. if there is a huge number of posts
	//         with a tag, toList() will accumulate a large collection in memory. 
	//         instead use a Stream-oriented way where results are streamed lazily
	//         in some fashion
		
	



	ctx := context.Background()

	find_opts := options.Find()
	find_opts.SetSort(map[string]interface{}{"creation_datetime": -1})
	find_opts.SetSkip(int64(p_page_index_int))
    find_opts.SetLimit(int64(p_page_size_int))

	cursor, gf_err := gf_core.MongoFind(bson.M{
			"t":        "post",
			"tags_lst": bson.M{"$in": []string{p_tag_str,}},
		},
		find_opts,
		map[string]interface{}{
			"tag_str":            p_tag_str,
			"page_index_int":     p_page_index_int,
			"page_size_int":      p_page_size_int,
			"caller_err_msg_str": fmt.Sprintf("failed to get posts with specified tag in DB"),
		},
		p_runtime_sys.Mongo_coll,
		ctx,
		p_runtime_sys)

	if gf_err != nil {
		return nil, gf_err
	}

	var posts_lst []*gf_publisher_core.Gf_post
	err := cursor.All(ctx, &posts_lst)
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to get posts with specified tag in DB",
			"mongodb_cursor_decode",
			map[string]interface{}{
				"tag_str":        p_tag_str,
				"page_index_int": p_page_index_int,
				"page_size_int":  p_page_size_int,
			},
			err, "gf_tagger_lib", p_runtime_sys)
		return nil, gf_err
	}

	/*err := p_runtime_sys.Mongodb_coll.Find(bson.M{
			"t":        "post",
			"tags_lst": bson.M{"$in": []string{p_tag_str,}},
		}).
		Sort("-creation_datetime"). // descending:true
		Skip(p_page_index_int).
		Limit(p_page_size_int).
		All(&posts_lst)

	if gf_err != nil {
		gf_err := gf_core.MongoHandleError("failed to get posts with specified tag",
			"mongodb_find_error",
			map[string]interface{}{
				"tag_str":        p_tag_str,
				"page_index_int": p_page_index_int,
				"page_size_int":  p_page_size_int,
			},
			err, "gf_tagger", p_runtime_sys)
		return nil, gf_err
	}*/

	return posts_lst, nil
}

//---------------------------------------------------
func db__add_tags_to_post(p_post_title_str string,
	p_tags_lst    []string,
	p_runtime_sys *gf_core.RuntimeSys) *gf_core.GFerror {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_tagger_db.db__add_tags_to_post()")

	ctx := context.Background()
	_, err := p_runtime_sys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":         "post",
			"title_str": p_post_title_str,
		},
		bson.M{"$push": bson.M{"tags_lst": p_tags_lst},
	})
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to update a gf_post with new tags in DB",
			"mongodb_update_error",
			map[string]interface{}{
				"post_title_str": p_post_title_str,
				"tags_lst":       p_tags_lst,
			},
			err, "gf_tagger_lib", p_runtime_sys)
		return gf_err
	}
	return nil
}

//---------------------------------------------------
// IMAGES
//---------------------------------------------------
func db__add_tags_to_image(pImageIDstr string,
	pTagsLst    []string,
	pRuntimeSys *gf_core.RuntimeSys) *gf_core.GFerror {

	ctx := context.Background()
	_, err := pRuntimeSys.Mongo_coll.UpdateMany(ctx, bson.M{
			"t":      "img",
			"id_str": pImageIDstr,
		},
		bson.M{"$push": bson.M{"tags_lst": pTagsLst},
	})
	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to update a gf_image with new tags in DB",
			"mongodb_update_error",
			map[string]interface{}{
				"image_id_str": pImageIDstr,
				"tags_lst":     pTagsLst,
			},
			err, "gf_tagger_lib", pRuntimeSys)
		return gf_err
	}
	return nil
}