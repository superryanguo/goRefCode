// Package main provides a Go wrapper for C read and write functions using cgo.
package main

// #include <stdio.h>
// #include <stdlib.h>
// #include <string.h>
// #include <unistd.h>
import "C"
import (
	"unsafe"
)

// Constants
const RV_LEN = 1024 // Buffer size for reading, matching typical C usage

// WriteToFD wraps the C write function to write a string to a file descriptor.
func WriteToFD(fd int, cmd string) (int, error) {
	// Convert Go string to C string
	cCmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(cCmd))

	// Call C write function
	n := C.write(C.int(fd), cCmd, C.size_t(C.strlen(cCmd)))
	if n < 0 {
		return 0, fmt.Errorf("write failed: %d", n)
	}
	return int(n), nil
}

// ReadFromFD wraps the C read function to read from a file descriptor into a buffer.
func ReadFromFD(fd int, bufLen int) ([]byte, error) {
	// Allocate buffer
	rv := make([]byte, bufLen)
	
	// Zero the buffer (emulating bzero)
	for i := range rv {
		rv[i] = 0
	}

	// Call C read function
	n := C.read(C.int(fd), unsafe.Pointer(&rv[0]), C.size_t(bufLen))
	if n < 0 {
		return nil, fmt.Errorf("read failed: %d", n)
	}
	return rv[:n], nil
}

// Example usage
func main() {
	// Example: Writing to and reading from a file descriptor
	// For demonstration, using stdout (fd=1)
	cmd := "Hello, cgo!\n"
	fd := 1 // stdout

	// Write to file descriptor
	n, err := WriteToFD(fd, cmd)
	if err != nil {
		fmt.Printf("Error writing: %v\n", err)
		return
	}
	fmt.Printf("Wrote %d bytes\n", n)

	// Note: Reading from stdout is not typical; this is just for demonstration.
	// In real usage, you might read from a device file or socket.
	data, err := ReadFromFD(fd, RV_LEN)
	if err != nil {
		fmt.Printf("Error reading: %v\n", err)
		return
	}
	fmt.Printf("Read %d bytes: %s\n", len(data), string(data))
}