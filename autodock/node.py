# Module:   node
# Date:     20th March 2014
# Author:   James Mills, prologic at shortcircuit dot net dot au


"""Peer to Peer Node Communcations

This module aims to build enough essential functionality for
an application to employ distributed communications.

Default Port::
    >>> s = "circuits.node"
    >>> xs = map(ord, s)
    >>> sum(xs)
    1338

TODO:
    [ ] Support UDP and TCP transports
    [ ] Support Websockets
    [ ] Support Web API
        [ ] PUT /event
        [ ] GET /event
    [ ] Support Serializations:
        [ ] json
        [ ] msgpack
"""


from __future__ import print_function


import sys
from traceback import format_exc


from circuits.net.events import write
from circuits.net.sockets import UDPServer
from circuits import handler, BaseComponent, Event


from .codecs import dumps, loads


class hello(Event):
    """hello Event"""


class Node(BaseComponent):

    channel = "node"

    def __init__(self, host="0.0.0.0", port=1338, channel=channel):
        super(Node, self).__init__(channel=channel)

        self.host = host
        self.port = port

        # Peers we keep track of
        self.peers = set()

        UDPServer((self.host, self.port), channel=self.channel).register(self)

    def broadcast(self, event):
        for peer in self.peers:
            self.send(event, peer)

    def send(self, event, peer):
        data = dumps(event)
        self.fire(write(peer, data))

    @handler("read")
    def _process_message(self, peer, data):
        # Event Packet
        try:
            self.peers.add(peer)
            event = loads(data)
            self.fire(event, *event.channels)
        except Exception as e:
            print(
                "ERROR: Could not parse packet.\n"
                "Error: {} Data: {}".format(e, data),
                file=sys.stderr
            )
            print(format_exc(), file=sys.stderr)
