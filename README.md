# gorrdpd

A simple statistics collector, written in Google Go.

## Installation

1. Make sure you have the a working Go environment. See the [install instructions](http://golang.org/doc/install.html). gorrdpd will always compile on the `release` tag.
2. Install [gorrd](http://github.com/kpumuk/gorrd) package.
2. `git clone git://github.com/kpumuk/gorrdpd.git`
3. `cd gorrdpd && make install`

By default gorrdpd will be install to the `/usr/local/gorrdpd` directory. To change it, use `DESTINATION` environment variable with `make install`:

    DESTINATION=/opt/gorrdpd make install

Please note: if you want to run `make install` with `sudo`, make sure that root user has all Go environment variables defined. Also you can use `-E` switch to preserve your current environment (needs `setenv` option set in `sudoers`):

    sudo -E make install

## Configuration

Configuration is stored in JSON format, and you can find an example in `gorrdpd.conf.example`. Every config option could be overridden using command-line arguments. Following options available at the moment:

* `Listen` (`-listen`) — set the port (+optional address) to listen at. Default is `"0.0.0.0:6311"`;
* `DataDir` (`-data`) — set the data directory. Default is `"./data"`;
* `LogLevel` (`-debug`) — set the debug level, the lower - the more verbose (0-5). Default is `1`;
* `SliceInterval` (`-slice`) — set the slice interval in seconds. Default is `10`;
* `WriteInterval` (`-write`) — set the write interval in seconds. Default is `60`.

Another command-line options:

* `-test` — validate the configuration file and exit.
* `-config` — path to the configuration file.

## Protocol details

Gorrdpd uses very simple UDP-based protocol for collecting metrics. Here is what it looks like:

1. `metric:value` — in this simplest case value will be collected in several RRD files; for each writer (see below) two files will be created: `writer-IP-metric.rrd` and `writer-all-metric.rrd`, where `writer` is a name of writer, `IP` — an IP address of the source host, `metric` — metric name.
2. `source@metric:value` — the same as previous, but instead of IP address of the source host, `source` will be used. If it's equal to `all`, no per-host RRD file will be created, only summary for all ones.
3. `metric:value;source@metric:value` — it's possible to send several metrics update in a single packet. Please note: you have to specify `source` for every metric (metrics without source will be saved to IP-based RRD files).

Examples:

    response_time:153
        quartiles-all-response_time.rrd, quartiles-10.0.0.1-response_time.rrd,
        yesno-all-response_time.rrd, yesno-10.0.0.1-response_time.rrd

    app01@response_time:153
        quartiles-all-response_time.rrd, quartiles-app01-response_time.rrd,
        yesno-all-response_time.rrd, yesno-app01-response_time.rrd

    all@response_time:153
        quartiles-all-response_time.rrd, yesno-all-response_time.rrd

    app01@response_time:153;all@requests:-1
        quartiles-all-response_time.rrd, quartiles-app01-response_time.rrd,
        yesno-all-response_time.rrd, yesno-app01-response_time.rrd,
        quartiles-all-requests.rrd, yesno-all-requests.rrd

## Writers

Writer is an implementation of a metrics aggregation algorithm. Each writer generates an RRD file with different (most probably) datasources and RRAs to store aggregated metrics.

There are two writers currently implemented:

1. `quartiles` — calculates [quartiles](http://en.wikipedia.org/wiki/Quartile) for input data. Creates following data sources: `q1` (first quartile), `q2` (second quartile), `q3` (third quartile), `hi` (max sample), `lo` (min sample), `total` (number of samples).
2. `yesno` — calculates number of successful (value >= `0`) and failes (value < `0`) events. Data sources: `ok` — number of successful events, `fail` — number of failed events.
