.PHONY: all dev clean build env-up env-down run

all: env-down clean env-up run

# dev: build run

# ##### BUILD    Tab 分隔符 不能用空格
# build:
# 	@echo "Build ..."
# 	@go build
# 	@echo "Build done"

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
	@go run main.go


##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /tmp/dragonwell-* dragonwell
	@docker rm -f -v `docker ps -a --no-trunc | grep "dragonwell" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker rmi `docker images --no-trunc | grep "dragonwell" | cut -d ' ' -f 1` 2>/dev/null || true
	@docker volume prune
	@docker network prune
	@echo "Clean up done"
