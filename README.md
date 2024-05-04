# Brainslurp

is pre alpha and very much work in progress.

## Building / Running it locally

You just need go v1.22 or later for this.

`go run cmd/brainslurp/main.go`

## Development

**Requirements**:

* go v1.22 or later [https://go.dev/dl/](https://go.dev/dl/)
* templ ( `go install github.com/a-h/templ/cmd/templ@v0.2.663` )
* protoc v5.27.0--rc1 [https://github.com/protocolbuffers/protobuf](https://github.com/protocolbuffers/protobuf/releases/tag/v27.0-rc1)
* protoc-gen-go v1.33.0 ( `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0` )
* ( optional but recommended ) air ( `go install github.com/cosmtrek/air@latest` )

After changes to `.proto` files just run `./gen-proto.sh` from the repo root.

With air you can just run `air` in the repo root to start the server and have a file watcher that restarts on
file changes. 