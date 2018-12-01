# 
#

file_list=web.go thermometer.go io.go application.go
root=/home/pi/spa_controller
remote=pi@spa.local
remote_root=$(remote):$(root)

all: run

install:
	go get -u github.com/mileusna/crontab
	go get -u github.com/stianeikeland/go-rpio

run:
	go run $(file_list)

build:
	env GOOS=linux GOARM=6 GOARCH=arm go build -o spa_controller $(file_list)

copy:
	scp -r *.go *.html *.service spa_controller static/ Makefile $(remote_root)

deploy: build
	scp spa_controller index.html api.html $(remote_root)

# open interactive ssh terminal
debug:
	ssh $(remote)

# install as systemd service
install_service:
	@echo "sudo cp $(root)/spa_controller.service /etc/systemd/system/" | ssh $(remote)

# start daemon service
start:
	@echo "sudo systemctl start $(script)" | ssh $(remote)

# stop daemon service
stop:
	@echo "sudo systemctl stop $(script)" | ssh $(remote)

# enable daemon service
enable:
	@echo "sudo systemctl enable  $(script)" | ssh $(remote)
