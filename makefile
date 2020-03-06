build:
	go build -mod vendor -o ./bin/influx_example main.go
run: build
	./bin/influx_example