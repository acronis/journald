# journald
[![Go Reference](https://pkg.go.dev/badge/github.com/acronis/journald.svg)](https://pkg.go.dev/github.com/acronis/journald)
[![Build Status](https://github.com/acronis/journald/actions/workflows/go.yml/badge.svg)](https://github.com/acronis/journald/actions/workflows/go.yml)

Package `journald` offers Go implementation of systemd Journal's native API for logging. Key features are:

* based on a connection-less socket
* work with messages of any size and type
* client can use any number of separate sockets

Requires Linux with systemd and Go 1.25 or newer.

## Installation

Install the package with:

```shell
go get github.com/acronis/journald
```

## Usage: The Best Way

The Best Way to use structured logs (systemd Journal, etc.) is [logf](https://github.com/acronis/logf) - the fast, asynchronous, structured logger in Go with zero allocation count and its journald driver [logfjournald](https://github.com/acronis/logfjournald). This driver uses `journald` package.
The following example creates the new `logf` logger with `logfjournald` appender.

```go
package main

import (
    "runtime"

    "github.com/acronis/logf"
    "github.com/acronis/logfjournald"
)

func main() {
    // Create journald Appender with default journald Encoder.
    appender, appenderClose := logfjournald.NewAppender(logfjournald.NewEncoder.Default())
    defer appenderClose()

    // Create ChannelWriter with journald Encoder.
    writer, writerClose := logf.NewChannelWriter(logf.ChannelWriterConfig{
        Appender: appender,
    })
    defer writerClose()

    // Create Logger with ChannelWriter.
    logger := logf.NewLogger(logf.LevelInfo, writer)

    logger.Info("got cpu info", logf.Int("count", runtime.NumCPU()))
}
```

The JSON representation of the journal entry this generates:

```json
{
  "TS": "2018-11-01T07:25:18Z",
  "PRIORITY": "6",
  "LEVEL": "info",
  "MESSAGE": "got cpu info",
  "COUNT": "4"
}
```

## Usage: AS-IS

Let's look at what the `journald` provides as Go APIs for logging:

```go
package main

import (
    "github.com/acronis/journald"
)

func main() {
    journald.Print(journald.PriorityInfo, "Hello World!")
}
```

The JSON representation of the journal entry this generates:

```json
{
    "PRIORITY": "6",
    "MESSAGE":  "Hello World!",
    "_PID":     "3965",
    "_COMM":    "simple",
    "...":      "..."
}
```

The primary reason for using the Journal's native logging APIs is not just plain logs: it is to allow passing additional structured log messages from the program into the journal. This additional log data may then be used to search the journal for, is available for consumption for other programs, and might help the administrator to track down issues beyond what is expressed in the human readable message text. Here's an example how to do that with `journald.Send`:

```go
package main

import (
    "fmt"
    "os"
    "runtime"

    "github.com/acronis/journald"
)

func main() {
    err := journald.Send("Hello World!", journald.PriorityInfo, map[string]interface{}{
        "HOME":        os.Getenv("HOME"),
        "TERM":        os.Getenv("TERM"),
        "N_GOROUTINE": runtime.NumGoroutine(),
        "N_CPUS":      runtime.NumCPU(),
        "TRACE":       runtime.ReadTrace(),
    })
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

This will write a log message to the journal much like the earlier examples. However, this time a few additional, structured fields are attached:

```json
{
    "PRIORITY":     "6",
    "MESSAGE":      "Hello World!",
    "HOME":         "/root",
    "TERM":         "xterm",
    "N_GOROUTINE":  "2",
    "N_CPUS":       "4",
    "TRACE":        [103,111,32,49,46,56,32,116,114,97,99,101,0,0,0,0],
    "_PID":         "4037",
    "_COMM":        "send",
    "...":          "..."
}
```

Our structured message includes seven fields. The first two we passed are well-known fields:

1. `MESSAGE=` is the actual human readable message part of the structured message.
1. `PRIORITY=` is the numeric message priority value as known from BSD syslog formatted as an integer string.

Applications may relatively freely define additional fields as they see fit (we defined five pretty arbitrary ones in our example). A complete list of the currently well-known fields is available [here](https://www.freedesktop.org/software/systemd/man/systemd.journal-fields.html).

## Your own Journal

`journald.Print` and `journald.Send` write through `journald.DefaultJournal`. Create your own `journald.Journal` when you need a separate socket or different settings:

```go
j := &journald.Journal{
    // Field names are passed through this function before sending.
    // Journal accepts only uppercase names, so this lets you write "home" instead of "HOME".
    NormalizeFieldNameFn: strings.ToUpper,
}
defer j.Close()

err := j.Send("Hello World!", journald.PriorityInfo, map[string]interface{}{"home": os.Getenv("HOME")})
```

Two more knobs:

* `TestModeEnabled: true` makes the journal discard all messages, handy in tests.
* `journald.IsNotExist()` reports whether the system journal socket is missing, for example in a container without systemd.

## Releases

Release notes are on the [releases page](https://github.com/acronis/journald/releases).

## License

MIT, see [LICENSE](LICENSE).

> Thanks to [http://0pointer.de/blog/](http://0pointer.de/blog/) for the inspiration.
