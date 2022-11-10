#!/bin/bash
ps -ef |  grep -v grep | grep ilogtail | awk '{print $2}' | xargs -I {} kill -s TERM {} &
ps -ef |  grep -v grep | grep vector | awk '{print $2}' | xargs -I {} kill -s TERM {}  &
ps -ef |  grep -v grep | grep filebeat | awk '{print $2}' | xargs -I {} kill -s TERM {}  &
ps -ef |  grep -v grep | grep fluent-bit| awk '{print $2}' | xargs -I {} kill -s TERM {}  &
ps -ef |  grep -v grep | grep rsyslogd| awk '{print $2}' | xargs -I {} kill -s TERM {}  &
wait
