.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: run
run:
	go run ./cmd/main.go
