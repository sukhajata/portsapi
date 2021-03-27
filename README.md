# portsapi
An API to retrieve a list of ports
# PortDomainService
To run tests 
```sh
cd service
go test ./...
```

# HTTP Client and file parser
To run tests
```sh
cd client
go test ./...
```
Benchmark and performance
```sh
go test ./... -bench=. -run=XX -benchmem
```
# Launch
docker-compose up --build

