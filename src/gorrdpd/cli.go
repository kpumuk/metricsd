package main

import (
	"exec"
	"flag"
	"os"
	"path"
	"gorrdpd/config"
)

var (
	configPath = flag.String("config", config.DEFAULT_CONFIG_PATH, "Set the path to config file")
	listenAddr = flag.String("listen", config.DEFAULT_LISTEN, "Set the port (+optional address) to listen at")
	dataPath = flag.String("data", config.DEFAULT_DATA_DIR, "Set the data directory")
	rootPath = flag.String("root", config.DEFAULT_ROOT_DIR, "Set the root directory")
	debugLevel = flag.Int("debug", int(config.DEFAULT_SEVERITY), "Set the debug level, the lower - the more verbose (0-5)")
	sliceInt = flag.Int("slice", config.DEFAULT_SLICE_INTERVAL, "Set the slice interval in seconds")
	writeInt = flag.Int("write", config.DEFAULT_WRITE_INTERVAL, "Set the write interval in seconds")
	batchWrites = flag.Bool("batch", config.DEFAULT_BATCH_WRITES, "Set the value indicating whether batch RRD updates should be used")
	dnsLookup = flag.Bool("lookup", config.DEFAULT_LOOKUP_DNS, "Set the value indicating whether reverse DNS lookup should be performed for sources")
	testAndExit = flag.Bool("test", false, "Validate config file and exit")
)

func parseCommandLineArguments() {
	flag.Parse()

	// Get root directory
	binaryFile, _ := exec.LookPath(os.Args[0])
	binaryRoot, _ := path.Split(binaryFile)

	// Make config file path absolute
	cfgpath := *configPath
	if !path.IsAbs(cfgpath) {
		cfgpath = path.Join(binaryRoot, cfgpath)
	}
	// Load config from a config file
	config.Global.Load(cfgpath)
	if *testAndExit {
		os.Exit(0)
	}

	// Override options with values passed in command line arguments
	// (but only if they have a value different from a default one)
	if *listenAddr != config.DEFAULT_LISTEN {
		config.Global.Listen = *listenAddr
	}
	if *dataPath != config.DEFAULT_DATA_DIR {
		config.Global.DataDir = *dataPath
	}
	if *rootPath != config.DEFAULT_ROOT_DIR {
		config.Global.RootDir = *rootPath
	}
	if *debugLevel != int(config.DEFAULT_SEVERITY) {
		config.Global.LogLevel = *debugLevel
	}
	if *sliceInt != config.DEFAULT_SLICE_INTERVAL {
		config.Global.SliceInterval = *sliceInt
	}
	if *writeInt != config.DEFAULT_WRITE_INTERVAL {
		config.Global.WriteInterval = *writeInt
	}
	if *batchWrites != config.DEFAULT_BATCH_WRITES {
		config.Global.BatchWrites = *batchWrites
	}
	if *dnsLookup != config.DEFAULT_LOOKUP_DNS {
		config.Global.LookupDns = *dnsLookup
	}

	// Make data directory path absolute
	if !path.IsAbs(config.Global.DataDir) {
		config.Global.DataDir = path.Join(binaryRoot, config.Global.DataDir)
	}

	// Make root dir path absolute
	if !path.IsAbs(config.Global.RootDir) {
		config.Global.RootDir = path.Join(binaryRoot, config.Global.RootDir)
	}
}