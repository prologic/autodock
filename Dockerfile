# Docker Image for autodock
#
# VERSION: 0.1

FROM prologic/crux-python
MAINTAINER James Mills <prologic@shortcircuitnet.au>

ADD . /

CMD ["/start"]
