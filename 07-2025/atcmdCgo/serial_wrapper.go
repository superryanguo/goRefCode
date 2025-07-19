// Package main provides Go wrappers for C serial port functions using cgo.
package main

// #cgo LDFLAGS: -L. -lserial
// #include <stdio.h>
// #include <stdlib.h>
// #include "serial.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// OpenSerialPort wraps the C open_serial_port function to open a serial port device.
func OpenSerialPort(device string) (int, error) {
	// Convert Go string to C string
	cDevice := C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))

	// Call C open_serial_port function
	fd := C.open_serial_port(cDevice)
	if fd == -1 {
		return 0, fmt.Errorf("failed to open serial port %s", device)
	}
	return int(fd), nil
}

// SetSerialPort wraps the C set_serial_port function to configure serial port parameters.
func SetSerialPort(fd int, speed int, bits int, stop int, parity byte) error {
	// Convert parity byte to C char
	cParity := C.char(parity)

	// Call C set_serial_port function
	result := C.set_serial_port(C.int(fd), C.int(speed), C.int(bits), C.int(stop), cParity)
	if result != 0 {
		return fmt.Errorf("failed to set serial port parameters: %d", result)
	}
	return nil
}

// Example usage
func main() {
	// Example: Open and configure a serial port
	device := "/dev/ttyUSB0" // Replace with your serial device
	speed := 9600            // Baud rate
	bits := 8                // Data bits
	stop := 1                // Stop bits
	parity := byte('N')      // Parity: 'N' (none), 'O' (odd), 'E' (even)

	// Open the serial port
	fd, err := OpenSerialPort(device)
	if err != nil {
		fmt.Printf("Error opening serial port: %v\n", err)
		return
	}
	fmt.Printf("Serial port opened, fd: %d\n", fd)

	// Configure the serial port
	err = SetSerialPort(fd, speed, bits, stop, parity)
	if err != nil {
		fmt.Printf("Error setting serial port: %v\n", err)
		return
	}
	fmt.Println("Serial port configured successfully")

	// Example: You can now use fd with read/write operations
	// For instance, integrate with previous read/write wrappers
}