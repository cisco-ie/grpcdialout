#!/bin/sh
#
# Copyright (c) 2018 Cisco
#
# Generate Python code from telemetry and dialout proto files
#
PROTO_ARCHIVE=/opt/git-repos/bigmuddy-network-telemetry-proto/proto_archive

# Check the directory exists
if [ ! -d "$PROTO_ARCHIVE" ]; then
    echo PROTO_ARCHIVE directory $PROTO_ARCHIVE does not exist
    exit 1
fi

export GOPATH=`pwd`

#
# boilerplate install stuff
#
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go

#
# need to put new stuff on the path
#
export PATH=`go env GOPATH`/bin:$PATH

#
# Generate Telemetry message code
#
mkdir -p src/telemetry
protoc \
    -I$PROTO_ARCHIVE \
    $PROTO_ARCHIVE/telemetry.proto \
    --go_out=plugins=grpc:src/telemetry

#
# Generate dialout code
#
mkdir -p src/mdt_dialout
protoc \
    -I$PROTO_ARCHIVE/mdt_grpc_dialout \
    $PROTO_ARCHIVE/mdt_grpc_dialout/mdt_grpc_dialout.proto \
    --go_out=plugins=grpc:src/mdt_dialout
