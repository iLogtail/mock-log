#!/bin/bash
./mock_log -log-type=nginx  -path="/home/aivin/mock/mock-log/access.log" -stdout=false -total-count=1000  -logs-per-sec=500
sleep 5
mv /home/aivin/mock/mock-log/access.log /home/aivin/mock/mock-log/access.log.0
echo "New file line" >> /home/aivin/mock/mock-log/access.log
echo "Last line" >> /home/aivin/mock/mock-log/access.log.0