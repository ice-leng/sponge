package websocket

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fatih/color"
)

type ErrSet struct {
	m sync.Map
}

func NewErrSet() *ErrSet {
	return &ErrSet{
		m: sync.Map{},
	}
}

func (s *ErrSet) Add(key string) {
	s.m.Store(key, struct{}{})
}

func (s *ErrSet) List() []string {
	var result []string

	s.m.Range(func(key, value interface{}) bool {
		if str, ok := key.(string); ok {
			result = append(result, str)
		}
		return true
	})

	return result
}

type statsCollector struct {
	// Connection
	activeConns         int64
	connectSuccessCount uint64
	connectFailureCount uint64
	totalConnectTime    uint64 // ns
	minConnectTime      uint64 // ns
	maxConnectTime      uint64 // ns

	// Message
	messageSentCount uint64
	messageRecvCount uint64
	sentBytesCount   uint64 // sent total bytes
	recvBytesCount   uint64 // received total bytes
	errorCount       uint64

	errSet *ErrSet // set of errors
}

type Statistics struct {
	URL      string  `json:"url"`      // performed request URL
	Duration float64 `json:"duration"` // seconds

	// Connections
	TotalConnections   uint64  `json:"total_connections"`
	SuccessConnections uint64  `json:"success_connections"`
	FailedConnections  uint64  `json:"failed_connections"`
	AvgConnectLatency  float64 `json:"avg_connect_latency"` // ms
	MinConnectLatency  float64 `json:"min_connect_latency"` // ms
	MaxConnectLatency  float64 `json:"max_connect_latency"` // ms

	// Messages
	TotalMessagesSent     uint64  `json:"total_messages_sent"`
	TotalMessagesReceived uint64  `json:"total_messages_received"`
	SentMessageQPS        float64 `json:"sent_message_qps"`     // sent messages per second
	ReceivedMessageQPS    float64 `json:"received_message_qps"` // received messages per second
	TotalBytesSent        uint64  `json:"total_bytes_sent"`     // bytes
	TotalBytesReceived    uint64  `json:"total_bytes_received"` // bytes

	ErrorCount uint64   `json:"error_count"` // total errors
	Errors     []string `json:"errors"`      // list of errors
}

