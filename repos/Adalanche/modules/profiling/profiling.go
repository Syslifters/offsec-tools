package profiling

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"runtime"
	runtimepprof "runtime/pprof"
	"sync"
	"time"

	"github.com/felixge/fgprof"
	"github.com/felixge/fgtrace"
	"github.com/lkarlslund/adalanche/modules/ui"
)

type Options struct {
	Datapath string

	EmbeddedProfiler  bool
	CPUProfile        bool
	CPUProfileTimeout int32
	MemProfile        bool
	MemProfileTimeout int32
	FGTrace           bool
	FGProf            bool
}

type profilerControl struct {
	stop             func()
	stopAfterProcess bool
}

type session struct {
	controls []profilerControl
	wg       sync.WaitGroup
}

var (
	sessionMu      sync.Mutex
	currentSession *session
)

func Start(options Options) error {
	sessionMu.Lock()
	defer sessionMu.Unlock()

	s := &session{}

	if options.EmbeddedProfiler {
		go func() {
			port := 6060
			for {
				err := http.ListenAndServe(fmt.Sprintf("localhost:%v", port), nil)
				if err != nil {
					ui.Error().Msgf("Profiling listener failed: %v, trying with new port", err)
					port++
					continue
				}
				ui.Info().Msgf("Profiling listener started on port %v", port)
				return
			}
		}()
	}

	if options.FGProf {
		traceFilename := filepath.Join(options.Datapath, "adalanche-fgprof-"+time.Now().Format("06010215040506")+".json")
		traceFile, err := os.Create(traceFilename)
		if err != nil {
			return fmt.Errorf("error creating fgprof file %v: %w", traceFilename, err)
		}

		traceStopper := fgprof.Start(traceFile, fgprof.FormatPprof)
		s.addControl(options.CPUProfileTimeout, func() {
			defer traceFile.Close()
			if err := traceStopper(); err != nil {
				ui.Error().Msgf("Problem stopping fgprof: %v", err)
			}
		})
	}

	if options.FGTrace {
		traceFilename := filepath.Join(options.Datapath, "adalanche-fgtrace-"+time.Now().Format("06010215040506")+".json")
		trace := fgtrace.Config{Dst: fgtrace.File(traceFilename)}.Trace()
		s.addControl(options.CPUProfileTimeout, func() {
			if err := trace.Stop(); err != nil {
				ui.Error().Msgf("Problem stopping fgtrace: %v", err)
			}
		})
	}

	if options.CPUProfile {
		profileFilename := filepath.Join(options.Datapath, "adalanche-cpuprofile-"+time.Now().Format("06010215040506")+".pprof")
		f, err := os.Create(profileFilename)
		if err != nil {
			return fmt.Errorf("could not set up CPU profiling in file %v: %w", profileFilename, err)
		}
		if err := runtimepprof.StartCPUProfile(f); err != nil {
			_ = f.Close()
			return fmt.Errorf("could not start CPU profiling in file %v: %w", profileFilename, err)
		}

		s.addControl(options.CPUProfileTimeout, func() {
			runtimepprof.StopCPUProfile()
			_ = f.Close()
		})
	}

	if options.MemProfile {
		profileFilename := filepath.Join(options.Datapath, "adalanche-memprofile-"+time.Now().Format("06010215040506")+".pprof")
		f, err := os.Create(profileFilename)
		if err != nil {
			return fmt.Errorf("could not set up memory profiling in file %v: %w", profileFilename, err)
		}

		s.addControl(options.MemProfileTimeout, func() {
			runtime.GC()
			if err := runtimepprof.WriteHeapProfile(f); err != nil {
				ui.Error().Msgf("Problem writing heap profile: %v", err)
			}
			_ = f.Close()
		})
	}

	currentSession = s
	return nil
}

func StopAfterProcessing() {
	sessionMu.Lock()
	s := currentSession
	sessionMu.Unlock()

	if s == nil {
		return
	}

	for _, control := range s.controls {
		if control.stopAfterProcess {
			control.stop()
		}
	}
}

func StopAll() {
	sessionMu.Lock()
	s := currentSession
	currentSession = nil
	sessionMu.Unlock()

	if s == nil {
		return
	}

	for _, control := range s.controls {
		control.stop()
	}
	s.wg.Wait()
}

func (s *session) addControl(timeoutSeconds int32, stopFn func()) {
	var once sync.Once
	stop := func() {
		once.Do(func() {
			stopFn()
			s.wg.Done()
		})
	}

	s.wg.Add(1)
	s.controls = append(s.controls, profilerControl{
		stop:             stop,
		stopAfterProcess: timeoutSeconds == -1,
	})

	if timeoutSeconds > 0 {
		go func() {
			<-time.After(time.Second * time.Duration(timeoutSeconds))
			stop()
		}()
	}
}
