package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v3"
)

var (
	version string

	serialPortFlag = &cli.StringFlag{
		Name:     "serial-port",
		Usage:    "serial port",
		Required: true,
		Sources:  cli.EnvVars("SERIAL_PORT"),
	}
	baudRateFlag = &cli.IntFlag{
		Name:    "baud-rate",
		Usage:   "baud rate",
		Value:   19200,
		Sources: cli.EnvVars("BAUD_RATE"),
	}
	dataBitsFlag = &cli.IntFlag{
		Name:    "data-bits",
		Usage:   "data bits",
		Value:   8,
		Sources: cli.EnvVars("DATA_BITS"),
	}
	parityFlag = &cli.StringFlag{
		Name:    "parity",
		Usage:   "parity",
		Value:   "N",
		Sources: cli.EnvVars("PARITY"),
		Action: func(ctx context.Context, cmd *cli.Command, s string) error {
			switch s {
			case "N", "O", "E", "M", "S":
				return nil
			default:
				return fmt.Errorf("unknown parity: %s", s)
			}
		},
	}
	stopBitsFlag = &cli.StringFlag{
		Name:    "stop-bits",
		Usage:   "stop bits",
		Value:   "1",
		Sources: cli.EnvVars("STOP_BITS"),
		Action: func(ctx context.Context, cmd *cli.Command, s string) error {
			switch s {
			case "1", "1.5", "2":
				return nil
			default:
				return fmt.Errorf("unknown stop bits: %s", s)
			}
		},
	}
	readTimeoutFlag = &cli.DurationFlag{
		Name:    "read-timeout",
		Usage:   "read timeout",
		Value:   5 * time.Second,
		Sources: cli.EnvVars("READ_TIMEOUT"),
	}

	slaveAddressArg = &cli.Uint8Arg{
		Name: "slave-address",
	}
	countArg = &cli.Uint8Arg{
		Name: "count",
	}
	dataArg = &cli.StringArg{
		Name: "data",
	}

	app = &cli.Command{
		Name:    "coptonix-rs232-i2c",
		Usage:   "Coptonix RS232 I2C CLI",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:   "scan-i2c",
				Usage:  "scan I2C",
				Action: doScanI2C,
			},
			{
				Name:  "check-slave-addr",
				Usage: "check slave addr",
				Arguments: []cli.Argument{
					slaveAddressArg,
				},
				Action: doCheckSlaveAddr,
			},
			{
				Name:  "read-i2c",
				Usage: "read I2C",
				Arguments: []cli.Argument{
					slaveAddressArg,
					countArg,
				},
				Action: doReadI2C,
			},
			{
				Name:  "write-i2c",
				Usage: "write I2C",
				Arguments: []cli.Argument{
					slaveAddressArg,
					dataArg,
				},
				Action: doWriteI2C,
			},
		},
		Flags: []cli.Flag{
			serialPortFlag,
			baudRateFlag,
			dataBitsFlag,
			parityFlag,
			stopBitsFlag,
			readTimeoutFlag,
		},
	}
)

func main() {
	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
