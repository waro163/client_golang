// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A minimal example of how to include Prometheus instrumentation.
package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	flag.Parse()
	rand.Seed(time.Now().Unix())

	httpReqs := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"code", "method"},
	)
	prometheus.MustRegister(httpReqs)

	httpReqs.WithLabelValues("404", "POST").Add(42)
	go func() {
		for {
			v := rand.Intn(10)
			if v%2 == 0 {
				httpReqs.With(prometheus.Labels{"code": "200", "method": "POST"}).Add(float64(v))
			} else {
				httpReqs.With(prometheus.Labels{"code": "200", "method": "GET"}).Inc()
			}
			if v%3 == 0 {
				httpReqs.With(prometheus.Labels{"code": "200", "method": "post"}).Add(float64(v))
			} else {
				httpReqs.With(prometheus.Labels{"code": "200", "method": "gGet"}).Inc()
			}
			time.Sleep(time.Microsecond * 300)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
