#! /bin/bash

service glusterfs-server start

mkdir -p /data
gluster volume create data $(hostname):/data force
gluster volume start data

#tail -f /var/log/glusterfs/etc-glusterfs-glusterd.vol.log

/usr/bin/glusterfs_exporter_go -version

/usr/bin/glusterfs_exporter_go

service glusterfs-server stop

exit 0