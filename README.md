[![Build Status](https://travis-ci.org/ofesseler/gluster_exporter.svg?branch=dev)](https://travis-ci.org/ofesseler/gluster_exporter)
[![CircleCI](https://circleci.com/gh/ofesseler/gluster_exporter/tree/dev.svg?style=svg)](https://circleci.com/gh/ofesseler/gluster_exporter/tree/dev)
# gluster_exporter
Gluster exporter for Prometheus


## Installation

```
go get github.com/ofesseler/gluster_exporter
./gluster_exporter
```

## Usage of `gluster_exporter`

| Option                    | Default             | Description
| ------------------------- | ------------------- | -----------------
| -h, --help                | -                   | Displays usage.
| --web.listen-address      | `:9189`             | The address to listen on for HTTP requests.
| --web.metrics-path        | `/metrics`          | URL Endpoint for metrics
| --gluster.volumes         | `_all`              | Comma separated volume names: vol1,vol2,vol3. Default is '_all' to scrape all metrics
| --gluster.executable-path | `/usr/sbin/gluster` | Path to gluster executable.
| --profile                 | `false`             | Enable gluster profiling reports.
| --quota                   | `false`             | Enable gluster quota reports.
| --log.format              | `logger:stderr`     | Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
| --log.level               | `info`              | Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
| --version                 | -                   | Prints version information

## Make
```
build: Go build
docker: build and run in docker container
gometalinter: run some linting checks
gotest: run go tests and reformats

```

**build**: runs go build for gluster_exporter

**docker**: runs docker build and copys new builded gluster_exporter

**gometalinter**: runs [gometalinter](https://github.com/alecthomas/gometalinter) lint tools

**gotest**: runs *vet* and *fmt* go tools

## Relevant Gluster Metrics
Commands within the exporter are executed with `--xml`.

### Command: `gluster volume info`

| Name          | type     | impl. state |
| ------------  | -------- | ------------|
| OpErrno       | Gauge    | implemented |
| opRet         | Gauge    | implemented |
| Status        | Gauge    | implemented |
| BrickCount    | Gauge    | implemented |
| Volumes.Count | Gauge    | implemented |
| Volume.Status | Gauge    | implemented |

### Command: `gluster peer status`

| Name                      | type     | impl. state |
| ------------------------- | -------- | ------------|
| peerStatus.peer.state     | Gauge    | pending     |
| peerStatus.peer.connected | Gauge    | implemented |

### Command: `gluster volume list`
with `gluster volume info` this is obsolete

| Name           | type     | impl. state |
| -------------- | -------- | ------------|
| volList.count  | Gauge    | pending     |
| volList.volume | string   | pending |

### Command: `gluster volume profile gv_test info cumulative`

| Name                                               | type     | impl. state |
| -------------------------------------------------- | -------- | ------------|
| VolProfile.ProfileOp                               | Gauge    | pending     |
| VolProfile.BrickCount                              | Gauge    | implemented     |
| VolProfile.CumulativeStatus.Duration               | Count    | implemented     |
| VolProfile.CumulativeStatus.TotalRead              | Count    | implemented     |
| VolProfile.CumulativeStatus.TotalWrite             | Count    | implemented     |
| VolProfile.CumulativeStats.FopStats.Fop.Name       | CREATE, ENTRYLK, FINODELK, FLUSH, FXATTROP, LOOKUP, OPENDIR, READDIR, STATFS, WRITE | implemented as label |
| VolProfile.CumulativeStats.FopStats.Fop.Hits       | Count    | implemented     |
| VolProfile.CumulativeStats.FopStats.Fop.AvgLatency | Gauge    | implemented     |
| VolProfile.CumulativeStats.FopStats.Fop.MinLatency | Gauge    | implemented     |
| VolProfile.CumulativeStats.FopStats.Fop.MaxLatency | Gauge    | implemented     |


### Command `gluster volume status all detail`
| Name | type | Labels | impl. state |
|------|------|--------|-------------|
| VolStatus.Volumes.Volume[].Node[].SizeFree  | Gauge | hostname, path, volume | implemented |
| VolStatus.Volumes.Volume[].Node[].SizeTotal | Count | hostname, path, volume | implemented |
| VolStatus.Volumes.Volume[].Node[].InodesFree  | Gauge | hostname, path, volume | implemented |
| VolStatus.Volumes.Volume[].Node[].InodesTotal | Count | hostname, path, volume | implemented |


### Metrics in prometheus
| Name          		| Descritpion     |
| ------------  		| -------- |
| up       				| Was the last query of Gluster successful.    |
| volumes_count         | How many volumes were up at the last query.    |
| volume_status        	| Status code of requested volume.    |
| node_size_free_bytes	| Free bytes reported for each node on each instance. Labels are to distinguish origins    |
| node_size_bytes_total	| Total bytes reported for each node on each instance. Labels are to distinguish origins    |
| node_inodes_free	| Free inodes reported for each node on each instance. Labels are to distinguish origins    |
| node_inodes_total	| Total inodes reported for each node on each instance. Labels are to distinguish origins    |
| brick_available 		| Number of bricks available at last query.    |
| brick_duration 		| Time running volume brick.    |
| brick_data_read 		| Total amount of data read by brick.    |
| brick_data_written 	| Total amount of data written by brick.    |
| brick_fop_hits_total		| Total amount of file operation hits.    |
| brick_fop_latency_avg | Average fileoperations latency over total uptime    |
| brick_fop_latency_min | Minimum fileoperations latency over total uptime    |
| brick_fop_latency_max | Maximum fileoperations latency over total uptime    |
| peers_connected 		| Is peer connected to gluster cluster.    |
| heal_info_files_count | File count of files out of sync, when calling 'gluster v heal VOLNAME info    |
| volume_writeable 		| Writes and deletes file in Volume and checks if it is writeable    |
| mount_successful 		| Checks if mountpoint exists, returns a bool value 0 or 1    |

## Troubleshooting
If the following message appears while trying to get some information out of your gluster. Increase scrape interval in `prometheus.yml` to at least 30s.

```
Another transaction is in progress for gv_cluster. Please try again after sometime
```

## Contributors
- coder-hugo
- mjtrangoni

## Similar Projects
glusterfs exporter for prometheus written in rust.
- https://github.com/ibotty/glusterfs-exporter
