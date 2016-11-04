# go-javactl

[![Build Status](https://travis-ci.org/ocadaruma/go-javactl.svg?branch=master)](https://travis-ci.org/ocadaruma/go-javactl)

Re-implementation of [mogproject/javactl](https://github.com/mogproject/javactl) in Go.

YAML-Configurable Java Application Wrapper

## Features

- Load settings from a YAML file
- Verify OS user name and Java version
- Check if the application has already been running when duplicate running is prohibited
- Execute pre-launch commands
- Log to syslog
- Launch the Java application with the proper options
- Execute post-launch commands

## Installation

T.B.D.

## Configuration Example

See the [example](./testdata/example.yml).

## Usage

- Dry-run mode

```
javactl --check /path/to/your-app.yml
```

- Launch the java application

```
javactl /path/to/your-app.yml
```

- Launch with arguments

```
javactl /path/to/your-app.yml --option-for-your-app arg1 arg2
```
