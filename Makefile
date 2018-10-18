APP=task-mgmt

all: build

build:
	go build -o ${APP} ./cmd/

package: build
	sudo docker build -t ${APP} .

deploy: package
	sudo docker tag ${APP} rivernet/${APP}
	sudo docker push rivernet/${APP}

clean:
	@go clean
	@rm -f ./${APP}
