# Contributing

Thanks for contributing to clidle-test. This repository is a sandbox clone of
[`ajeetdsouza/clidle`](https://github.com/ajeetdsouza/clidle), but changes here
follow the usual Go project conventions.

## Prerequisites

- [Go](https://go.dev/dl/) 1.22 or newer (see the `go` directive in `go.mod`)
- `git`

## Getting started

```sh
git clone https://github.com/dispatrick/clidle-test.git
cd clidle-test
go mod download
go build ./...
```

To run the game locally over SSH:

```sh
go run .
```

Then connect to it from another terminal:

```sh
ssh localhost -p 3000
```

## Running the tests

Run the whole suite from the repository root:

```sh
go test ./...
```

While iterating, you can scope the run to a single package or test:

```sh
go test -v ./store
go test -run TestName ./store
```

Note that the repository currently has no test files, so `go test ./...`
reports `[no test files]` for every package. New tests are welcome — add
`*_test.go` files alongside the code they cover.

Before opening a pull request, also check formatting and vet:

```sh
gofmt -l .
go vet ./...
```

`gofmt -l .` should print nothing.

## Submitting changes

1. Branch off `main`.
2. Keep commits focused, with a short descriptive message.
3. Make sure `go build ./...`, `go test ./...`, and `go vet ./...` all pass.
4. Open a pull request describing what changed and why, linking the related
   issue.
