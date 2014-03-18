#!/usr/bin/env python

from __future__ import print_function


from docker import Client


c = Client("unix://var/run/docker.sock")
print("\n".join([x["Id"][:10] for x in c.containers()]))
