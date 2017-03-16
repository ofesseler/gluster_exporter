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
Help is displayed with `-h`.

| Option                   | Default             | Description
| ------------------------ | ------------------- | -----------------
| -help                    | -                   | Displays usage.
| -gluster_executable_path | `/usr/sbin/gluster` | Path to gluster executable.
| -listen-address          | `:9189`             | The address to listen on for HTTP requests.
| -log.format              | `logger:stderr`     | Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
| -log.level               | `info`              | Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
| -metrics-path            | `/metrics`          | URL Endpoint for metrics
| -profile                 | `false`             | When profiling reports in gluster are enabled, set ' -profile true' to get more metrics
| -version                 | -                   | Prints version information
| -volumes                 | `_all`              | Comma separated volume names: vol1,vol2,vol3. Default is '_all' to scrape all metrics


## Make


```
build: Go build
docker: build and run in docker container

```

**build**: runs go build for gluster_exporter

**docker**: runs docker build and copys new builded gluster_exporter


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
| volProfile.profileOp                               | Gauge    | pending     |
| volProfile.brickCount                              | Gauge    | pending     |
| volProfile.cumulativeStatus.duration               | Count    | implemented     |
| volProfile.cumulativeStatus.totalRead              | Count    | implemented     |
| volProfile.cumulativeStatus.totalWrite             | Count    | implemented     |
| volProfile.cumulativeStats.fopStats.fop.Name       | WRITE, STATFS, FLUSH, OPENDIR, CREATE, LOOKUP, READDIR, FINODELK, ENTRYLK, FXATTROP | pending | 
| volProfile.cumulativeStats.fopStats.fop.hits       | count    | implemented     |
| volProfile.cumulativeStats.fopStats.fop.avgLatency | Gauge    | implemented     |
| volProfile.cumulativeStats.fopStats.fop.minLatency | Gauge    | implemented     |
| volProfile.cumulativeStats.fopStats.fop.maxLatency | Gauge    | implemented     |


### Command `gluster volume status all detail`
| Name | type | Labels | impl. state |
|------|------|--------|-------------|
| volStatus.volumes.volume[].node[].sizeFree  | Gauge | hostname, path, volume | implemented |
| volStatus.volumes.volume[].node[].sizeTotal | Gauge | hostname, path, volume | implemented |


### Metrics in prometheus
| Name          		| Descritpion     |
| ------------  		| -------- |
| up       				| Was the last query of Gluster successful.    |
| volumes_count         | How many volumes were up at the last query.    |
| volume_status        	| Status code of requested volume.    |
| node_size_free_bytes	| Free bytes reported for each node on each instance. Labels are to distinguish origins    |
| node_size_total_bytes | Total bytes reported for each node on each instance. Labels are to distinguish origins    |
| brick_count 			| Number of bricks at last query.    |
| brick_duration 		| Time running volume brick.    |
| brick_data_read 		| Total amount of data read by brick.    |
| brick_data_written 	| Total amount of data written by brick.    |
| brick_fop_hits 		| Total amount of file operation hits.    |
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

## Similar Projects
glusterfs exporter for prometheus written in rust.
- https://github.com/ibotty/glusterfs-exporter
