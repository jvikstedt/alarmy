# Alarmy
Alarmy is under development.

## Development
#### Install Dependencies
Creates a vendor folder and downloads dependencies
```sh
$ glide install
```
#### Run server
```sh
$ go run main.go server
```
#### Running tests
```sh
$  go test $(go list ./... | grep -v /vendor/) -v
```
