package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ngyewch/coptonix-rs232-i2c"
	"github.com/urfave/cli/v3"
	"go.bug.st/serial"
)

func doReadI2C(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() != 2 {
		return fmt.Errorf("incorrect number of arguments")
	}
	slaveAddress := cmd.Uint8Arg(slaveAddressArg.Name)
	count := cmd.Uint8Arg(countArg.Name)

	serialPort, err := newSerialPort(ctx, cmd)
	if err != nil {
		return err
	}
	defer func(serialPort serial.Port) {
		_ = serialPort.Close()
	}(serialPort)

	d := coptonixrs232i2c.New(serialPort)

	data, success, err := d.ReadI2C(slaveAddress, count)
	if err != nil {
		return err
	}
	if success {
		for _, v := range data {
			fmt.Printf("0x%02x ", v)
		}
		fmt.Println()
	} else {
		fmt.Println(success)
	}

	return nil
}

func doWriteI2C(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() != 2 {
		return fmt.Errorf("incorrect number of arguments")
	}
	slaveAddress := cmd.Uint8Arg(slaveAddressArg.Name)
	dataString := cmd.StringArg(dataArg.Name)

	data, err := hex.DecodeString(dataString)
	if err != nil {
		return err
	}

	serialPort, err := newSerialPort(ctx, cmd)
	if err != nil {
		return err
	}
	defer func(serialPort serial.Port) {
		_ = serialPort.Close()
	}(serialPort)

	d := coptonixrs232i2c.New(serialPort)

	success, err := d.WriteI2C(slaveAddress, data)
	if err != nil {
		return err
	}
	fmt.Println(success)

	return nil
}

func doScanI2C(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() != 0 {
		return fmt.Errorf("incorrect number of arguments")
	}
	serialPort, err := newSerialPort(ctx, cmd)
	if err != nil {
		return err
	}
	defer func(serialPort serial.Port) {
		_ = serialPort.Close()
	}(serialPort)

	d := coptonixrs232i2c.New(serialPort)

	slaveAddresses, err := d.ScanI2C()
	if err != nil {
		return err
	}
	for _, slaveAddress := range slaveAddresses {
		fmt.Printf("0x%02x ", slaveAddress)
	}
	fmt.Println()

	return nil
}

func doCheckSlaveAddr(ctx context.Context, cmd *cli.Command) error {
	if cmd.NArg() != 1 {
		return fmt.Errorf("incorrect number of arguments")
	}
	slaveAddress := cmd.Uint8Arg(slaveAddressArg.Name)

	serialPort, err := newSerialPort(ctx, cmd)
	if err != nil {
		return err
	}
	defer func(serialPort serial.Port) {
		_ = serialPort.Close()
	}(serialPort)

	d := coptonixrs232i2c.New(serialPort)

	connected, err := d.CheckSlaveAddress(slaveAddress)
	if err != nil {
		return err
	}
	fmt.Println(connected)

	return nil
}
