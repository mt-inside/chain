run:
	go run .

build:
        go build .

install-static:
        go install -ldflags '-w -linkmode external -extldflags "-static"'

package:
        docker build -t mtinside/chain:0.0.1 .

push:
	docker push mtinside/chain:0.0.1
