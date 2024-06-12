.PHONY: local-server
local-server:
	go run ./cmd server local

.PHONE: migrate
migrate:
	go run ./cmd migrate
