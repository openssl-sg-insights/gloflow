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

package gf_eth_monitor_core

import (
	"fmt"
	"github.com/gloflow/gloflow/go/gf_core"
)

//-------------------------------------------------
type GF_py_plugins struct {
	Base_dir_path_str string
}

//-------------------------------------------------
func py__run_plugin__get_contract_info(p_plugins_info *GF_py_plugins,
	p_runtime_sys *gf_core.Runtime_sys) *gf_core.Gf_error {


	py_path_str       := fmt.Sprintf("%s/gf_plugin__get_contract_info.py", p_plugins_info.Base_dir_path_str)
	stdout_prefix_str := "GF_OUT:"

	outputs_lst, gf_err := gf_core.CLI_py__run(py_path_str, stdout_prefix_str, p_runtime_sys)
	if gf_err != nil {
		return gf_err
	}


	for _, o := range outputs_lst {

		fmt.Println(o)

	}

	return nil
}