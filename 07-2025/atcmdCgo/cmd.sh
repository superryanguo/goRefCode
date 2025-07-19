gcc -c serial.c -o serial.o
ar rcs libserial.a serial.o
go run serial_cmd.go /dev/ttyUSB0 "test command"
