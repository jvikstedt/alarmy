# Alarmy
Alarmy is under development.

## Development
#### Run server
```sh
$ go run main.go server
```
#### Running tests
```sh
$  go test $(go list ./... | grep -v /vendor/) -v
```
