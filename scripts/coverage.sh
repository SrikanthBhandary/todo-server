#!/bin/bash
go test -v -coverpkg=./service,./router,./entity,./worker,./config -coverprofile cover.out ./...
go tool cover -html cover.out -o cover.html
open cover.html