// Package main implements concurrent serial port reading and writing using goroutines.
package main

// #cgo LDFLAGS: -L. -lserial
// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include <unistd.h>
// #include "serial.h"
import "C"
import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
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

// readSerial continuously reads from the serial port and sends data to a channel.
func readSerial(ctx context.Context, fd int, ch chan<- []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data, err := ReadFromFD(fd, RV_LEN)
			if err != nil {
				fmt.Printf("Error reading from serial port: %v\n", err)
				continue
			}
			if len(data) > 0 {
				ch <- data
			}
			time.Sleep(100 * time.Millisecond) // Prevent tight loop
		}
	}
}

// writeSerial reads user input and writes commands to the serial port.
func writeSerial(ctx context.Context, fd int) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter commands (press Ctrl+C to exit):")
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			cmd := scanner.Text() + "\r" // Append carriage return
			if len(cmd) > CMD_LEN {
				fmt.Printf("Command too long, max length is %d\n", CMD_LEN)
				continue
			}
			_, err := WriteToFD(fd, cmd)
			if err != nil {
				fmt.Printf("Error writing to serial port: %v\n", err)
				continue
			}
			// Wait briefly to allow response
			time.Sleep(700 * time.Millisecond)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}

func main() {
	// Check command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("usage:\n\t./serial_cmd <device>")
		os.Exit(1)
	}

	// Get device from arguments
	device := os.Args[1]

	// Open serial port
	fd, err := OpenSerialPort(device)
	if err != nil {
		fmt.Printf("Error opening serial port: %v\n", err)
		os.Exit(1)
	}
	defer CloseFD(fd) // Ensure fd is closed on exit

	// Configure serial port (115200 baud, 8 bits, 1 stop bit, no parity)
	err = SetSerialPort(fd, 115200, 8, 1, 'N')
	if err != nil {
		fmt.Printf("Error setting serial port: %v\n", err)
		os.Exit(1)
	}

	// Create context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel for read results
	readCh := make(chan []byte, 10)

	// Start read goroutine
	go readSerial(ctx, fd, readCh)

	// Start write goroutine
	go writeSerial(ctx, fd)

	// Handle SIGINT (Ctrl+C) for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	// Print read results
	for {
		select {
		case data := <-readCh:
			fmt.Printf("Received: %s\n", string(data))
		case <-sigCh:
			fmt.Println("\nShutting down...")
			cancel()
			time.Sleep(100 * time.Millisecond) // Allow goroutines to exit
			return
		}
	}
}