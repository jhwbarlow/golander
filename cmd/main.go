package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/jhwbarlow/golander/pkg/cgroup"
)

const megabyte = 1024 * 1024

var (
	delay     = flag.Duration("d", 0, "the initial delay before starting memory growth")
	sleepTime = flag.Duration("s", 10*time.Second, "the time to sleep between memory growth cycles")
	increment = flag.Int("i", 1, "the amount of memory to grow each cycle (MB)")
	maxMem    = flag.Int("m", 100, "the maximum value that the memory will be grown to (MB)")
)

func main() {
	flag.Parse()

	sinkCap := (*maxMem / *increment) + 1
	memSink := make([][]byte, 0, sinkCap)

	cgroupsEnabled, err := cgroup.CGroupsEnabled()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error determining if cgroups are enabled: %v", err)
	}

	time.Sleep(*delay)
	
	currentlyAlloced := 0
	for i := 0; ; i++ {
		runtime.GC()	
		fmt.Printf("==[Round %d (current allocation: %dMB)]========================================\n", i, currentlyAlloced)

		fmt.Printf("--[/proc/self/status]-----------------------------------------------------------\n")
		status, err := ioutil.ReadFile("/proc/self/status")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading status: %v", err)
			os.Exit(1)
		}
		fmt.Printf("%s", string(status))

		fmt.Printf("--[runtime memstats]------------------------------------------------------------\n")
		memStats := new(runtime.MemStats)
		runtime.ReadMemStats(memStats)
		fmt.Printf("runtime heap alloc: %dMB\n", memStats.HeapAlloc/megabyte)
		fmt.Printf("runtime heap in use: %dMB\n", memStats.HeapInuse/megabyte)
		fmt.Printf("runtime heap idle: %dMB\n", memStats.HeapIdle/megabyte)
		fmt.Printf("runtime heap total: %dMB\n", memStats.TotalAlloc/megabyte)
		fmt.Printf("runtime heap sys: %dMB\n", memStats.HeapSys/megabyte)
		fmt.Printf("runtime sys: %dMB\n", memStats.Sys/megabyte)
		fmt.Printf("runtime heap released: %dMB\n", memStats.HeapReleased/megabyte)

		fmt.Printf("--[/proc/meminfo]---------------------------------------------------------------\n")
		meminfo, err := ioutil.ReadFile("/proc/meminfo")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading meminfo: %v", err)
			os.Exit(1)
		}
		fmt.Printf("%s", string(meminfo))

		if cgroupsEnabled {
			fmt.Printf("--[%s]-------------------------------------------\n", cgroup.CGroupMemStatsPath)
			memstat, err := cgroup.ReadCGroupStats()
			if err != nil {
				fmt.Fprintf(os.Stderr, "error reading cgroup memory stat: %v", err)
				os.Exit(1)
			}
			fmt.Printf("%s", memstat)
		}

		fmt.Println()
		time.Sleep(*sleepTime)

		if !(currentlyAlloced >= *maxMem) {
			memAllocSize := megabyte * *increment
			mem := make([]byte, memAllocSize)
			memSink = append(memSink, mem)
			currentlyAlloced = len(memSink) * *increment
		}
	}
}
