# Gedis

Gedis is a simple in-memory key-value store written in Go, similar to Redis. It supports basic Redis operations like `SET`, `GET`, `PING`, and `ECHO`.

## Usage

### Method 1: Download the Pre-built Application

You can download pre-built versions of the application for Windows or Linux from the [releases](LINK_TO_RELEASES_PAGE) page.

### Method 2: Build and Run from Source

To run the server from the source code, you'll need to have [Go installed](https://go.dev/dl/) on your system. Once Go is installed, navigate to the directory containing the "main.go" file in your terminal and execute the following command:
```shell
go run main.go
```

## Comands

- PING - Tests if the service is running. If the service is online, it prints "PONG". Use it as follows:
```shell
% gedis-cli PING
```
- ECHO - Prints the provided message back to the user:
```shell
% gedis-cli ECHO Hello, World!
```
- SET - Sets the provided key to the provided value. It supports setting an optional expiry time for the key. The expiry can be set in seconds using EX flag or milliseconds using PX flag:
```shell
% gedis-cli SET myKey myValue // Without expiry time
% gedis-cli SET myKey myValue EX 10 // With expiry time in seconds
% gedis-cli SET myKey myValue PX 1000 // With expiry time in milliseconds
```
- GET - Retrieves the value of the provided key:
```shell
% gedis-cli GET myKey
```
Do note that if a key with an expiry time is accessed after it has expired, it will return (nil).