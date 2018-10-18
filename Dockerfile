FROM centos:7
MAINTAINER yarntime@163.com

RUN yum install -y \
    libyaml 

ADD ./task-mgmt /usr/local/bin/task-mgmt

ENTRYPOINT ["/usr/local/bin/task-mgmt"]
