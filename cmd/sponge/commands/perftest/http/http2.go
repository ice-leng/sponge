package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/cmd/sponge/commands/perftest/common"
)

// PerfTestHTTP2CMD creates a new cobra.Command for HTTP/2 performance test.
func PerfTestHTTP2CMD() *cobra.Command {
	var (
		targetURL string
		method    string
		body      string
		bodyFile  string
		headers   []string

		worker   int
		total    uint64
		duration time.Duration

		out               string
		pushURL           string
		prometheusJobName string
	)

	cmd := &cobra.Command{
		Use:   "http2",
		Short: "Run performance test for HTTP/2 API",
		Long:  "Run performance test for HTTP/2 API.",
		Example: color.HiBlackString(`  # Default mode: worker=CPU*3, 5000 requests, GET request
  sponge perftest http2 --url=https://localhost:6443/user/1

  # Fixed number of requests: 50 workers, 500k requests, GET request
  sponge perftest http2 --worker=50 --total=500000 --url=https://localhost:6443/user/1

  # Fixed number of requests: 50 workers, 500k requests, POST request with JSON body
  sponge perftest http2 --worker=50 --total=500000 --url=https://localhost:6443/user --method=POST --body={\"name\":\"Alice\",\"age\":25}

  # Fixed duration: 50 workers, duration 10s, GET request
  sponge perftest http2 --worker=50 --duration=10s --url=https://localhost:6443/user/1

  # Fixed duration: 50 workers, duration 10s, POST request with JSON body
  sponge perftest http2 --worker=50 --duration=10s --url=https://localhost:6443/user --method=POST --body={\"name\":\"Alice\",\"age\":25}

  # Fixed number of requests: 50 workers, 500k requests, GET request, push statistics to custom HTTP server every 1s
  sponge perftest http2 --worker=50 --total=500000 --url=https://localhost:6443/user/1 --push-url=http://localhost:7070/report

  # Fixed duration: 50 workers, duration 10s, get request, push statistics to Prometheus (job=xxx)
  sponge perftest http2 --worker=50 --duration=10s --url=https://localhost:6443/user/1 --push-url=http://localhost:9091/metrics --prometheus-job-name=perftest-http2`),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			bodyBytes, headerMap, err := common.ParseHTTPParams(method, headers, body, bodyFile)
			if err != nil {
				return err
			}

			params := &HTTPReqParams{
				URL:     targetURL,
				Method:  method,
				Headers: headerMap,
				Body:    bodyBytes,
				version: "HTTP/2",
			}

			b := PerfTestHTTP{
				ID:                common.NewID(),
				Client:            newHTTP2Client(),
				Params:            params,
				Worker:            worker,
				TotalRequests:     total,
				Duration:          duration,
				PushURL:           pushURL,
				PrometheusJobName: prometheusJobName,
			}
			if err = b.checkParams(); err != nil {
				return err
			}

			var stats *Statistics
			if duration > 0 {
				stats, err = b.RunWithFixedDuration()
			} else {
				stats, err = b.RunWithFixedRequestsNum()
			}
			if err != nil {
				return err
			}
			if out != "" && stats != nil {
				err = stats.Save(out)
				if err != nil {
					fmt.Println()
					return fmt.Errorf("failed to save statistics to file: %s", err)
				}
				fmt.Printf("\nsave statistics to '%s' successfully\n", out)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&targetURL, "url", "u", "", "request URL")
	_ = cmd.MarkFlagRequired("url")
	cmd.Flags().StringVarP(&method, "method", "m", "GET", "request method")
	cmd.Flags().StringSliceVarP(&headers, "header", "e", nil, "request headers")
	cmd.Flags().StringVarP(&body, "body", "b", "", "request body (priority higher than --body-file)")
	cmd.Flags().StringVarP(&bodyFile, "body-file", "f", "", "request body file")

	cmd.Flags().IntVarP(&worker, "worker", "w", runtime.NumCPU()*3, "number of workers concurrently processing requests")
	cmd.Flags().Uint64VarP(&total, "total", "t", 5000, "total requests")
	cmd.Flags().DurationVarP(&duration, "duration", "d", 0, "duration of the test, e.g., 10s, 1m (priority higher than --total)")

	cmd.Flags().StringVarP(&out, "out", "o", "", "save statistics to JSON file")
	cmd.Flags().StringVarP(&pushURL, "push-url", "p", "", "push statistics to target URL once per second ")
	cmd.Flags().StringVarP(&prometheusJobName, "prometheus-job-name", "j", "", "if not empty, the push-url parameter value indicates prometheus url")

	return cmd
}

func newHTTP2Client() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, // Skip certificate validation
			MaxIdleConns:        1000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     90 * time.Second,
			ForceAttemptHTTP2:   true,
		},
		Timeout: 10 * time.Second,
	}
}
