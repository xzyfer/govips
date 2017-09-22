// Package govips provides a library for transforming images that is built on lipvips. Libvips
// is an extremely fast C-library. Therefore, govips requires that libvips 8+ be installed
// and available on the target environment.
package vips

//go:generate scripts/codegen.sh

// #cgo pkg-config: vips
// #include "vips/vips.h"
import "C"
import (
	"fmt"
	"runtime"
	"sync"
)

// VipsVersion if the primary version of libvips
const VipsVersion = string(C.VIPS_VERSION)

// VipsMajorVersion is the major version of libvips
const VipsMajorVersion = int(C.VIPS_MAJOR_VERSION)

// VipsMinorVersion if the minor vesrion of libvips
const VipsMinorVersion = int(C.VIPS_MINOR_VERSION)

var (
	running  = false
	initLock sync.Mutex
)

// Config allows fine-tuning of libvips library
type Config struct {
	ConcurrencyLevel int
	MaxCacheMem      int
	MaxCacheSize     int
	ReportLeaks      bool
	CacheTracing     bool
}

// TODO(d): Tune these. Concurrency is set to a safe level but assumes
// openslide is not enabled.
const (
	defaultConcurrencyLevel = 1
	defaultMaxCacheMem      = 100 * 1024 * 1024
	defaultMaxCacheSize     = 1000
)

// Startup sets up the libvips support and ensures the versions are correct. Pass in nil for
// default configuration.
func Startup(cfg *Config) {
	initLock.Lock()
	runtime.LockOSThread()
	defer initLock.Unlock()
	defer runtime.UnlockOSThread()

	if running {
		panic("libvips already running")
	}

	if C.VIPS_MAJOR_VERSION < 8 {
		panic("Requires libvips version 8+")
	}

	cName := C.CString("govips")
	defer freeCString(cName)

	err := C.vips_init(cName)
	if err != 0 {
		panic(fmt.Sprintf("Failed to start vips code=%d", err))
	}

	running = true

	C.vips_concurrency_set(defaultConcurrencyLevel)
	C.vips_cache_set_max(defaultMaxCacheSize)
	C.vips_cache_set_max_mem(defaultMaxCacheMem)

	if cfg != nil {
		C.vips_leak_set(toGboolean(cfg.ReportLeaks))

		if cfg.ConcurrencyLevel > 0 {
			C.vips_concurrency_set(C.int(cfg.ConcurrencyLevel))
		}
		if cfg.MaxCacheMem > 0 {
			C.vips_cache_set_max_mem(C.size_t(cfg.MaxCacheMem))
		}
		if cfg.MaxCacheSize > 0 {
			C.vips_cache_set_max(C.int(cfg.MaxCacheSize))
		}
	}

	initTypes()
}

// Shutdown libvips
func Shutdown() {
	initLock.Lock()
	defer initLock.Unlock()

	if running {
		C.vips_shutdown()
	}
	running = false
}

func printVipsObjects() {
	C.vips_object_print_all()
}

func startupIfNeeded() {
	if !running {
		debug("libvips was forcibly started automatically, consider calling Startup/Shutdown yourself")
		Startup(nil)
	}
}

// ShutdownThread clears the cache for for the given thread
func ShutdownThread() {
	C.vips_thread_shutdown()
}

// VipsCacheDropAll drops the vips operation cache, freeing the allocated memory.
func CacheDropAll() {
	C.vips_cache_drop_all()
}

// VipsDebugInfo outputs to stdout libvips collected data. Useful for debugging.
func DebugInfo() {
	C.im__print_all()
}

// VipsMemoryInfo represents the memory stats provided by libvips.
type VipsMemoryInfo struct {
	Memory          int64
	MemoryHighwater int64
	Allocations     int64
}

// VipsMemory gets memory info stats from libvips (cache size, memory allocs...)
func Memory() VipsMemoryInfo {
	return VipsMemoryInfo{
		Memory:          int64(C.vips_tracked_get_mem()),
		MemoryHighwater: int64(C.vips_tracked_get_mem_highwater()),
		Allocations:     int64(C.vips_tracked_get_allocs()),
	}
}

func LeakTest(fn func()) {
	cfg := &Config{
		ReportLeaks: true,
	}
	Startup(cfg)
	fn()
	runtime.GC()
	ShutdownThread()
	runtime.GC()
	Shutdown()
	runtime.GC()
}
