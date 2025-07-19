// Package main replicates the C program for serial communication using cgo wrappers.
package main

// #cgo LDFLAGS: -L. -lserial
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include <unistd.h>
// #include "serial.h"
import "C"
import (
	"fmt"
	"os"
	"time"
	"unsafe"
)

// Constants
const (
	RV_LEN  = 512 // Buffer size for reading
	CMD_LEN = 512 // Buffer size for command
)

// OpenSerialPort wraps the C open_serial_port function.
func OpenSerialPort(device string) (int, error) {
	cDevice := C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))
	fd := C.open_serial_port(cDevice)
	if fd == -1 {
		return 0, fmt.Errorf("failed to open serial port %s", device)
	}
	return int(fd), nil
}

// SetSerialPort wraps the C set_serial_port function.
func SetSerialPort(fd int, speed int, bits int, stop int, parity byte) error {
	cParity := C.char(parity)
	result := C.set_serial_port(C.int(fd), C.int(speed), C.int(bits), C.int(stop), cParity)
	if result != 0 {
		return fmt.Errorf("failed to set serial port parameters: %d", result)
	}
	return nil
}

// WriteToFD wraps the C write function.
func WriteToFD(fd int, cmd string) (int, error) {
	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))
	n := C.write(C.int(fd), cCmd, C.size_t(C.strlen(cCmd)))
	if n < 0 {
		return 0, fmt.Errorf("write failed: %d", n)
	}
	return int(n), nil
}

// ReadFromFD wraps the C read function.
func ReadFromFD(fd int, bufLen int) ([]byte, error) {
	rv := make([]byte, bufLen)
	for i := range rv {
		rv[i] = 0
	}
	n := C.read(C.int(fd), unsafe.Pointer(&rv[0]), C.size_t(bufLen))
	if n < 0 {
		return nil, fmt.Errorf("read failed: %d", n)
	}
	return rv[:n], nil
}

// CloseFD wraps the C close function.
func CloseFD(fd int) error {
	result := C.close(C.int(fd))
	if result != 0 {
		return fmt.Errorf("failed to close file descriptor: %d", result)
	}
	return nil
}

func main() {
	// Check command-line arguments
	if len(os.Args) < 3 {
		fmt.Println("usage:\n\t./serial_cmd <device> <command>")
		os.Exit(1)
	}

	// Get device and command from arguments
	device := os.Args[1]
	cmd := os.Args[2] + "\r" // Append carriage return, mimicking strcat(cmd, "\r")

	// Open serial port
	fd, err := OpenSerialPort(device)
	if err != nil {
		fmt.Printf("Error opening serial port: %v\n", err)
		os.Exit(1)
	}

	// Configure serial port (115200 baud, 8 bits, 1 stop bit, no parity)
	err = SetSerialPort(fd, 115200, 8, 1, 'N')
	if err != nil {
		fmt.Printf("Error setting serial port: %v\n", err)
		C.close(C.int(fd)) // Close fd on error
		os.Exit(1)
	}

	// Write command to serial port
	_, err = WriteToFD(fd, cmd)
	if err != nil {
		fmt.Printf("Error writing to serial port: %v\n", err)
		C.close(C.int(fd))
		os.Exit(1)
	}

	// Wait 700ms (mimicking usleep(700000))
	time.Sleep(700 * time.Millisecond)

	// Read response
	rv, err := ReadFromFD(fd, RV_LEN)
	if err != nil {
		fmt.Printf("Error reading from serial port: %v\n", err)
		C.close(C.int(fd))
		os.Exit(1)
	}

	// Print response
	fmt.Print(string(rv))

	// Close serial port
	err = CloseFD(fd)
	if err != nil {
		fmt.Printf("Error closing serial port: %v\n", err)
		os.Exit(1)
	}
}