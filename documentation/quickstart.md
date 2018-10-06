# Quickstart

### dependencies

- [golang](https://golang.org/doc/install)
- [dep](https://github.com/golang/dep#installation)


### cloning

Before you clone this repo, create the appropriate directory for development by running the following command:

```bash
 $ mkdir -p ~/go/src/github.com/angles-n-daemons
```

Then clone the repo to the created directory:

```bash
 $ git clone https://github.com/angles-n-daemons/popsql.git ~/go/src/github.com/angles-n-daemons
```

### running

To run the process, execute the following command:

```bash
 $ go run cmd/main.go
```

### testing

To run the tests, run the following command:

```bash
 $ go test ./...
```

If you would like to test a specific file, simply replace `...` with the relative path to the file name.
