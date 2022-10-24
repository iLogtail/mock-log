FROM centos:centos7.9.2009

WORKDIR /bin
COPY bin/mock_log /bin/mock_log
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && echo 'Asia/Shanghai' >/etc/timezone
