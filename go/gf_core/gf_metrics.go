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

package gf_core

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

//-------------------------------------------------
type GFmetrics struct {

	// ERRORS_COUNTER - number of GF errors that were thrown
	ErrorsCounter prometheus.Counter
}

//-------------------------------------------------
// INIT
func MetricsInit(p_metrics_endpoint_str string, // "/metrics"
	pPortInt int) *GFmetrics {
	
	// ERRORS_COUNTER
	errorsCounter := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "gf_core__errors_num",
		Help: "number of gf_errors thrown in the system",
	})
	prometheus.MustRegister(errorsCounter)
	

	// START_SERVER
	go func() {
		metrics_router := mux.NewRouter()
		metrics_router.Handle(p_metrics_endpoint_str, promhttp.Handler())


		metricsServer := http.Server{
			Handler: metrics_router,
			Addr:    fmt.Sprintf(":%d", pPortInt),
		}
		
		// ADD!! - check for returned error here,
		//         and report this in some way to the user.
		metricsServer.ListenAndServe()
	}()

	metrics := &GFmetrics{
		ErrorsCounter: errorsCounter,
	}

	return metrics
}