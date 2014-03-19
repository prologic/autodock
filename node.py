# Module:   node
# Date:     20th March 2014
# Author:   James Mills, prologic at shortcircuit dot net dot au

"""Peer to Peer Node Communcations

This module aims to build enough essential functionality for
an application to employ distributed communications.

Protocol:
    \x00 Hello
    \x06 ACknowledged

TODO:
    - Support UDP and TCP transports
    - Support Websockets
    - Support Web API
      - PUT /event
      - GET /event
    - Support pickle, json and msgpack
    - Support 0mq, mqtt
"""


from __future__ import print_function


import sys
from pickle import dumps, loads


from circuits.net.sockets import UDPServer
from circuits import handler, BaseComponent
from circuits.net.events import broadcast, write


class Node(BaseComponent):

    channel = "node"

    def __init__(self, host="0.0.0.0", port=1338, channel=channel):
        super(Node, self).__init__(channel=channel)

        self.host = host
        self.port = port

        # Peers we keep track of
        # {
        #     (host, port): {
        #         "ls": 1234
        #     }
        # }
        #
        # ls -- Last Seen (no. of ticks since we've seen this peer)
        self._peers = {}

        UDPServer((self.host, self.port), channel=self.channel).register(self)

    def broadcast(self, event, port):
        data = dumps(event)
        self.fire(broadcast(data, port))

        for peer in self._peers.keys():
            self.send(event, *peer)

    def send(self, event, host, port):
        data = dumps(event)
        self.fire(write((host, port), data))

    @handler("read")
    def process_message(self, peer, data):
        if len(data) == 1 and data == b"\x00":
            # Hello Packet
            self._peers[peer] = {"ls": 0}
            self.fire(write(peer, b"\x06"))
        else:
            # Event Packet
            try:
                event = loads(data)
                print("Received Event: {0:s}".format(repr(event)), file=sys.stderr)
                self.fire(event, *event.channels)
            except:
                print("Ignoring Bad Packet: {0:s}".format(repr(data)))
