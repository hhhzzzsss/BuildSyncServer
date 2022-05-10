# Build Sync Server
 The local server that automatically loads and displays builds made in golang (or anything that spits out a valid REGION_DUMP file).

This repo includes the golang code that generates the various regions. Each folder in `GeneratorPrograms/cmds` contains a separate entrypoint.
My workflow usually involves opening the `GeneratorPrograms` in VSCode, writing whatever build code I want, and then running it with a command like `go ./cmds/blue_mandelbulb`.
