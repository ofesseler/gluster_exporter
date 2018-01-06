#! /bin/bash

service glusterfs-server start

gluster volume create data $(hostname):/data force
gluster volume start data
gluster volume profile data start

glusterfs --volfile-server=localhost --volfile-id=data /mnt/data

dd if=/dev/zero of=/mnt/data/test.zero bs=1M count=10
dd if=/dev/urandom of=/mnt/data/test.random bs=1M count=10

/usr/bin/gluster_exporter --version

/usr/bin/gluster_exporter --profile

service glusterfs-server stop

exit 0
