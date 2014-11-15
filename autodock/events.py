# Module:   events
# Date:     15th November 2014
# Author:   James Mills, prologic at shortcircuit dot net dot au


"""Events"""


from circuits import Event


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


class container_exported(docker_event):
    """Container Exported Event"""


class container_paused(docker_event):
    """Container Paused Event"""


class container_restarted(docker_event):
    """Container Restarted Event"""


class container_unpaused(docker_event):
    """Container Unpaused Event"""


class image_untagged(docker_event):
    """Image Untagged Event"""


class image_deleted(docker_event):
    """Image Delete Event"""


DOCKER_EVENTS = {
    u"create": container_created,
    u"destroy": container_destroyed,
    u"start": container_started,
    u"stop": container_stopped,
    u"kill": container_killed,
    u"die": container_died,
    "export": container_exported,
    "pause": container_paused,
    "restart": container_restarted,
    "unpause": container_unpaused,
    "untag": image_untagged,
    "delete": image_deleted,
}
