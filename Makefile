.PHONY: run
run:
	go run ./cmd/main.go

.PHONY: templ-generate
templ-generate:
	templ generate