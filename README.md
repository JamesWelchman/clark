
# Clark

Clark is a daemon implementing the i3bar protocol.

## Overview

Clark implements 100% of the i3bar protocol, including click events.

Clark implements all the vanilla blocks which we would expect in a status bar.
Specifically: Clock/Network/Memory/CPU/Battery.

Clark has no configuration files and is designed to be small and hackable.
It has a suckless style config.h (conf/conf.go).

Clark has an asynchronous architecture.

   - Strictly one goroutine per block
   - Click events are handled in the same goroutine as block writing
