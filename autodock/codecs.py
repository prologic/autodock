# Module:   codecs
# Date:     15th November 2014
# Author:   James Mills, prologic at shortcircuit dot net dot au


"""Codecs"""


import json
from functools import partial


from circuits import Event
from circuits.six import bytes_to_str, text_type


class JSONEncoder(json.JSONEncoder):

    def default(self, obj):
        if isinstance(obj, Event):
            return {
                "name": obj.name,
                "args": obj.args,
                "kwargs": obj.kwargs,
                "success": obj.success,
                "failure": obj.failure,
                "channels": obj.channels,
                "notify": obj.notify
            }

        # Let the base class default method raise the TypeError
        return json.JSONEncoder.default(self, obj)


class JSONDecoder(json.JSONDecoder):

    def decode(self, data):
        obj = json.loads(data)

        name = bytes_to_str(obj["name"].encode("utf-8"))

        args = []
        for arg in obj["args"]:
            if isinstance(arg, text_type):
                arg = arg.encode("utf-8")
            args.append(arg)

        kwargs = {}
        for k, v in obj["kwargs"].items():
            if isinstance(v, text_type):
                v = v.encode("utf-8")
            kwargs[str(k)] = v

        e = Event.create(name, *args, **kwargs)

        e.success = bool(obj["success"])
        e.failure = bool(obj["failure"])
        e.notify = bool(obj["notify"])
        e.channels = tuple(obj["channels"])

        return e


dumps = partial(json.dumps, cls=JSONEncoder)
loads = partial(json.loads, cls=JSONDecoder)
