#!/usr/bin/env python

from __future__ import print_function


import sys
from threading import Thread


from docker import Client


from circuits import handler, Event, Component
from circuits.web import Server, Controller, Logger


class docker_event(Event):
    """Docker Event"""


class container_created(docker_event):
    """Container Created Event"""


class container_destroyed(docker_event):
    """Container Destroyed Event"""


class container_started(docker_event):
    """Container Started Event"""


class container_stopped(docker_event):
    """Container Stopped Event"""


class container_killed(docker_event):
    """Container Killed Event"""


class container_died(docker_event):
    """Container Died Event"""


DOCKER_EVENTS = {
    u"create": container_created,
    u"destroy": container_destroyed,
    u"start": container_started,
    u"stop": container_stopped,
    u"kill": container_killed,
    u"die": container_died,
}


class DockerEventManager(Thread):

    def __init__(self, manager, url):
        super(DockerEventManager, self).__init__()

        self.manager = manager
        self.url = url

        self.daemon = True

        self.client = Client(self.url)

    def run(self):
        for event in self.client.events():
            status = event.pop("status")
            docker_event = DOCKER_EVENTS.get(status)
            if docker_event is not None:
                self.manager.fire(docker_event(**event), "docker")
            else:
                print("WARNING: Unknown Docker Event <{0:s}({1:s})>".format(status, repr(event)), file=sys.tderr)

    def stop(self):
        self.client.close()


class DockerEventLogger(Component):

    channel = "docker"

    @handler("*")
    def log_docker_event(self, event, *args, **kwargs):
        if isinstance(event, docker_event):
            print(repr(event))


class App(Component):

    def init(self, url):
        DockerEventManager(self, url).start()
        DockerEventLogger().register(self)


class Root(Controller):

    def index(self):
        return "Hello World!"


def main():
    app = App("unix://var/run/docker.sock")
    Server(("0.0.0.0", 5000)).register(app)
    Logger().register(app)
    Root().register(app)

    #from circuits import Debugger
    #Debugger().register(app)

    app.run()


if __name__ == "__main__":
    main()