func (s *Statistics) Save(filePath string) error {
	err := ensureFileExists(filePath)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// --- Atomic operations for thread safety ---

func (s *statsCollector) AddConnectSuccess() {
	atomic.AddUint64(&s.connectSuccessCount, 1)
	atomic.AddInt64(&s.activeConns, 1)
}

func (s *statsCollector) AddConnectFailure() {
	atomic.AddUint64(&s.connectFailureCount, 1)
}

func (s *statsCollector) AddDisconnect() {
	atomic.AddInt64(&s.activeConns, -1)
}

func (s *statsCollector) RecordConnectTime(d time.Duration) {
	ns := uint64(d.Nanoseconds())
	atomic.AddUint64(&s.totalConnectTime, ns)

	minTime := atomic.LoadUint64(&s.minConnectTime)
	if minTime == 0 {
		atomic.CompareAndSwapUint64(&s.minConnectTime, minTime, ns)
	} else if ns < minTime {
		atomic.CompareAndSwapUint64(&s.maxConnectTime, minTime, ns)
	}
	maxTime := atomic.LoadUint64(&s.maxConnectTime)
	if ns > maxTime {
		atomic.CompareAndSwapUint64(&s.maxConnectTime, maxTime, ns)
	}
}

func (s *statsCollector) AddMessageSent() {
	atomic.AddUint64(&s.messageSentCount, 1)
}

func (s *statsCollector) AddMessageRecv() {
	atomic.AddUint64(&s.messageRecvCount, 1)
}

func (s *statsCollector) AddSentBytes(bytes uint64) {
	atomic.AddUint64(&s.sentBytesCount, bytes)
}

func (s *statsCollector) AddRecvBytes(bytes uint64) {
	atomic.AddUint64(&s.recvBytesCount, bytes)
}

func (s *statsCollector) AddError() {
	atomic.AddUint64(&s.errorCount, 1)
}

// Snapshot creates a read-only copy of the current stats.
func (s *statsCollector) Snapshot() statsCollector {
	return statsCollector{
		connectSuccessCount: atomic.LoadUint64(&s.connectSuccessCount),
		connectFailureCount: atomic.LoadUint64(&s.connectFailureCount),
		activeConns:         atomic.LoadInt64(&s.activeConns),
		totalConnectTime:    atomic.LoadUint64(&s.totalConnectTime),
		minConnectTime:      atomic.LoadUint64(&s.minConnectTime),
		maxConnectTime:      atomic.LoadUint64(&s.maxConnectTime),
		messageSentCount:    atomic.LoadUint64(&s.messageSentCount),
		messageRecvCount:    atomic.LoadUint64(&s.messageRecvCount),
		sentBytesCount:      atomic.LoadUint64(&s.sentBytesCount),
		recvBytesCount:      atomic.LoadUint64(&s.recvBytesCount),
		errorCount:          atomic.LoadUint64(&s.errorCount),
	}
}

// PrintReport prints a formatted report of the current stats.
func (s *statsCollector) PrintReport(duration time.Duration, targetURL string) *Statistics {
	snapshot := s.Snapshot()

	totalConnections := snapshot.connectSuccessCount + snapshot.connectFailureCount
	avgConnectTimeMs := 0.0
	if snapshot.connectSuccessCount > 0 {
		avgConnectTimeMs = float64(snapshot.totalConnectTime) / float64(snapshot.connectSuccessCount) / 1e6
	}
	minConnTimeMs := float64(snapshot.minConnectTime) / 1e6
	maxConnTimeMs := float64(snapshot.maxConnectTime) / 1e6
	errSet := s.errSet.List()

	fmt.Printf("\n========== WebSocket Performance Test Report ==========\n\n")
	if snapshot.connectSuccessCount == 0 {
		_, _ = color.New(color.Bold).Println("[Connections]")
		fmt.Printf("  • %-20s%d\n", "Total:", totalConnections)
		fmt.Printf("  • %-20s%d%s\n", "Successful:", 0, color.RedString(" (0%)"))
		fmt.Printf("  • %-20s%d%s\n", "Failed:", snapshot.connectFailureCount, color.RedString(" ✗"))
		fmt.Printf("  • %-20smin: %.2f ms, avg: %.2f ms, max: %.2f ms\n\n", "Latency:", minConnTimeMs, avgConnectTimeMs, maxConnTimeMs)
		if len(errSet) > 0 {
			_, _ = color.New(color.Bold).Println("[Error Details]")
			for _, errStr := range errSet {
				fmt.Printf("  • %s\n", color.RedString(errStr))
			}
		}
		return nil
	}

	sentThroughput := 0.0
	recvThroughput := 0.0
	if seconds := duration.Seconds(); seconds > 0 {
		sentThroughput = float64(snapshot.messageSentCount) / seconds
		recvThroughput = float64(snapshot.messageRecvCount) / seconds
	}

	st := &Statistics{
		URL:      targetURL,
		Duration: duration.Seconds(),

		TotalConnections:   totalConnections,
		SuccessConnections: snapshot.connectSuccessCount,
		FailedConnections:  snapshot.connectFailureCount,
		AvgConnectLatency:  avgConnectTimeMs,
		MinConnectLatency:  minConnTimeMs,
		MaxConnectLatency:  maxConnTimeMs,

		SentMessageQPS:        sentThroughput,
		ReceivedMessageQPS:    recvThroughput,
		TotalMessagesSent:     snapshot.messageSentCount,
		TotalMessagesReceived: snapshot.messageRecvCount,
		TotalBytesSent:        snapshot.sentBytesCount,
		TotalBytesReceived:    snapshot.recvBytesCount,

		ErrorCount: snapshot.errorCount,
		Errors:     errSet,
	}

	_, _ = color.New(color.Bold).Println("[Connections]")
	fmt.Printf("  • %-20s%d\n", "Total:", st.TotalConnections)
	successStr := fmt.Sprintf("  • %-20s%d", "Successful:", st.SuccessConnections)
	failureStr := fmt.Sprintf("  • %-20s%d", "Failed:", st.FailedConnections)
	if st.TotalConnections > 0 {
		if st.TotalConnections == st.SuccessConnections {
			successStr += color.GreenString(" (100%)")
		} else if st.FailedConnections > 0 {
			if st.SuccessConnections == 0 {
				successStr += color.RedString(" (0%)")
			} else {
				successStr += color.YellowString(" (%d%%)", int(float64(st.SuccessConnections)/float64(st.TotalConnections)*100))
			}
			failureStr += color.RedString(" ✗")
		}
	}
	fmt.Println(successStr)
	fmt.Println(failureStr)
	fmt.Printf("  • %-20smin: %.2f ms, avg: %.2f ms, max: %.2f ms\n\n", "Latency:", st.MinConnectLatency, st.AvgConnectLatency, st.MaxConnectLatency)

	_, _ = color.New(color.Bold).Println("[Messages Sent]")
	fmt.Printf("  • %-20s%d\n", "Total Messages:", st.TotalMessagesSent)
	fmt.Printf("  • %-20s%d\n", "Total Bytes:", st.TotalBytesSent)
	fmt.Printf("  • %-20s%.2f msgs/sec\n\n", "Throughput (QPS):", st.SentMessageQPS)

	_, _ = color.New(color.Bold).Println("[Messages Received]")
	fmt.Printf("  • %-20s%d\n", "Total Messages:", st.TotalMessagesReceived)
	fmt.Printf("  • %-20s%d\n", "Total Bytes:", st.TotalBytesReceived)
	fmt.Printf("  • %-20s%.2f msgs/sec\n\n", "Throughput (QPS):", st.ReceivedMessageQPS)

	if len(errSet) > 0 {
		_, _ = color.New(color.Bold).Println("[Error Details]")
		for _, errStr := range errSet {
			fmt.Printf("  • %s\n", color.RedString(errStr))
		}
	}

	return st
}

func ensureFileExists(filePath string) error {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer file.Close()
	}
	return nil
}
