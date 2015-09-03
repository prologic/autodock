#!/usr/bin/env python

from setuptools import setup, find_packages


from autodock.version import version


def parse_requirements(filename):
    with open(filename, "r") as f:
        for line in f:
            if line and line[:2] not in ("#", "-e"):
                yield line.strip()


setup(
    name="autodock",
    version=version,
    description="Daemon for Docker Automation",
    long_description=open("README.rst", "r").read(),
    author="James Mills",
    author_email="James Mills, prologic at shortcircuit dot net dot au",
    url="https://github.com/prologic/autodock",
    download_url="https://github.com/prologic/autodock/archive/master.zip",
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Environment :: Console",
        "Intended Audience :: Developers",
        "Intended Audience :: System Administrators",
        "License :: OSI Approved :: MIT License",
        "Natural Language :: English",
        "Operating System :: POSIX :: BSD",
        "Operating System :: POSIX :: Linux",
        "Operating System :: MacOS :: MacOS X",
        "Programming Language :: Python :: 2.6",
        "Programming Language :: Python :: 2.7",
        "Programming Language :: Python :: Implementation :: CPython",
        "Programming Language :: Python :: Implementation :: PyPy",
    ],
    license="MIT",
    keywords="docker automation daemon",
    platforms="POSIX",
    packages=find_packages("."),
    install_requires=list(parse_requirements("requirements.txt")),
    entry_points={
        "console_scripts": [
            "autodock=autodock.main:main"
        ]
    },
    zip_safe=True
)
