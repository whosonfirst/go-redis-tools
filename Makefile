tools:
	go build -mod vendor -o bin/publish cmd/publish/main.go
	go build -mod vendor -o bin/subscribe cmd/subscribe/main.go
	go build -mod vendor -o bin/pubsubd cmd/pubsubd/main.go

fmt:
	# go fmt cmd/*.go
	go fmt pubsub/*.go
	go fmt resp/*.go

pub:
	./bin/publish -redis-channel debug -pubsubd -debug -

sub:
	./bin/subscribe -redis-channel debug
