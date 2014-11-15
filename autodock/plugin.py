# Module:   plugin
# Date:     15th November 2014
# Author:   James Mills, prologic at shortcircuit dot net dot au


"""Plugin

Subclass :class:`Plugin` to create autodock plugins with standarized CLI Options and API.
"""


from __future__ import print_function


from os import environ
from argparse import ArgumentDefaultsHelpFormatter, ArgumentParser


from circuits import Component, Debugger


from autodock.node import hello, Node
from autodock.utils import parse_bind


def parse_args(parse=True):
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
        default=environ.get("URL", environ.get("AUTODOCK_PORT", "udp://127.0.0.1:1338")),
        help="autodock Daemon URL"
    )

    return parser.parse_args() if parse else parser


class Plugin(Component):

    def init(self, parse_args_cb=None):
        # Allow ArgumentsParser to be extended.
        if parse_args_cb is not None:
            self.args = parse_args_cb(parse_args(False)).parse_args()
        else:
            self.args = parse_args()

        self.bind = parse_bind(self.args.bind)
        self.url = parse_bind(self.args.url)

    def started(self, *args):
        if self.args.debug:
            Debugger().register(self)

        self.node = Node(*self.bind).register(self)

    def ready(self, *args):
        self.node.send(hello(), self.url)

    def container_created(self, event, **data):
        """Container created Event

        Override this in a subclass to receiver the created event
        """

    def container_started(self, event, **data):
        """Container started Event

        Override this in a subclass to receiver the started event
        """

    def container_stopped(self, event, **data):
        """Container stopped Event

        Override this in a subclass to receiver the stopped event
        """

    def container_died(self, event, **data):
        """Container died Event

        Override this in a subclass to receiver the died event
        """

    def container_killed(self, event, **data):
        """Container killed Event

        Override this in a subclass to receiver the killed event
        """

    def container_destroyed(self, event, **data):
        """Container destroyed Event

        Override this in a subclass to receiver the destroyed event
        """
