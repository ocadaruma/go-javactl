# go-javactl

YAML-Configurable Java Application Wrapper

[![Build Status](https://travis-ci.org/ocadaruma/go-javactl.svg?branch=master)](https://travis-ci.org/ocadaruma/go-javactl)

Re-implementation of [mogproject/javactl](https://github.com/mogproject/javactl) in Go.

## Description

Since JVM has many parameters, writing wrapper shell script for launching Java application is painful.

You can easily configure and launch Java applications via `go-javactl` (also original `mogproject/javactl` ofcourse).

## Features

- Load configurations from a YAML file
- Verify OS user name and Java version
- Check if the application has already been running when duplicate running is prohibited
- Execute pre-launch commands
  - if some pre-launch commads failed, main Java application won't executed.
- Log to syslog
- Launch the Java application with the proper options
- Execute post-launch commands

## Installation

`go-javactl` is available in binary format for Linux and Mac OSX.

[Download latest release here](https://github.com/ocadaruma/go-javactl/releases/latest) and put the binary to somewhere you want.

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
