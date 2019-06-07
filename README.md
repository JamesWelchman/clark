
# Clark

Clark is a daemon implementing the i3bar protocol.

[godoc](https://godoc.org/github.com/JamesWelchman/clark)

## Overview

Clark implements 100% of the i3bar protocol, including click events.

Clark implements all the vanilla blocks we would expect in a status bar:

   - Battery Status/Charge Percentage
   - Cpu loads
   - Memory Usage
   - Network I/O
   - Clock

Clark has no configuration files and is designed to be small and hackable.
It has a suckless style config.h (conf/conf.go) and (colors/colors.go).

Clark has an asynchronous architecture.

   - Strictly one goroutine per block
   - Click events are handled in the same goroutine as block writing
