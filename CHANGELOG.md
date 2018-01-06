# v0.2.7 / 2017-05-09

* [Fix] Bug #11 Incorrect volume name or error with PR #13 Fix unmarshalling of volume status
* [Feature] #12 Feature/add quota metrics
* [Feature] #9 ￼ Add metrics used in promethueus

# v0.2.6 / 2017-02-16

* [Feature] PR#8 Adds new Metrics:
  - mount_successful: Checks if mountpoints of volumes occur in 'mount' at system
  - volume_writeable: Issues Create and Remove file on mounted volume
  - heal_info_files_count: adds all files out of sync together. If this is > 0 your gluster is probably out of sync

# v0.2.5 / 2017-12-14

* [Feature] Added new metrics exposed by command gluster volume status all details:
  - `gluster_node_size_free_bytes`
  - `gluster_node_size_total_bytes`

# v0.2.4 / 2016-12-12

* [Feature] Adds FOP Metrics.

# v0.2.3 / 2016-12-05

* [Feature] Now builds with promu.
* [Feature] `gluster_exporter -version` shows the correct information.

# v0.2.2 / 2016-11-29

* [Feature] Added metrics: `brick duration`, `totalRead`, `totalWrite`
* [Feature] Extended README doc
* [Feature] Added status badges for cirlce and travis
* [Fix] A lot of formatting, fixing typos and refactoring.
* [Fix] Fixed warning of golint.

# v0.1.2 / 2016-11-28

* [Feature] Exposes new value: `peers_connected`.
* [Fix] Removes label "node".
* [Fix] Refactoring.

# v0.1.1 / 2016-11-22

* [Feature] ConstMetrics
* [Feature] Build with promu
* [Feature] vendor.json
* [Fix] Better error messages

# v0.1.0 / 2016-11-21

* First version of `gluster_exporter`.
  ATM it only runs "gluster volume info --xml" as a command and parses the xml output.
