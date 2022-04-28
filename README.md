uProg

UART programming wrapper tool.
This integrates a _very_ basic serial port terminal, handling closing and opening the serial port around a programming tool.

When it is first run it will act as a serial port terminal viewer tool; but then running the programming command will close the viewer, and execute the programming command.
Once the programming command is finished, any key can be pressed to drop back into the serial port viewer.

## Usage

By default, the program will use `/dev/ttyACMA0` and a baud rate of `115200`.
(These can both be changed via command line flags obviously)

The program expects key combinations to be pressed to perform actions.
At the moment the sequence is Ctrl-k, followed by another key.

- `p` -> Runs the programming command
- `c` -> Closes the serial port
- `o` -> Opens the serial port
- `l` -> Clears the current serial terminal

The program will parse for arguments until `--` after which, all following arguments will be interpreted as the programmer command to use in programming mode.

For example:
```
uprog --port /dev/ttyUSB0 -- bflb_mcu_tool --firmware test.bin --port /dev/ttyUSB0
```

This will by default open a serial monitor on /dev/ttyUSB0 at 115200 baud.
Once `Ctrl-k,p` is pressed it will close the port, then run `bflb_mcu_tool --firmware test.bin --port /dev/ttyUSB0` in the current folder and show its output.
When the command finishes, pressing any key will close the log and return to the serial data viewer


## Planned features:
- [ ] DTS/RTS control if possible
- [ ] Config file support for saving commands
- [ ] Nicer overlay menu for configuartion without restarting