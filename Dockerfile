FROM centos:centos7.9.2009

WORKDIR /app/mock
COPY bin/mock_log /app/mock/mock_log
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone