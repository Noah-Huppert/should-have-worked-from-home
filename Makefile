.PHONY: run

MAIN=main/main.go

# run starts the application
run:
	. ./.env && \
	go run "${MAIN}"
