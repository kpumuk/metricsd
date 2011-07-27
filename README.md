# MetricsD

A simple statistics collector, written in Google Go.

## Installation

First make sure you have the a working Go environment. See the [install instructions](http://golang.org/doc/install.html). MetricsD will always compile on the `release` tag.

Then you will be able to clone and compile MetricsD:

    git clone --recursive --branch release git://github.com/kpumuk/metricsd.git
    cd metricsd
    make install

By default MetricsD will be install to the `/usr/local/metricsd` directory. To change it, use `DESTINATION` environment variable with `make install`:

    DESTINATION=/opt/metricsd make install

Please note: if you want to run `make install` with `sudo`, make sure that root user has all Go environment variables defined. Also you can use `-E` switch to preserve your current environment (needs `setenv` option set in `sudoers`):

    sudo -E make install

## Configuration

Configuration is stored in JSON format, and you can find an example in `metricsd.conf.example`. Every config option could be overridden using command-line arguments. Following options available at the moment:

* `Listen` (`-listen`) — set the port (+optional address) to listen at. Default is `"0.0.0.0:6311"`;
* `DataDir` (`-data`) — set the data directory. Default is `"./data"`;
* `LogLevel` (`-debug`) — set the debug level, the lower - the more verbose (0-5). Default is `1`;
* `SliceInterval` (`-slice`) — set the slice interval in seconds. Default is `10`;
* `WriteInterval` (`-write`) — set the write interval in seconds. Default is `60`;
* `BatchWrites` (`-batch`) — set the value indicating whether batch RRD updates should be used. Default is `false`;
* `LookupDns` (`-lookup`) — set the value indicating whether reverse DNS lookup should be performed for sources.

Another command-line options:

* `-test` — validate the configuration file and exit.
* `-config` — path to the configuration file.

## Protocol details

MetricsD uses very simple UDP-based protocol for collecting metrics. Here is what it looks like:

1. `metric:value` — in this simplest case value will be collected in several RRD files; for each writer (see below) two files will be created: `IP/metric-writer.rrd` and `all/metric-writer.rrd`, where `writer` is a name of writer, `IP` — an IP address of the source host, `metric` — metric name.
2. `source@metric:value` — the same as previous, but instead of IP address of the source host, `source` will be used. If it's equal to `all`, no per-host RRD file will be created, only summary for all ones.
3. `metric:value;source@metric:value` — it's possible to send several metrics update in a single packet. Please note: you have to specify `source` for every metric (metrics without source will be saved to IP-based RRD files).
3. `group$metric:value` — metrics could be grouped in UI based on the `group`
value.

Examples:

    response_time:153
        all/response_time-quartiles.rrd, 10.0.0.1/response_time-quartiles.rrd,
        all/response_time-yesno.rrd, 10.0.0.1/response_time-yesno.rrd

    app01@response_time:153
        all/response_time-quartiles.rrd, app01/response_time-quartiles.rrd,
        all/response_time-yesno.rrd, app01/response_time-yesno.rrd

    all@response_time:153
        all/response_time-quartiles.rrd, all/response_time-yesno.rrd

    app01@response_time:153;all@requests:-1
        all/response_time-quartiles.rrd, app01/response_time-quartiles.rrd,
        all/response_time-yesno.rrd, app01/response_time-yesno.rrd,
        all/requests-quartiles.rrd, all/requests-yesno.rrd

## Writers

Writer is an implementation of a metrics aggregation algorithm. Each writer generates an RRD file with different (most probably) datasources and RRAs to store aggregated metrics.

There are two writers currently implemented:

1. `count` — calculates number of successful (value > `0`) and failes (value < `0`) events. Data sources: `ok` — number of successful events, `fail` — number of failed events.
2. `quartiles` — calculates [quartiles](http://en.wikipedia.org/wiki/Quartile) for input data. Creates following data sources: `q1` (first quartile), `q2` (second quartile), `q3` (third quartile), `hi` (max sample), `lo` (min sample), `total` (number of samples).
3. `percentiles` — calculates 90th and 95th [percentiles](http://en.wikipedia.org/wiki/Percentile) for input data, along with [mean value](http://en.wikipedia.org/wiki/Arithmetic_mean) and [standard deviation](http://en.wikipedia.org/wiki/Standard_deviation) for values under the percentile. Creates following data sources: `pct90` (90th percentile), `pct90mean` (mean of values under 90th percentile), `pct90dev` (standard deviation of values under 95th percentile), `pct95` (95th percentile), `pct95mean` (mean of values under 95th percentile), `pct95dev` (standard deviation of values under 95th percentile).

## Screenshots

![MetricsD: Index Page](http://kpumuk.github.com/metricsd/images/index.png)

![MetricsD: Metric Details](http://kpumuk.github.com/metricsd/images/metric.png)
