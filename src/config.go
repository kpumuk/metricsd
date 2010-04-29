package config

import (
    "net"
)

type Configuration struct {
    Listen        string
    UDPAddress    *net.UDPAddr
    Data          string
    LogLevel      int
    SliceInterval int
    WriteInterval int
}

var GlobalConfig = &Configuration {}
