#!/bin/bash
nohup vector -c /home/aivin/mock/mock-tools/vector/vector.toml -w > vector_run.log  2>&1  &
nohup ~/mock/mock-tools/ilogtail/ilogtail > ilogtail_run.log 2>&1  &
nohup rsyslogd -f /home/aivin/mock/mock-tools/rsyslog/rsys.conf  -i /home/aivin/mock/mock-tools/rsyslog/rsys.pid  > rsys_run.log 2>&1  &
nohup filebeat -path.data=~/mock/mock-tools/filebeat  --path.config=/home/aivin/mock/mock-tools/filebeat  > filebeat_run.log 2>&1  &
nohup fluent-bit -c fluent-bt/fluent-bit.conf  > filebeat_run.log 2>&1  &
