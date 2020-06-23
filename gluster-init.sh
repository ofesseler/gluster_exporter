#! /bin/bash

#GlusterFS configuration variables
VOLNAME="data"

# Start gluster manually (systemd is not running)
/usr/sbin/glusterd -p /var/run/glusterd.pid --log-level INFO &
# Wait to start configuring gluster
sleep 10
# Create a volume
gluster volume create "$VOLNAME" "$(hostname)":/"$VOLNAME" force
# Start Gluster volume
gluster volume start "$VOLNAME"
# Enable gluster profile
gluster volume profile "$VOLNAME" start

# Mount the volume
glusterfs --volfile-server=localhost --volfile-id="$VOLNAME" /mnt/"$VOLNAME"

# Write something to the volume
dd if=/dev/zero of=/mnt/"$VOLNAME"/test.zero bs=1M count=10
dd if=/dev/urandom of=/mnt/"$VOLNAME"/test.random bs=1M count=10

# Show the exporter version
/usr/bin/gluster_exporter --version

# Start gluster_exporter
/usr/bin/gluster_exporter --gluster.volumes="data" --profile

# Stop glusterfs
kill -9 "$(cat /var/run/glusterd.pid)"

exit 0
