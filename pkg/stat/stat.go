// Package stat provides for counting system and process cpu and memory information, alarm notification support.
package stat

import (
	"math"
	"runtime"
	"time"

	"go.uber.org/zap"

	"github.com/zhufuyi/sponge/pkg/stat/cpu"
	"github.com/zhufuyi/sponge/pkg/stat/mem"
)

var (
	printInfoInterval = time.Minute // minimum 1 second
	zapLog, _         = zap.NewProduction()

	notifyCh = make(chan struct{})
)

// Option set the stat options field.
type Option func(*options)

type options struct {
	enableAlarm bool
	zapFields   []zap.Field
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithPrintInterval set print interval
func WithPrintInterval(d time.Duration) Option {
	return func(o *options) {
		if d < time.Second {
			return
		}
		printInfoInterval = d
	}
}

// WithLog set zapLog
func WithLog(l *zap.Logger) Option {
	return func(o *options) {
		if l == nil {
			return
		}
		zapLog = l
	}
}

// WithPrintField set print field
func WithPrintField(fields ...zap.Field) Option {
	return func(o *options) {
		o.zapFields = fields
	}
}

// WithAlarm enable alarm and notify, except windows
func WithAlarm(opts ...AlarmOption) Option {
	return func(o *options) {
		if runtime.GOOS == "windows" {
			return
		}
		ao := &alarmOptions{}
		ao.apply(opts...)
		o.enableAlarm = true
	}
}

// Init initialize statistical information
func Init(opts ...Option) {
	o := &options{}
	o.apply(opts...)

	//nolint
	go func() {
		printTick := time.NewTicker(printInfoInterval)
		defer printTick.Stop()
		sg := newStatGroup()

		for {
			select {
			case <-printTick.C:
				data := printUsageInfo(o.zapFields...)
				if o.enableAlarm {
					if sg.check(data) {
						sendSystemSignForLinux()
					}
				}
			}
		}
	}()
}

// nolint
func sendSystemSignForLinux() {
	select {
	case notifyCh <- struct{}{}:
	default:
	}
}

func printUsageInfo(fields ...zap.Field) *statData {
	defer func() { _ = recover() }()

	mSys := mem.GetSystemMemory()
	mProc := mem.GetProcessMemory()
	cSys := cpu.GetSystemCPU()
	cProc := cpu.GetProcess()

	var cors int32
	for _, ci := range cSys.CPUInfo {
		cors += ci.Cores
	}

	sys := system{
		CPUUsage: cSys.UsagePercent,
		CPUCores: cors,
		MemTotal: mSys.Total,
		MemFree:  mSys.Free,
		MemUsage: float64(int(math.Round(mSys.UsagePercent))), // rounding
	}
	proc := process{
		CPUUsage:   cProc.UsagePercent,
		RSS:        cProc.RSS,
		VMS:        cProc.VMS,
		Alloc:      mProc.Alloc,
		TotalAlloc: mProc.TotalAlloc,
		Sys:        mProc.Sys,
		NumGc:      mProc.NumGc,
		Goroutines: runtime.NumGoroutine(),
	}

	fields = append(fields, zap.Any("system", sys), zap.Any("process", proc))
	zapLog.Info("statistics", fields...)

	return &statData{
		sys:  sys,
		proc: proc,
	}
}

type system struct {
	CPUUsage float64 `json:"cpu_usage"` // system cpu usage, unit(%)
	CPUCores int32   `json:"cpu_cores"` // cpu cores, multiple cpu accumulation
	MemTotal uint64  `json:"mem_total"` // system total physical memory, unit(M)
	MemFree  uint64  `json:"mem_free"`  // system free physical memory, unit(M)
	MemUsage float64 `json:"mem_usage"` // system memory usage, unit(%)
}

type process struct {
	CPUUsage   float64 `json:"cpu_usage"`   // process cpu usage, unit(%)
	RSS        uint64  `json:"rss"`         // use of physical memory, unit(M)
	VMS        uint64  `json:"vms"`         // use of virtual memory, unit(M)
	Alloc      uint64  `json:"alloc"`       // allocated memory capacity, unit(M)
	TotalAlloc uint64  `json:"total_alloc"` // cumulative allocated memory capacity, unit(M)
	Sys        uint64  `json:"sys"`         // requesting memory capacity from the system, unit(M)
	NumGc      uint32  `json:"num_gc"`      // number of GC cycles
	Goroutines int     `json:"goroutines"`  // number of goroutines
}

type statData struct {
	sys  system
	proc process
}
