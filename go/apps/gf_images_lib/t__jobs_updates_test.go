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
	"testing"
	"github.com/davecgh/go-spew/spew"
	"github.com/gloflow/gloflow/go/gf_core"
	"github.com/gloflow/gloflow/go/apps/gf_images_lib/gf_images_jobs"
)
//---------------------------------------------------
func Test__jobs_updates(p_test *testing.T) {
	
	//-----------------
	//TEST_DATA
	
	test__http_server_host_str             := "localhost:8000"
	test__mongodb_host_str                 := "127.0.0.1"
	test__mongodb_db_name_str              := "gf_tests"
	test__images_local_dir_path_str        := "./tests_data"
	test__images_thumbs_local_dir_path_str := "./tests_data/thumbnails"
	test__s3_bucket_name_str               := "gf--test--img"
	test__image_client_type_str            := "test_run"
	test__image_flows_names_lst            := []string{"test_flow",}
	test__image_url_str                    := fmt.Sprintf("http://%s/test_image_01.jpeg", test__http_server_host_str)
	test__origin_page_url_str              := "https://some_test_domain.com/page_1"

	fmt.Println(test__http_server_host_str)
	//-------------
	log_fun      := gf_core.Init_log_fun()
	mongodb_db   := gf_core.Mongo__connect(test__mongodb_host_str, test__mongodb_db_name_str, log_fun)
	mongodb_coll := mongodb_db.C("data_symphony")
	
	runtime_sys := &gf_core.Runtime_sys{
		Service_name_str:"gf_images_tests",
		Log_fun:         log_fun,
		Mongodb_coll:    mongodb_coll,
	}
	//-------------
	//S3
	s3_info := gf_core.T__get_s3_info(runtime_sys)
	//-------------
	//JOBS_MNGR
	jobs_mngr := gf_images_jobs.Jobs_mngr__init(test__images_local_dir_path_str,
		test__images_thumbs_local_dir_path_str,
		test__s3_bucket_name_str,
		s3_info,
		runtime_sys)
	//-------------

	test__job_start(test__image_url_str,
		test__image_flows_names_lst,
		test__image_client_type_str,
		test__origin_page_url_str,
		jobs_mngr,
		p_test,
		runtime_sys)
}
//---------------------------------------------------
func test__job_start(p_test__image_url_str string,
	p_test__image_flows_names_lst []string,
	p_test__image_client_type_str string,
	p_test__origin_page_url_str   string,
	p_jobs_mngr                   gf_images_jobs.Jobs_mngr,
	p_test                        *testing.T,
	p_runtime_sys                 *gf_core.Runtime_sys) {


	images_to_process_lst := []gf_images_jobs.Image_to_process{
		gf_images_jobs.Image_to_process{
			Source_url_str:     p_test__image_url_str,
			Origin_page_url_str:p_test__origin_page_url_str,
		},
	}

	running_job, output, gf_err := gf_images_jobs.Job__start(p_test__image_client_type_str,
		images_to_process_lst,
		p_test__image_flows_names_lst,
		p_jobs_mngr,
		p_runtime_sys)
	if gf_err != nil {
		panic(gf_err.Error)
	}


	fmt.Println(running_job)
	spew.Dump(output)



	//-------------
	//TEST_JOB_UPDATES
	test_job_id_str := running_job.Id_str
	job_updates_ch  := gf_images_jobs.Job__get_update_ch(test_job_id_str, p_jobs_mngr, p_runtime_sys)

	for ;; {

		fmt.Println("\n\n------------------------- TESTING - GET_JOB_UPDATE -----")
		job_update := <-job_updates_ch

		spew.Dump(job_update)

		job_update_type_str := job_update.Type_str
		if job_update_type_str == gf_images_jobs.JOB_UPDATE_TYPE__ERROR {
			panic("job encountered an error while processing")
		}

		if !(job_update_type_str == gf_images_jobs.JOB_UPDATE_TYPE__OK || job_update_type_str == gf_images_jobs.JOB_UPDATE_TYPE__COMPLETED) {
			panic(fmt.Sprintf("job_update is expected to be of type 'ok' but instead is - %s", job_update_type_str))
		}
		
		//test complete
		if job_update_type_str == gf_images_jobs.JOB_UPDATE_TYPE__COMPLETED {
			break
		}
	}
	//-------------
}