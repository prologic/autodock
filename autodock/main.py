#!/usr/bin/env python
# Module:   main
# Date:     15th November 2014
# Author:   James Mills, prologic at shortcircuit dot net dot au


"""Daemon for Docker Automation"""


from __future__ import print_function


import sys
from time import time
from os import environ
from json import loads
from threading import Thread
from argparse import ArgumentDefaultsHelpFormatter, ArgumentParser


from docker import Client
from circuits.web import Server, JSONRPC
from circuits import handler, Component, Debugger


from .node import Node
from .utils import parse_bind
from .events import docker_event, DOCKER_EVENTS


class DockerRPCServer(Component):

    channel = "rpc"

    def init(self, bind, url, channel=channel):
        self.client = Client(url)

        Server(bind).register(self)
        JSONRPC(rpc_channel=self.channel).register(self)

    def ping(self, ts):
        return time() - ts

    def docker(self, method, *args, **kwargs):
        # TODO: Make this async
        return getattr(self.client, method)(*args, **kwargs)


class DockerEventManager(Thread):

    def __init__(self, manager, url=None):
        super(DockerEventManager, self).__init__()

        self.manager = manager
        self.url = url

        self.daemon = True

        self.client = Client(self.url)

    def run(self):
        for payload in self.client.events():
            event = loads(payload)
            status = event.pop("status")
            docker_event = DOCKER_EVENTS.get(status)
            if docker_event is not None:
                self.manager.fire(docker_event(**event), "docker")
            else:
                print(
                    "WARNING: Unknown Docker Event <{0:s}({1:s})>".format(
                        status, repr(event)
                    ),
                    file=sys.stderr
                )

    def stop(self):
        self.client.close()


class EventBroadcaster(Component):

    def init(self, host="127.0.0.1", port=1338):
        self.host = host
        self.port = port

        self.node = Node(self.host, self.port).register(self)

    @handler("*", channel="docker")
    def broadcast_docker_event(self, event, *args, **kwargs):
        if isinstance(event, docker_event):
            self.node.broadcast(event)


class App(Component):

    def init(self, args):
        if args.debug:
            Debugger().register(self)

        bind = parse_bind(args.bind)

        DockerRPCServer(bind, args.url).register(self)
        DockerEventManager(self, args.url).start()
        EventBroadcaster(*bind).register(self)

    def signal(self, *args):
        raise SystemExit(0)


def parse_args():
    parser = ArgumentParser(
        description=__doc__,
        formatter_class=ArgumentDefaultsHelpFormatter
    )

    parser.add_argument(
        "-b", "--bind", action="store", dest="bind", metavar="INT", type=str,
        default=environ.get("BIND", "0.0.0.0:1338"),
        help="Interface and Port to Bind to"
    )

    parser.add_argument(
        "-d", "--debug", action="store_true", dest="debug",
        default=environ.get("DEBUG", False),
        help="Enable Debug Mode"
    )

    parser.add_argument(
        "-u", "--url", action="store", dest="url", metavar="URL", type=str,
        default=environ.get("URL", None),
        help="Docker Daemon URL"
    )

    return parser.parse_args()


def main():
    App(parse_args()).run()


if __name__ == "__main__":
    main()
