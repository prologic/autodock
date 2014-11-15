autodock
========

autodock is a Daemon for Docker Automation.

Installation
------------

Either pull the prebuilt [Docker](http://docker.com/) image:

    $ docker pull prologic/autodock

Or install from the development repository:

    $ hg clone https://bitbucket.org/prologic/autodock
    $ cd autodock
    $ pip install -r requirements.txt

Example Usage \#1 -- Logging Docker Events
------------------------------------------

Start the daemon:

    $ docker run -d -v /var/run/docker.sock:/var/run/docker.sock --name autodock:autodock prologic/autodock

Link and start an autodock plugin:

    $ docker run -i -t --link autodock prologic/autodock-logger

Now whenever you start a new container autodock will listen for Docker events. The `autodock-logger` plugin will log all Docker Events received by autodock.

Example Usage \#2 -- Automatic Virtual Hosting with hipache
-----------------------------------------------------------

> **note**
>
> This example is still in development.

Start the daemon:

    $ docker run -d --name autodock prologic/autodock

Link and start an autodock plugin:

    $ docker run -d --link autodock prologic/autodock-hipache

Now whenever you start a new container autodock will listen for Docker events and discover containers that have been started. The `autodock-hipache` plugin will specifically listen for starting containers that have a `VIRTUALHOST` environment variable and reconfigure the running `hipache` container.

Start a "Hello World" Web Application:

    $ docker run -d -e VIRTUALHOST=hello.local prologic/hello

Now assuming you had `hello.local` configured in your `/etc/hosts` pointing to your `hipache` container you can now visit <http://hello.local/>

    echo "127.0.0.1 hello.local" >> /etc/hosts
    curl -q -o - http://hello.local/
    Hello World!
