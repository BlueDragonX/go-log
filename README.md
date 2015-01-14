A Simple Go Logger
==================
A simplified log library with support for logging to console (stderr), file,
and syslog.  It supports three log levels: debug, info, and error.

[![Build Status](https://travis-ci.org/BlueDragonX/simplelog.svg?branch=master)](https://travis-ci.org/BlueDragonX/simplelog)

Usage
-----
In its simplest form log will output use the `info` log level and output
messages to the console. For example:

	logger := log.New()
	logger.Info("Hello, world!")

It may be configured to log elsewhere by configuring it with a target:

	target := log.Target("file:///var/log/example.log")
	logger := log.New(target)
	logger.Info("Hello, world!")

By default a logger has its level set to `info` which will log `info` and
`error` messages. To change the level:

	level := log.Level("debug")
	logger := log.New(level)
	logger.Debug("Too much information here.")

You can combine the two:

	target := log.Syslog
	level := log.Level("error")
	logger := log.New(target, level)
	logger.Error("Oops! Something went wrong.")

The `log.Syslog` value is a target which points to the local syslog server.

Targets
-------
The `Target` function takes the URI (as a string) of a file or network socket
to log to. If the file is a socket it is assumed to be attached to a running
syslog server. Network sockets assume the same. The special values "syslog",
"stderr", and "stdout" will cause the target to log to syslog, stderr, and
stdout respectively. The `log.Syslog` and `log.Console` values are shortcuts to
calling `log.Target("syslog")` and `log.Target("stderr")`.

Levels
------
The `Level` function takes the log level (as as string) as its only argument.
Case is ignored. Valid levels are "debug", "info", and "error". Log messages at
or greater than that level will be logged. The special values `log.Debug`,
`log.Info`, and `log.Error` are equivelent to `log.Level("debug")`,
`log.Level("info")`, and `log.Level("error")`.

Logging
-------
The following methods log to the logger:

- `Print`: Log a message at the provided level.
- `Printf`: Log a formatted message at the provided level.
- `Panic`: Log a message at the `error` level and call panic().
- `Panicf`: Log a formatted message at the `error` level and call panic().
- `Fatal`: Log a message at the `error` level and call os.Exit(1).
- `Fatalf`: Log a formatted message at the `error` level and call os.Exit(1).
- `Debug`: Log a message at the `debug` level.
- `Debugf`: Log a formatted message at the `debug` level.
- `Info`: Log a message at the `info` level.
- `Infof`: Log a formatted message at the `info` level.
- `Error`: Log a message at the `error` level.
- `Errorf`: Log a formatted message at the `error` level.

License
-------
Copyright (c) 2014 Ryan Bourgeois. Licensed under BSD-Modified. See the LICENSE
file for a copy of the license.
