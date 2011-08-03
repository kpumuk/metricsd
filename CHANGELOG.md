## 0.6.0 (August 3, 2011)

Project has been renamed to MetricsD (was Gorrdpd)

Features:

  - Started migration from "$" groups separator to "."
  - Added dark mode for graphs
  - Added parallel RRD updates

Bugfixes:

  - Fixed favicon installation
  - Fixed legend for percentiles graph
  - Fixed failed items display for count graph
  - Startup script waits for server to shutdown before kill -9
  - Image URLs include .png in the end (nicer playing with Campfire)
  - Fixed bug with standard deviation for 95th percentile graphing

## 0.5.6 (July 26, 2011)

Features:

  - Better build system using goinstall
  - Added favicon.ico
  - Nicer count graphs
  - Added percentiles writer to calculate 90th and 95th percentiles, along with mean and standard deviation of values below them

Bugfixes:

  - Removed dependency on vector package in favour of slices
  - Simplified configuration usage over the project code
  - Fixed bug when 0 was considered as successful event
  - Fixed quartiles calculation logic

## 0.5.5 (May 4, 2011)

Bugfixes:

  - Detect data directory relatively to the executable file full path
  - Updates for Go release r57.1

## 0.5.4 (October 27, 2010)

Bugfixes:

  - Fixed root directory detection

## 0.5.3 (October 27, 2010)

Bugfixes:

  - Better root directory detection

## 0.5.2 (October 27, 2010)

Bugfixes:

  - Updates for Go release 2010-10-20

## 0.5.1 (October 26, 2010)

Bugfixes:

  - Fixed wrong gorrd library repository URL in Makefile

## 0.5.0 (October 17, 2010)

Features:

  - Added live filter on the summary page
  - Renamed YesNo writer to Count
  - Added internal gorrdpd statistics (messages count, traffic, memory usage)
  - Added message source and metric name validation

Bugfixes:

  - RRD sequence start time was invalid
  - Removed "slicer" thread as useless

## 0.4.1 (September 14, 2010)

Bugfixes:

  - Fixed "make install" so it fill delete old executable file backups

## 0.4.0 (September 14, 2010)

Features:

  - Implemented Web UI using web.go
  - Added groups support to protocol and Web UI (group$metric:value)
  - Parallel metrics reporting in benchmarking tool + some other configuration parameters

## 0.3.0 (September 6, 2010)

Features:

  - Added batch RRD updates
  - Added DNS lookups for metric source

Bugfixes:

  - Fixed error in "listen" parameter parsing

## 0.2.2 (August 28, 2010)

Bugfixes:

  - Makefile refactoring

## 0.2.1 (August 20, 2010)

Bugfixes:

  - Fixed bug in quartiles calculation

## 0.2.0 (August 19, 2010)

Features:

  - Added JSON configuration file
  - Added message source to protocol: source@metric:value
  - Added ability to group metrics in a single network packet using ";"
  - Switched from rrdtool to gorrd library (librrd wrapper)
  - Added gorrdpd.sh to start/stop gorrdpd daemon

## 0.1.0 (August 4, 2010)

Features:

  - Simple protocol metric:value
  - Calculates quartiles and fail/success events
  - Updates RRD files using rrdtool
  - Includes udp_generator tool to benchmark the server
  - Handles signals to stop or data dump
