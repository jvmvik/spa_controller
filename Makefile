file_list=web.go thermometer.go io.go application.go
remote=pi@spa.local:/home/pi/spa_controller

all: run

install:
	go get -u github.com/mileusna/crontab
	go get -u github.com/stianeikeland/go-rpio

run:
	go run $(file_list)

build:
	go build $(file_list)

deploy:
	scp -r *.go index.html static/ Makefile $(remote)

release:
	env GOARM=6 GOARCH=arm go build $(file_list) ; scp ./web $(remote)

ssh:
	ssh pi@spa.local