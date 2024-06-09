.PHONY: local-server
local-server:
	go run ./cmd server local

.PHONE: migrate and initialize
migrate:
	go run ./cmd migrate and initialize
