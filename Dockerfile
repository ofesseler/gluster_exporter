FROM debian:jessie
MAINTAINER Oliver Fesseler <oliver@fesseler.info>

EXPOSE 9189
EXPOSE 24007
EXPOSE 24008

# Gluster debian Repo
ADD http://download.gluster.org/pub/gluster/glusterfs/3.8/3.8.5/rsa.pub /tmp
RUN apt-key add /tmp/rsa.pub && rm -f /tmp/rsa.pub

# Add gluster debian repo and update apt
RUN echo deb http://download.gluster.org/pub/gluster/glusterfs/3.8/3.8.5/Debian/jessie/apt jessie main > /etc/apt/sources.list.d/gluster.list
RUN apt-get update

# Install Gluster server
RUN DEBIAN_FRONTEND=noninteractive apt-get -y install glusterfs-server

# Clean
RUN apt-get clean

# Copy gluster_exporter
COPY gluster_exporter /usr/bin/gluster_exporter

# Create gluster volume, start gluster service and gluster_exporter
RUN mkdir -p /mnt/gv_test
COPY gluster-init.sh /usr/bin/gluster-init.sh
RUN chmod a+x /usr/bin/gluster-init.sh

ENTRYPOINT /usr/bin/gluster-init.sh
