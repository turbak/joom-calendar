build:
	CGO_ENABLED=0 && go build cmd/joom-calendar bin/joom-calendar

run:
	CGO_ENABLED=0 && go run cmd/joom-calendar