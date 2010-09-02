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
)

type Configuration struct {
    // Configurable values
    Listen        string
    DataDir       string
    LogLevel      int
    SliceInterval int
    WriteInterval int
    BatchWrites   bool
    // Values for internal usage
    UDPAddress    *net.UDPAddr
    Logger        logger.Logger
}

var GlobalConfig = &Configuration {
    Listen:        DEFAULT_LISTEN,
    DataDir:       DEFAULT_DATA_DIR,
    LogLevel:      int(DEFAULT_SEVERITY),
    SliceInterval: DEFAULT_SLICE_INTERVAL,
    WriteInterval: DEFAULT_WRITE_INTERVAL,
    BatchWrites:   DEFAULT_BATCH_WRITES,
}

func (config *Configuration) Load(path string) {
    file, error := os.Open(path, os.O_RDONLY, 0)
    if error != nil {
        fmt.Printf("Config file does not exist: %s\n", path)
        return
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    error = decoder.Decode(GlobalConfig)
    if error != nil {
        fmt.Printf("Failed to parse config file: %s\n", error)
        os.Exit(1)
    }
}

func (config *Configuration) String() string {
    return fmt.Sprintf(
        "Configuration:\nListen: \t%s\nData dir:\t%s\nLog level:\t%s\nSlice interval:\t%d\nWrite interval:\t%d\nBatch Writes:\t%t\n",
        GlobalConfig.Listen,
        GlobalConfig.DataDir,
        logger.Severity(GlobalConfig.LogLevel),
        GlobalConfig.SliceInterval,
        GlobalConfig.WriteInterval,
        GlobalConfig.BatchWrites,
    )
}
