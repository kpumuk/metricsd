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
    DEFAULT_RRD_TOOL_PATH  = "/usr/bin/rrdtool"
    DEFAULT_SEVERITY       = logger.INFO
    DEFAULT_SLICE_INTERVAL = 10
    DEFAULT_WRITE_INTERVAL = 60
)

type Configuration struct {
    // Configurable values
    Listen        string
    DataDir       string
    RrdToolPath   string
    LogLevel      int
    SliceInterval int
    WriteInterval int
    // Values for internal usage
    UDPAddress    *net.UDPAddr
    Logger        logger.Logger
}

var GlobalConfig = &Configuration {
    DEFAULT_LISTEN,
    DEFAULT_DATA_DIR,
    DEFAULT_RRD_TOOL_PATH,
    int(DEFAULT_SEVERITY),
    DEFAULT_SLICE_INTERVAL,
    DEFAULT_WRITE_INTERVAL,
    nil,
    nil,
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
        "Configuration:\nListen: \t%s\nData dir:\t%s\nRRDTool path:\t%s\nLog level:\t%s\nSlice interval:\t%d\nWrite interval:\t%d",
        GlobalConfig.Listen,
        GlobalConfig.DataDir,
        GlobalConfig.RrdToolPath,
        logger.Severity(GlobalConfig.LogLevel),
        GlobalConfig.SliceInterval,
        GlobalConfig.WriteInterval,
    )
}
