BINARY_NAME=prometheus_libvirt_exporter

build:
	go build -o .build/${BINARY_NAME}

run:
	CGO_ENABLED=0 go build -o .build/${BINARY_NAME}
	./.build/${BINARY_NAME}

clean:
	go clean
	rm -r ./.build/
