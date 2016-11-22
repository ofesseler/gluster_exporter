# gluster_exporter
Gluster exporter for Prometheus

## Installation 

```
go get github.com/ofesseler/gluster_exporter
./main
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

| Name          | type     |
| ------------  | -------- |
| OpErrno       | Gauge    |
| opRet         | Gauge    |
| Status        | Gauge    |
| BrickCount    | Gauge    |
| Volumes.Count | Gauge    |

### Command: `gluster peer status`

| Name         | type     |
| ------------ | -------- |
| peerStatus.peer.state | Gauge
| peerStatus.peer.connected | Gauge

### Command: `gluster volume list`

| Name         | type     |
| ------------ | -------- |
| volList.count | Gauge
| volList.volume | string

### Command: `gluster volume profile gv_test info cumulative`

| Name         | type     |
| ------------ | -------- |
| volProfile.profileOp | Gauge <- kein plan was das soll
| volProfile.brickCount | Gauge
| volProfile.cumulativeStatus.duration | Count
| volProfile.cumulativeStatus.totalRead | Count
| volProfile.cumulativeStatus.totalWrite | Count
| volProfile.cumulativeStats.fopStats.fop.Name | WRITE, STATFS, FLUSH, OPENDIR, CREATE, LOOKUP, READDIR, FINODELK, ENTRYLK, FXATTROP, 
| volProfile.cumulativeStats.fopStats.fop.hits | count
| volProfile.cumulativeStats.fopStats.fop.avgLatency | Gauge
| volProfile.cumulativeStats.fopStats.fop.minLatency | Gauge
| volProfile.cumulativeStats.fopStats.fop.maxLatency | Gauge

