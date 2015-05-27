FROM crux/python:onbuild

EXPOSE 1338/udp 1338/tcp

ENTRYPOINT ["autodock"]
CMD []
