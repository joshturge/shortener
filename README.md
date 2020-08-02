# URL Shortener
A URL shortener made in Go. This program hosts a web front-end and uses either an in-memory or a redis
key store for storing shortened URL's. 

## Configuration
All of the programs configuration is done through command line
arguments. The default configuration variables are the following:

| Field Name         | Value          |
|--------------------|----------------|
| Listen Address     | localhost:8080 |
| Redis Address      | localhost:6379 |
| Redis Password     | None 		  |
| Test Mode          | False          |
| Timeout            | 5 Minutes   	  |

To see usage see `./shortener -h` for more details

## Building

### POSIX
I have provided a [Makefile](Makefile) which can be used on POSIX operating systems.
The file comes with two targets, the first being `build` which will build the project
and the second target `clean` will just remove the bin directory and call `go clean`. Example usage:

```bash
make build
```

**NOTE**: The build target will place the resulting binary within the `bin/` directory

### Windows
Use the command below to build on windows:

```bash
mkdir bin/
go build -o bin/shortener.exe cmd/main.go
```

## Unit Tests
Most of the packages have unit tests which test the core functionality of the program. The unit tests are
automatically done within the [Makefile](Makefile) when running or building the program. For windows you
may need to run the command below:
```bash
go test -v ./pkg/...
```

## Built With
* [go-redis](https://github.com/go-redis/redis)
* [xxhash](https://github.com/cespare/xxhash)

## License
This project is licensed under the BSD License - see the [LICENSE](LICENSE)
file for details.
