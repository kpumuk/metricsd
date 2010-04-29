package config

import (
    "net"
    "./logger"
)

const (
    DEFAULT_LISTEN         = "0.0.0.0:6311"
    DEFAULT_RRD_TOOL_PATH  = "/usr/bin/rrdtool"
    DEFAULT_DATA_DIR       = "./data"
    DEFAULT_SEVERITY       = logger.INFO
    DEFAULT_SLICE_INTERVAL = 10
    DEFAULT_WRITE_INTERVAL = 60
)

type Configuration struct {
    Listen        string
    UDPAddress    *net.UDPAddr
    DataDir       string
    RrdToolPath   string
    LogLevel      int
    Logger        logger.Logger
    SliceInterval int
    WriteInterval int
}

var GlobalConfig = &Configuration { }
