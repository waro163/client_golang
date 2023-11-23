package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

const (
	addressTemp = "https://prometheus.%s.k8s.tesla.com"
)

func main() {
	client, err := api.NewClient(api.Config{
		Address: fmt.Sprintf(addressTemp, "cn-pvg16-eng-general"),
		RoundTripper: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// api  query
	result, warnings, err := v1api.Query(ctx, "kube_resourcequota{namespace=\"scc-release-stage\"}", time.Now(), v1.WithTimeout(5*time.Second))
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	// fmt.Printf("Result:\n%#v\n", result)
	switch result.Type() {
	case model.ValVector:
		vecRes := result.(model.Vector)
		for _, sample := range vecRes {
			// fmt.Println(sample.String())
			fmt.Println("metric: ", sample.Metric.String())
			for key, value := range sample.Metric {
				fmt.Println("key: ", key, " value: ", value)
			}
			fmt.Println("value: ", sample.Value)
			fmt.Println("time: ", sample.Timestamp)
			// fmt.Println("Histogram: ", *sample.Histogram)
			fmt.Println("--------------------")
			break
		}
	default:
		fmt.Println("other type")
	}

	// api query range
	// r := v1.Range{
	// 	Start: time.Now().Add(-time.Hour),
	// 	End:   time.Now(),
	// 	Step:  time.Minute,
	// }
	// result, warnings, err = v1api.QueryRange(ctx, "rate(prometheus_tsdb_head_samples_appended_total[5m])", r, v1.WithTimeout(5*time.Second))
	// if err != nil {
	// 	fmt.Printf("Error querying Prometheus: %v\n", err)
	// 	os.Exit(1)
	// }
	// if len(warnings) > 0 {
	// 	fmt.Printf("Warnings: %v\n", warnings)
	// }
	// fmt.Printf("Result:\n%v\n", result)
}
