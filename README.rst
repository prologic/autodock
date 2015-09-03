.. _docker: http://docker.com/
.. _dotCloud: http://dotcloud.com/
.. _hipache: https://github.com/hipache/hipache
.. _autodock: https://github.com/prologic/autodock
.. _autodock-cron: https://github.com/prologic/autodock
.. _autodock-logger: https://github.com/prologic/autodock-logger
.. _autodock-hipache: https://github.com/prologic/autodock-hipache


autodock
========

`autodock`_ is a Daemon for Docker Automation.

autodock is MIT licensed.

Installation
------------

Either pull the automatically updated `Docker`_ image::
    
    $ docker pull prologic/autodock

Or install from the development repository::
    
    $ git clone https://github.com/prologic/autodock.git
    $ cd autodock
    $ pip install -r requirements.txt


Plugins
-------

autodock comes with a number of plugins where each piece of functionality is
provided by a separate plugin. Each plugin is "linked" to autodock to receive
Docker events and issue new Docker API commands. The following list is a list
of the currently available plugins for production use:

- `autodock-cron`_      -- Provides a Cron-like scheduler for Containers
- `autodock-logger`_    -- Logs Dockers Events
- `autodock-hipache`_   -- Automatically registers virtualhosts with `hipache`_


Example #1 -- Logging Docker Events
-----------------------------------

.. note:: See: `autodock Logger plugin <http://github.com/prologic/autodock-logger>`_

Start the daemon::
    
    $ docker run -d -v /var/run/docker.sock:/var/run/docker.sock --name autodock:autodock prologic/autodock

Link and start an autodock plugin::
    
    $ docker run -i -t --link autodock prologic/autodock-logger

Now whenever you start a new container autodock will listen for Docker events.
The ``autodock-logger`` plugin will log all Docker Events received by autodock.


Example #2 -- Automatic Virtual Hosting with hipache
----------------------------------------------------

.. note:: See `autodock Hipache plugin <http://github.com/prologic/autodock-hipache>`_

Start the daemon::
    
    $ docker run -d --name autodock prologic/autodock

Link and start an autodock plugin::
    
    $ docker run -d --link autodock prologic/autodock-hipache

Now whenever you start a new container autodock will listen for Docker events
and discover containers that have been started. The ``autodock-hipache`` plugin
will specifically listen for starting containers that have a ``VIRTUALHOST``
environment variable and reconfigure the running ``hipache`` container.

Start a "Hello World" Web Application::
    
    $ docker run -d -e VIRTUALHOST=hello.local prologic/hello

Now assuming you had ``hello.local`` configured in your ``/etc/hosts``
pointing to your ``hipache`` container you can now visit http://hello.local/

::
    
    echo "127.0.0.1 hello.local" >> /etc/hosts
    curl -q -o - http://hello.local/
    Hello World!

.. note:: This method of hosting and managing webapps and websites is in production deployments and talked about in more detail in the post `A Docker-based mini-PaaS <http://shortcircuit.net.au/~prologic/blog/article/2015/03/24/a-docker-based-mini-paas/>`_.


Example #3 -- Cron-like Scheduling of Containers
------------------------------------------------

.. note:: See `autodock Cron plugin <http://github.com/prologic/autodock-cron>`_

Start the daemon::
    
    $ docker run -d --name autodock prologic/autodock

Link and start an autodock plugin::
    
    $ docker run -d --link autodock prologic/autodock-cron

Now whenever you create a new container autodock will listen for Docker events
and discover containers that have been created. The ``autodock-cron`` plugin
will specifically listen for created containers that have a ``CRON``
environment variable and schedule a job based on the cron expression supplied
and re-run that container when its scheduled has triggered.

Start a "Hello" Busybox Container::
    
    $ docker run -i -t --name hello busybox sh -c "echo Hello"

After about three minutes or so you should see the following in the logs::
    
    $ docker logs hello
    Hello
    Hello
    Hello
