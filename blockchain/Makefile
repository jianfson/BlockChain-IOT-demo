.PHONY: all dev clean build env-up env-down run

all: env-down clean env-up build run

# dev: build run

##### BUILD    用 tab 不能用空格
build:
	@echo "Build ..."
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd ./fixtures && docker-compose up --force-recreate -d
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd ./fixtures && docker-compose down
	@echo "Environment down"

##### RUN
run:
	@./blc


##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/dragonwell-* dragonwell
	@docker rm -f -v `docker ps -a --no-trunc | grep "dragonwell" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "dragonwell" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker volume prune
	@docker network prune
	@echo "Clean up done"
