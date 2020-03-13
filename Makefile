PROJECT_NAME=db
CMD_ROOT=db

-include ./makefile.properties

deps:
	go mod vendor -v
	go mod tidy -v
dev:
	USER_ID=$$(id -u) docker-compose -f ./deploy/dev/docker-compose.yaml up -d
	docker-compose -f ./deploy/dev/docker-compose.yaml logs -f
dev_stop:
	docker-compose -f ./deploy/dev/docker-compose.yaml down
dev_clean: dev_stop
	docker-compose -f ./deploy/dev/docker-compose.yaml rm
run:
	go run ./cmd/$(CMD_ROOT)
test:
	go test -v ./... -cover -coverprofile c.out
build:
	go build -o ./bin/$(CMD_ROOT) ./cmd/$(CMD_ROOT)_${GOOS}_${GOARCH}
build_production:
	CGO_ENABLED=0 \
	go build -a -v \
		-ldflags "-X main.Commit=$(GIT_COMMIT) \
			-X main.Version=$(GIT_TAG) \
			-X main.Timestamp=$(TIMESTAMP) \
			-extldflags 'static' \
			-s -w" \
		-o ./bin/$(CMD_ROOT)_$$(go env GOOS)_$$(go env GOARCH)${BIN_EXT} \
		./cmd/$(CMD_ROOT)
compress_production:
	ls -lah ./bin/$(CMD_ROOT)_$$(go env GOOS)_$$(go env GOARCH)${BIN_EXT}
	upx -9 -v -o ./bin/.$(CMD_ROOT)_$$(go env GOOS)_$$(go env GOARCH)${BIN_EXT} \
		./bin/$(CMD_ROOT)_$$(go env GOOS)_$$(go env GOARCH)${BIN_EXT}
	upx -t ./bin/.$(CMD_ROOT)_$$(go env GOOS)_$$(go env GOARCH)${BIN_EXT}
	rm -rf ./bin/$(CMD_ROOT)_$$(go env GOOS)_$$(go env GOARCH)${BIN_EXT}
	mv ./bin/.$(CMD_ROOT)_$$(go env GOOS)_$$(go env GOARCH)${BIN_EXT} \
		./bin/$(CMD_ROOT)_$$(go env GOOS)_$$(go env GOARCH)${BIN_EXT}
	ls -lah ./bin/$(CMD_ROOT)_$$(go env GOOS)_$$(go env GOARCH)${BIN_EXT}

package:
	docker build --file ./deploy/Dockerfile --tag $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest .
save:
	mkdir -p ./build
	docker save --output ./build/db.tar.gz $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest
load:
	docker load --input ./build/db.tar.gz
publish_dockerhub:
	docker push $(DOCKER_NAMESPACE)/$(DOCKER_IMAGE_NAME):latest

see_ci:
	xdg-open https://gitlab.com/usvc/modules/go/db/pipelines

.ssh:
	mkdir -p ./.ssh
	ssh-keygen -t rsa -b 8192 -f ./.ssh/id_rsa -q -N ""
	cat ./.ssh/id_rsa | base64 -w 0 > ./.ssh/id_rsa.base64
