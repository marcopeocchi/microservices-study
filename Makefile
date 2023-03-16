default:
	go build -o fuu -o fuu cmd/server/main.go
	mkdir -p build
	mv fuu* build

clean:
	rm -r build

app: 
	cd cmd/server/frontend && pnpm build

linuxamd64:
	GOOS=linux GOARCH=amd64 go build -o fuu cmd/server/main.go

fuu:
	go run cmd/server/main.go \
		-c "/Users/marco/dev/homebrew/fuu/Fuufile"

fuutest:
	TESTING=true go run cmd/server/main.go \
		-c "/Users/marco/dev/homebrew/fuu/Fuufile"

knight:
	RMQ_ENDPOINT=amqp://user:oseopilota@10.0.0.2:5672/ \
		go run cmd/knight/*.go

perceval:
	JAEGER_ENDPOINT=http://10.0.0.2:14268/api/traces \
		go run cmd/perceval/*.go \
		-c "/Users/marco/dev/homebrew/fuu/cmd/perceval/PercevalFile"
