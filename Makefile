DOCKER_REGISTRY=index.docker.io
DOCKER_IMAGE=mrecco/vault-backup
DOCKER_TAG=v1.0.0

TEST_ADDR=http://me:8200
TEST_TOKEN=s.111111111111111111111111

build:
	@docker build -t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG) .

push: auth
	@docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)

rmi:
	@docker rmi $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)

auth:
	@aws ecr get-login | sed "s/\ -e\ none//g" | sh

debug:
	@go run ./code/*.go -addr="$(TEST_ADDR)" -token=$(TEST_TOKEN) -authtype=token backup GI
	# @go run ./code/*.go -addr="$(TEST_ADDR)" -token=$(TEST_TOKEN) -authtype=token restore bkp.yml
