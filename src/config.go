package config

import (
    "fmt"
    "json"
    "net"
    "os"
    "./logger"
)

const (
    DEFAULT_CONFIG_PATH    = "./gorrdpd.conf"
    DEFAULT_LISTEN         = "0.0.0.0:6311"
    DEFAULT_DATA_DIR       = "./data"
    DEFAULT_SEVERITY       = logger.INFO
    DEFAULT_SLICE_INTERVAL = 10
    DEFAULT_WRITE_INTERVAL = 60
    DEFAULT_BATCH_WRITES   = false
    DEFAULT_LOOKUP_DNS     = false
)

type Configuration struct {
    Listen        string        // port and address to listen at
    DataDir       string        // data directory
    LogLevel      int           // debug level, the lower - the more verbose (0-5)
    SliceInterval int           // slice interval in seconds
    WriteInterval int           // write interval in seconds
    BatchWrites   bool          // value indicating whether batch RRD updates should be used
    LookupDns     bool          // value indicating whether reverse DNS lookup should be performed for sources
    UDPAddress    *net.UDPAddr  // address to listen at (for internal usage)
    Logger        logger.Logger // logger instance
}

var Global = &Configuration{
    Listen:        DEFAULT_LISTEN,
    DataDir:       DEFAULT_DATA_DIR,
    LogLevel:      int(DEFAULT_SEVERITY),
    SliceInterval: DEFAULT_SLICE_INTERVAL,
    WriteInterval: DEFAULT_WRITE_INTERVAL,
    BatchWrites:   DEFAULT_BATCH_WRITES,
    LookupDns:     DEFAULT_LOOKUP_DNS,
}

func (config *Configuration) Load(path string) {
    file, error := os.Open(path, os.O_RDONLY, 0)
    if error != nil {
        fmt.Printf("Config file does not exist: %s\n", path)
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    error = decoder.Decode(config)
    if error != nil {
        fmt.Printf("Failed to parse config file: %s\n", error)
        os.Exit(1)
    }
}

func (config *Configuration) String() string {
    return fmt.Sprintf(
        "Configuration:\nListen: \t%s\nData dir:\t%s\nLog level:\t%s\nSlice interval:\t%d\nWrite interval:\t%d\nBatch Writes:\t%t\nLookup DNS:\t%t\n",
        config.Listen,
        config.DataDir,
        logger.Severity(config.LogLevel),
        config.SliceInterval,
        config.WriteInterval,
        config.BatchWrites,
        config.LookupDns,
    )
}
