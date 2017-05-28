PID      = /tmp/logvoyage-api.pid
GO_FILES = $(wildcard *.go)
APP      = ./logvoyage
serve: restart
	@fswatch -o . | xargs -n1 -I{}  make restart || make kill
kill:
	@kill `cat $(PID)` || true
before:
	@echo "actually do nothing"
$(APP): $(GO_FILES)
	@go build -o logvoyage
restart: kill before $(APP)
	@./logvoyage start api & echo $$! > $(PID)
test:
	@LV_MODE=test go test -v

.PHONY: serve restart kill before
