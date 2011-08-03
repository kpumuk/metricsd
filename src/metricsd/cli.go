package main

import (
	"exec"
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"metricsd/config"
)

var (
	configPath       = flag.String("config", config.DEFAULT_CONFIG_PATH, "Set the path to config file")
	listenAddr       = flag.String("listen", config.DEFAULT_LISTEN, "Set the port (+optional address) to listen at")
	dataPath         = flag.String("data", config.DEFAULT_DATA_DIR, "Set the data directory")
	rootPath         = flag.String("root", config.DEFAULT_ROOT_DIR, "Set the root directory")
	debugLevel       = flag.Int("debug", int(config.DEFAULT_SEVERITY), "Set the debug level, the lower - the more verbose (0-5)")
	sliceInt         = flag.Int("slice", config.DEFAULT_SLICE_INTERVAL, "Set the slice interval in seconds")
	writeInt         = flag.Int("write", config.DEFAULT_WRITE_INTERVAL, "Set the write interval in seconds")
	rrdUpdateThreads = flag.Int("threads", config.DEFAULT_RRD_UPDATE_THREADS, "Set the number of RRD update threads")
	batchWrites      = flag.Bool("batch", config.DEFAULT_BATCH_WRITES, "Set the value indicating whether batch RRD updates should be used")
	dnsLookup        = flag.Bool("lookup", config.DEFAULT_LOOKUP_DNS, "Set the value indicating whether reverse DNS lookup should be performed for sources")
	testAndExit      = flag.Bool("test", false, "Validate config file and exit")
)

func parseCommandLineArguments() {
	flag.Parse()

	// Get root directory
	binaryRoot, error := getBinaryRootDir()
	if error != nil {
		fmt.Print(error)
		os.Exit(1)
	}

	// Make config file path absolute
	cfgpath := *configPath
	if !path.IsAbs(cfgpath) {
		cfgpath = path.Join(binaryRoot, cfgpath)
	}

	// Load config from a config file
	config.Load(cfgpath)
	if *testAndExit {
		os.Exit(0)
	}

	// Override options with values passed in command line arguments
	// (but only if they have a value different from a default one)
	if *listenAddr != config.DEFAULT_LISTEN {
		config.Listen = *listenAddr
	}
	if *dataPath != config.DEFAULT_DATA_DIR {
		config.DataDir = *dataPath
	}
	if *rootPath != config.DEFAULT_ROOT_DIR {
		config.RootDir = *rootPath
	}
	if *debugLevel != int(config.DEFAULT_SEVERITY) {
		config.LogLevel = *debugLevel
	}
	if *sliceInt != config.DEFAULT_SLICE_INTERVAL {
		config.SliceInterval = *sliceInt
	}
	if *writeInt != config.DEFAULT_WRITE_INTERVAL {
		config.WriteInterval = *writeInt
	}
	if *rrdUpdateThreads != config.DEFAULT_RRD_UPDATE_THREADS {
		config.RrdUpdateThreads = *rrdUpdateThreads
	}
	if *batchWrites != config.DEFAULT_BATCH_WRITES {
		config.BatchWrites = *batchWrites
	}
	if *dnsLookup != config.DEFAULT_LOOKUP_DNS {
		config.LookupDns = *dnsLookup
	}

	// Make data directory path absolute
	if !path.IsAbs(config.DataDir) {
		config.DataDir = path.Join(binaryRoot, config.DataDir)
	}

	// Make root dir path absolute
	if !path.IsAbs(config.RootDir) {
		config.RootDir = path.Join(binaryRoot, config.RootDir)
	}
}

func getBinaryRootDir() (binaryRoot string, err os.Error) {
	var binaryFile string
	var error os.Error
	if binaryFile, error = exec.LookPath(os.Args[0]); error != nil {
		err = os.NewError(fmt.Sprintf("Failed to retrieve metricsd executable file path: %s\n", error))
		return
	}
	if binaryFile, error = filepath.Abs(binaryFile); error != nil {
		err = os.NewError(fmt.Sprintf("Failed to get absolute path of the metricsd executable file: %s\n", error))
		return
	}
	binaryRoot, _ = path.Split(binaryFile)
	if r, d := path.Split(binaryRoot[:len(binaryRoot)-1]); d == "bin" {
		binaryRoot = r
	}
	return
}
