run-server:
	go run ./cmd/main.go server --name="$(ARGS)" --cert-file="./tls.crt" --cert-key-file="./tls.key"

run-client:
	go run ./cmd/main.go client --cert-file="./tls.crt" --cert-key-file="./tls.key" --skip-tls-verify :8080 $(ARGS)