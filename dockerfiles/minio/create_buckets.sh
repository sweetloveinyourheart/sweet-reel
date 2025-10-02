#!/bin/bash

set -e
set -u

(minio server /data &)
while [ "$(mc ready local)" != "The cluster is ready" ]; do
    sleep 1;
done

for bucket in $(echo $MINIO_DEFAULT_BUCKETS | tr ',' ' '); do
	mc mb  /data/$bucket
done

while true; do
    sleep 10;
done