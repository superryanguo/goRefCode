SHELL:=/bin/zsh
.PHONY: dhcpscan

export GO111MODULE=on

# Path Related
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
RELEASE_DIR := ${MKFILE_DIR}bin
VERSION=$(shell git describe --tags --dirty || echo "unknown version")
BDTIME=$(shell date --utc "+%Y-%m-%d %H:%M:%S" || echo "unknown date")
 #-X "github.com/wireless/wirelesscpe/utils.BuildTime=$(BDTIME)"
#-X "github.com/wireless/wirelesscpe/utils.Version=$(VERSION)" 
 #-ldflags '-w -s -buildid=$(VERSION)'
GOBUILD=go build -v -trimpath


# Image Name
NAME=dhcpscan
# Version
RELEASE?=v0.1

# Git Related
GIT_REPO_INFO=$(shell cd ${MKFILE_DIR} && git config --get remote.origin.url)
ifndef GIT_COMMIT
  GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif

# Rules

dhcpscan: ${MKFILE_DIR}/dhcpscan/dhcpscan.go
	@echo "building... dhcpscan"
	cd ${MKFILE_DIR}/dhcpscan && \
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(RELEASE_DIR)/$@

clean:
	rm -rf ${RELEASE_DIR}

ftp: dhcpscan
	scp ${RELEASE_DIR}/dhcpscan  root@192.168.1.1:~/

tftp:
	cd ${RELEASE_DIR} && tftp -m binary 172.21.1.94 -c put dhcpscan

run: dhcpscan
	sudo ${RELEASE_DIR}/dhcpscan
