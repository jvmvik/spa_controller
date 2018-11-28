# 
#

file_list=web.go thermometer.go io.go application.go
remote=pi@spa.local:/home/pi/spa_controller

all: run

install:
	go get -u github.com/mileusna/crontab
	go get -u github.com/stianeikeland/go-rpio

run:
	go run $(file_list)

build:
	env GOOS=linux GOARM=6 GOARCH=arm go build -o spa_controller $(file_list)

copy:
	scp -r *.go index.html static/ Makefile $(remote)

deploy: build
	scp spa_controller index.html api.html $(remote)

ssh:
	ssh pi@spa.local
