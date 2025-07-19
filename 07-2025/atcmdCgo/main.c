#include<unistd.h>
#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include"serial.h"
#define RV_LEN 512
#define CMD_LEN 512

int main(int argc, char *argv[])
{
    if(argc < 3)
    {
	printf("usag:\n\taa dev cmd \n");
	exit(1);
    }

    int dev;
    char cmd[CMD_LEN], rv[RV_LEN];
    strcpy(cmd, argv[2]);
    strcat(cmd, "\r");

    dev = open_serial_port(argv[1]);
    set_serial_port(dev, 115200, 8, 1, 'N');
    write(dev, cmd, strlen(cmd));
    usleep(700000);
    bzero(rv, RV_LEN);
    read(dev, rv, RV_LEN);
    printf("%s", rv);
    close(dev);

    return 0;
}

