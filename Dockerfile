FROM prologic/python-runtime:2.7

EXPOSE 1338/udp 1338/tcp

ENTRYPOINT ["autodock"]
CMD []

RUN apk -U add git && \
    rm -rf /var/cache/apk/*

COPY requirements.txt /tmp/requirements.txt
RUN pip install -r /tmp/requirements.txt && rm /tmp/requirements.txt

WORKDIR /app
COPY . /app/
RUN pip install .
