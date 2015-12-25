FROM prologic/python-runtime:2.7-onbuild

EXPOSE 1338/udp 1338/tcp

ENTRYPOINT ["autodock"]
CMD []
