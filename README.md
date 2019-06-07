
# Clark

Clark is a daemon implementing the i3bar protocol.

Docs: [godoc](https://godoc.org/github.com/JamesWelchman/clark)

Build: [![CircleCI](https://circleci.com/gh/JamesWelchman/clark/tree/master.svg?style=svg)](https://circleci.com/gh/JamesWelchman/clark/tree/master)

## Overview

Clark implements 100% of the i3bar protocol, including click events.

Clark implements all the vanilla blocks we would expect in a status bar:

   - Battery Status/Charge Percentage
   - Cpu loads
   - Memory Usage
   - Network I/O
   - Clock
   - Pacman Updates Block

Clark has no configuration files and is designed to be small and hackable.
It has a suckless style config.h (conf/conf.go) and (colors/colors.go).

Clark has an asynchronous architecture.

   - Strictly one goroutine per block
   - Click events are handled in the same goroutine as block writing


## Install
To install clark

```bash
$ go get github.com/jameswelchman/clark
```

This should clone this repo and build a clark binary.
See documentation about go get for more details.
It *probably* installed clark to `$HOME/go/bin/clark`.


## Debugging
All errors are written to stderr.
This is a snippet from my i3 config. Note the commented out line.

```
bar {
   status_command /home/james/go/bin/clark
   # status_command /home/james/bin/clark_debug.sh
}
```

And the contents of `/home/james/bin/clark_debug.sh`.

```bash
#!/usr/bin/env bash

/home/james/go/bin/clark 2> /tmp/clark_error.log
```


## TODO
   1. i3 current-window and current-layout block
   2. PulseAudio block
