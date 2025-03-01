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

package gf_home_lib

import (
	"text/template"
	"context"
	"github.com/gloflow/gloflow/go/gf_core"
)

//------------------------------------------------
// RENDER_DASHBOARD
func PipelineRenderDashboard(pTmpl *template.Template,
	pSubtemplatesNamesLst []string,
	pCtx                  context.Context,
	pRuntimeSys           *gf_core.RuntimeSys) (string, *gf_core.GFerror) {

	templateRenderedStr, gfErr := viewRenderTemplateDashboard(pTmpl,
		pSubtemplatesNamesLst,
		pRuntimeSys)
	if gfErr != nil {
		return "", gfErr
	}

	return templateRenderedStr, nil
}

//---------------------------------------------------
func homeVizCreateID(pUserIdentifierStr string,
	pCreationUNIXtimeF float64) gf_core.GF_ID {

	fieldsForIDlst := []string{
		pUserIdentifierStr,
	}
	gfIDstr := gf_core.IDcreate(fieldsForIDlst,
		pCreationUNIXtimeF)

	return gfIDstr
}