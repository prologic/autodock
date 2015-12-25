# Module:   events
# Date:     15th November 2014
# Author:   James Mills, prologic at shortcircuit dot net dot au


"""Events"""


from circuits import Event


class docker_event(Event):
    """Docker Event"""


class container_attached(docker_event):
    """Container Attached Event"""


class container_committed(docker_event):
    """Container Commited Event"""


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


class container_renamed(docker_event):
    """Container Renamed Event"""


class container_resized(docker_event):
    """Container Resized Event"""


class container_restarted(docker_event):
    """Container Restarted Event"""


class container_unpaused(docker_event):
    """Container Unpaused Event"""


class image_untagged(docker_event):
    """Image Untagged Event"""


class image_deleted(docker_event):
    """Image Delete Event"""


class pull(docker_event):
    """Pull Event"""


DOCKER_EVENTS = {
    u"attach": container_attached,
    u"commit": container_committed,
    u"create": container_created,
    u"destroy": container_destroyed,
    u"start": container_started,
    u"stop": container_stopped,
    u"kill": container_killed,
    u"die": container_died,
    u"export": container_exported,
    u"pause": container_paused,
    u"rename": container_renamed,
    u"resize": container_resized,
    u"restart": container_restarted,
    u"unpause": container_unpaused,
    u"untag": image_untagged,
    u"delete": image_deleted,
    u"pull": pull,
}
