// Package config implements configuration management for gorrdpd.
package config

import (
	"fmt"
	"json"
	"net"
	"os"
	"gorrdpd/logger"
)

const (
	DEFAULT_CONFIG_PATH		= "./gorrdpd.conf"
	DEFAULT_LISTEN			= "0.0.0.0:6311"
	DEFAULT_DATA_DIR		= "./data"
	DEFAULT_ROOT_DIR		= "."
	DEFAULT_SEVERITY		= logger.INFO
	DEFAULT_SLICE_INTERVAL	= 10
	DEFAULT_WRITE_INTERVAL	= 60
	DEFAULT_BATCH_WRITES	= false
	DEFAULT_LOOKUP_DNS		= false
)

var (
	Listen			string			= DEFAULT_LISTEN			// port and address to listen at
	DataDir			string			= DEFAULT_DATA_DIR			// data directory
	RootDir			string			= DEFAULT_ROOT_DIR			// root directory
	LogLevel		int				= int(DEFAULT_SEVERITY)		// debug level, the lower - the more verbose (0-5)
	SliceInterval	int				= DEFAULT_SLICE_INTERVAL	// slice interval in seconds
	WriteInterval	int				= DEFAULT_WRITE_INTERVAL	// write interval in seconds
	BatchWrites		bool			= DEFAULT_BATCH_WRITES		// value indicating whether batch RRD updates should be used
	LookupDns		bool			= DEFAULT_LOOKUP_DNS		// value indicating whether reverse DNS lookup should be performed for sources
	UDPAddress		*net.UDPAddr								// address to listen at (for internal usage)
	Logger			logger.Logger								// logger instance
)

// Load loads configuration from a JSON file.
func Load(path string) {
	file, error := os.Open(path)
	if error != nil {
		fmt.Printf("Config file does not exist or failed to read the file: %s. Original error: %s\n", path, error)
		return
	}
	defer file.Close()

	config := make(map[string]interface{})
	decoder := json.NewDecoder(file)
	error = decoder.Decode(&config)
	if error != nil {
		fmt.Printf("Failed to parse config file: %s\n", error)
		os.Exit(1)
	}

	if listen, found := config["Listen"]; found {
		Listen = listen.(string)
	}
	if dataDir, found := config["DataDir"]; found {
		DataDir = dataDir.(string)
	}
	if logLevel, found := config["LogLevel"]; found {
		LogLevel = (int)(logLevel.(float64))
	}
	if sliceInterval, found := config["SliceInterval"]; found {
		SliceInterval = (int)(sliceInterval.(float64))
	}
	if writeInterval, found := config["WriteInterval"]; found {
		WriteInterval = (int)(writeInterval.(float64))
	}
	if batchWrites, found := config["BatchWrites"]; found {
		BatchWrites = batchWrites.(bool)
	}
	if lookupDns, found := config["LookupDns"]; found {
		LookupDns = lookupDns.(bool)
	}
}

// String returns a string representation of current configuration.
func String() string {
	return fmt.Sprintf(
		"Configuration:\nListen: \t%s\nData dir:\t%s\nRoot dir:\t%s\nLog level:\t%s\nSlice interval:\t%d\nWrite interval:\t%d\nBatch Writes:\t%t\nLookup DNS:\t%t\n",
		Listen,
		DataDir,
		RootDir,
		logger.Severity(LogLevel),
		SliceInterval,
		WriteInterval,
		BatchWrites,
		LookupDns,
	)
}
