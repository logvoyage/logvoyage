PID      = /tmp/awesome-golang-project.pid
GO_FILES = $(wildcard *.go)
APP      = ./api
serve: restart
	@fswatch -o . | xargs -n1 -I{}  make restart || make kill

kill:
	@kill `cat $(PID)` || true

before:
	@echo "actually do nothing"
$(APP): $(GO_FILES)
	@go build -o api
restart: kill before $(APP)
	@./api & echo $$! > $(PID)
test:
	@LV_MODE=test go test -v

.PHONY: serve restart kill before # let's go to reserve rules names
