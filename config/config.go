package config

import "flag"

type Config struct {
	ProgrammingCommand []string
	SerialPort         string
	BaudRate           uint
	LocalEcho          bool
}

func (c *Config) ParseFlags() {
	serialPort := flag.String("port", "/dev/ttyACMA0", "Serial Port for serial terminal")
	baudRate := flag.Uint("baud", 115200, "Serial Port Baud rate for serial terminal")
	flag.Parse()
	c.ProgrammingCommand = flag.Args() // Remaining command is the programming command from the user
	c.SerialPort = *serialPort
	c.BaudRate = *baudRate
	c.LocalEcho = false
}
