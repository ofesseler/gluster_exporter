[![Build Status](https://travis-ci.org/ofesseler/gluster_exporter.svg?branch=dev)](https://travis-ci.org/ofesseler/gluster_exporter)
[![CircleCI](https://circleci.com/gh/ofesseler/gluster_exporter/tree/dev.svg?style=svg)](https://circleci.com/gh/ofesseler/gluster_exporter/tree/dev)
# gluster_exporter
Gluster exporter for Prometheus

## Installation 

```
go get github.com/ofesseler/gluster_exporter
./gluster_exporter
```

## Similar Projects
glusterfs exporter for prometheus written in rust. 
- https://github.com/ibotty/glusterfs-exporter

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
| volProfile.cumulativeStats.fopStats.fop.hits       | count    | pending     |
| volProfile.cumulativeStats.fopStats.fop.avgLatency | Gauge    | pending     |
| volProfile.cumulativeStats.fopStats.fop.minLatency | Gauge    | pending     |
| volProfile.cumulativeStats.fopStats.fop.maxLatency | Gauge    | pending     |

