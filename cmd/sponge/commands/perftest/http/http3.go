package http

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/fatih/color"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/spf13/cobra"

	"github.com/go-dev-frame/sponge/cmd/sponge/commands/perftest/common"
)

// PerfTestHTTP3CMD creates a new cobra.Command for HTTP/3 performance test.
func PerfTestHTTP3CMD() *cobra.Command {
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
		Use:   "http3",
		Short: "Run performance test for HTTP/3 API",
		Long:  "Run performance test for HTTP/3 API.",
		Example: color.HiBlackString(`  # Default mode: worker=CPU*3, 5000 requests, GET request
  sponge perftest http3 --url=https://localhost:8443/user/1

  # Fixed number of requests: 50 workers, 500k requests, GET request
  sponge perftest http3 --worker=50 --total=500000 --url=https://localhost:8443/user/1

  # Fixed number of requests: 50 workers, 500k requests, POST request with JSON body
  sponge perftest http3 --worker=50 --total=500000 --url=https://localhost:8443/user --method=POST --body={\"name\":\"Alice\",\"age\":25}

  # Fixed duration: 50 workers, duration 10s, GET request
  sponge perftest http3 --worker=50 --duration=10s --url=https://localhost:8443/user/1

  # Fixed duration: 50 workers, duration 10s, POST request with JSON body
  sponge perftest http3 --worker=50 --duration=10s --url=https://localhost:8443/user --method=POST --body={\"name\":\"Alice\",\"age\":25}

  # Fixed number of requests: 50 workers, 500k requests, GET request, push statistics to custom HTTP server every 1s
  sponge perftest http3 --worker=50 --total=500000 --url=https://localhost:8443/user/1 --push-url=http://localhost:7070/report

  # Fixed duration: 50 workers, duration 10s, get request, push statistics to Prometheus (job=xxx)
  sponge perftest http3 --worker=50 --duration=10s --url=https://localhost:8443/user/1 --push-url=http://localhost:9091/metrics --prometheus-job-name=perftest-http3`),
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
				version: "HTTP/3",
			}

			b := PerfTestHTTP{
				ID:                common.NewID(),
				Client:            newHTTP3Client(),
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

func newHTTP3Client() *http.Client {
	return &http.Client{
		Transport: &http3.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip certificate validation

			// quic.Config provides fine control over the underlying QUIC connections
			QUICConfig: &quic.Config{
				MaxIdleTimeout:                 10 * time.Second,
				KeepAlivePeriod:                5 * time.Second,
				InitialStreamReceiveWindow:     6 * 1024 * 1024,
				InitialConnectionReceiveWindow: 15 * 1024 * 1024,
				MaxStreamReceiveWindow:         6 * 1024 * 1024,
				MaxConnectionReceiveWindow:     15 * 1024 * 1024,
				MaxIncomingStreams:             1000,
				MaxIncomingUniStreams:          1000,
				HandshakeIdleTimeout:           5 * time.Second,
			},
		},
		Timeout: 10 * time.Second,
	}
}
