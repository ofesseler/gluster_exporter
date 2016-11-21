#! /bin/bash

service glusterfs-server start

mkdir -p /data
gluster volume create data $(hostname):/data force
gluster volume start data

/usr/bin/gluster_exporter -version

/usr/bin/gluster_exporter

service glusterfs-server stop

exit 0