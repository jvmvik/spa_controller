all: run


install:
	go get -u github.com/mileusna/crontab
	go get -u github.com/stianeikeland/go-rpio

run:
	go run web.go thermometer.go io.go application.go

build:
	go build web.go thermometer.go io.go application.go

