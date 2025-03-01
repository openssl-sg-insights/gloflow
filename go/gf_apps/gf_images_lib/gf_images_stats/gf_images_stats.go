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

package gf_images_stats

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
func Get_query_funs(p_runtime_sys *gf_core.RuntimeSys) map[string]func(*gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_images_stats.Init()")

	stats_funs_map := map[string]func(*gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {
		"image_jobs_errors":                  stats__image_jobs_errors,
		"completed_image_jobs_runtime_infos": stats__completed_image_jobs_runtime_infos,
	}
	return stats_funs_map
}

//-------------------------------------------------
func stats__image_jobs_errors(p_runtime_sys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_images_stats.stats__image_jobs_errors()")


	ctx := context.Background()

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{
				{"t",            "img_running_job"},
				{"start_time_f", bson.M{"$exists": true}}, // filter for new start_time_f/end_time_f format
			}},
		},
		{
			{"$project", bson.D{
				{"id_str",       true},
				{"errors_lst",   true},
				{"start_time_f", true}, // include field
				{"errors_num_i", bson.M{"$size": "$errors_lst",}},
			}},
		},
		{
			{"$sort", bson.D{
				{"start_time_f", 1},
			}},
		},
	}

	/*pipe := p_runtime_sys.Mongo_coll.Pipe([]bson.M{
		bson.M{"$match": bson.M{
				"t":            "img_running_job",
				"start_time_f": bson.M{"$exists": true}, // filter for new start_time_f/end_time_f format
			},
		},
		bson.M{"$project": bson.M{
				"id_str":       true,
				"errors_lst":   true,
				"start_time_f": true, // include field
				"errors_num_i": bson.M{"$size":"$errors_lst",},
			},
		},
		bson.M{"$sort":bson.M{
				"start_time_f":1,
			},
		},
	})*/

	cursor, err := p_runtime_sys.Mongo_coll.Aggregate(ctx, pipeline)
	if err != nil {

		gf_err := gf_core.MongoHandleError("failed to run DB aggregation to get img_running_jobs not complete yet",
			"mongodb_aggregation_error",
			nil,
			err, "gf_images_stats", p_runtime_sys)
		return nil, gf_err
	}
	defer cursor.Close(ctx)

	/*results_lst := []map[string]interface{}{}
	err         := pipe.All(&results_lst)

	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to get errors of all img_running_job's",
			"mongodb_aggregation_error",
			nil,err,"gf_images_stats",p_runtime_sys)
		return nil,gf_err
	}*/

	results_lst := []map[string]interface{}{}
	for cursor.Next(ctx) {

		var r map[string]interface{}
		err := cursor.Decode(&r)
		if err != nil {
			gf_err := gf_core.MongoHandleError("failed to decode mongodb result of image_jobs aggregation",
				"mongodb_cursor_decode",
				nil,
				err, "gf_images_stats", p_runtime_sys)
			return nil, gf_err
		}
	
		results_lst = append(results_lst, r)
	}

	data_map := map[string]interface{}{
		"image_jobs_errors_lst": results_lst,
	}

	return data_map,nil
}

//-------------------------------------------------
func stats__completed_image_jobs_runtime_infos(p_runtime_sys *gf_core.RuntimeSys) (map[string]interface{}, *gf_core.GFerror) {
	p_runtime_sys.LogFun("FUN_ENTER", "gf_images_stats.stats__completed_image_jobs_runtime_infos()")



	ctx := context.Background()
	pipeline := mongo.Pipeline{
		{
			// {"$match", bson.D{{"t", "peer_new_lifecycle"}}},
			{"$match", bson.D{
				{"t",            "img_running_job"},
				{"status_str",   "complete"},
				{"start_time_f", bson.M{"$exists": true}}, // filter for new start_time_f/end_time_f format
			}},
		},
		{
			{"$project", bson.D{
				{"_id",          false},
				{"status_str",   true},
				{"start_time_f", true}, // include field
				{"end_time_f",   true}, // include field

				{"runtime_duration_sec_f", bson.M{"$subtract": []string{"$end_time_f", "$start_time_f"},}},
				{"processed_imgs_num_i",   bson.M{"$size":     "$images_urls_to_process_lst",}},
			}},
		},
		{
			{"$sort", bson.D{
				{"start_time_f", 1},
			}},
		},
	}



	/*pipe := p_runtime_sys.Mongo_coll.Pipe([]bson.M{
		bson.M{"$match":bson.M{
				"t":            "img_running_job",
				"status_str":   "complete",
				"start_time_f": bson.M{"$exists": true}, // filter for new start_time_f/end_time_f format
			},
		},
		bson.M{"$project":bson.M{
				"_id":                    false, // suppression of the "_id" field
				"status_str":             true,  // include field
				"start_time_f":           true,  // include field
				"end_time_f":             true,  // include field
				"runtime_duration_sec_f": bson.M{"$subtract": []string{"$end_time_f", "$start_time_f"},},
				"processed_imgs_num_i":   bson.M{"$size":     "$images_urls_to_process_lst",},
			},
		},
		bson.M{"$sort":bson.M{"start_time_f":1},},
	})*/

	cursor, err := p_runtime_sys.Mongo_coll.Aggregate(ctx, pipeline)
	if err != nil {

		gf_err := gf_core.MongoHandleError("failed to run DB aggregation to get img_running_jobs that are complete",
			"mongodb_aggregation_error",
			map[string]interface{}{},
			err, "gf_images_stats", p_runtime_sys)
		return nil, gf_err
	}
	defer cursor.Close(ctx)

	/*results_lst := []map[string]interface{}{}
	err         := pipe.All(&results_lst)

	if err != nil {
		gf_err := gf_core.MongoHandleError("failed to run an aggregation pipeline to get runtime_info of all img_running_job's",
			"mongodb_aggregation_error",
			nil, err, "gf_images_stats", p_runtime_sys)
		return nil, gf_err
	}*/
	
	results_lst := []map[string]interface{}{}
	for cursor.Next(ctx) {

		var r map[string]interface{}
		err := cursor.Decode(&r)
		if err != nil {
			gf_err := gf_core.MongoHandleError("failed to decode mongodb result of image_jobs aggregation",
				"mongodb_cursor_decode",
				map[string]interface{}{},
				err, "gf_images_stats", p_runtime_sys)
			return nil, gf_err
		}
	
		results_lst = append(results_lst, r)
	}



	data_map := map[string]interface{}{
		"completed_image_jobs_runtime_infos_lst": results_lst,
	}

	return data_map,nil
}