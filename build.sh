#!/usr/bin/env sh
#

set -xe

# generate sql queries
sqlc generate &&

# generate templ templates
templ generate &&

# generate tailwind css file
tailwindcss -i ./assets/app.css -o ./assets/dist/css/style.css &&

# build program
go build -o ./tmp/main ./cmd/web/main.go
