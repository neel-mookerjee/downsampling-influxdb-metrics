AWS_ACCOUNT	:= "12345678901"
IMAGE_NAME	:= "metrics-downsampling-job"
REPOSITORY_NAME	:= "$(IMAGE_NAME)"
ECR_REPOSITORY	:= "$(AWS_ACCOUNT).dkr.ecr.us-west-2.amazonaws.com/$(REPOSITORY_NAME)"

check-var = $(if $(strip $($1)),,$(error var for "$1" is empty))
STACK_LOWER := $(shell echo $(stack) | tr A-Z a-z)
QUERY_ID := $(shell echo $(queryid))
RELEASE := metrics-downsampling-job-$(STACK_LOWER)-$(QUERY_ID)

default: help

require_tag:
	$(call check-var,tag)

require_stack:
	$(call check-var,stack)

require_queryid:
	$(call check-var,queryid)

go/compile:             ## compile go programs
						@CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix cgo -o bin/main .

docker/tags:            ## list the existing tagged images
						@aws ecr list-images --registry-id $(AWS_ACCOUNT) --repository-name $(REPOSITORY_NAME) --filter tagStatus=TAGGED | jq -M -r '.imageIds|.[]|.imageTag'

docker/build:           validate_tag ## build and tag the Docker image. vars:tag
						@docker build -t $(IMAGE_NAME) .
						@docker tag $(REPOSITORY_NAME) $(ECR_REPOSITORY):$(tag)

validate_tag:           require_tag
						#@aws ecr list-images --registry-id $(AWS_ACCOUNT) --repository-name $(REPOSITORY_NAME) --filter tagStatus=TAGGED | jq -M -r '.imageIds|.[]|.imageTag' | tr '\n' ' ' | grep -q -v $(tag)[^-] || { echo "error using the tag"; exit 1;}

docker/push:            validate_tag ## push the Docker image to ECR. vars:tag
						@aws ecr get-login --region us-west-2 | sh -
						@docker push $(ECR_REPOSITORY):$(tag)

helm/install:           require_stack require_queryid ## Deploy the stack into kubernetes. vars: stack, queryid (e.g. stack=test queryid=q37f89d)
						@helm install --name $(RELEASE) --values=./chart/values/values-$(STACK_LOWER).yaml --set QUERY_ID=Q-$(QUERY_ID) ./chart

helm/delete:            require_stack require_queryid ## delete stack from reference. vars: stack, queryid (e.g. stack=test queryid=q37f89d)
						@helm delete $(RELEASE) --purge

helm/reinstall:         require_stack require_queryid ## delete stack from reference and then deploy. vars: stack, queryid (e.g. stack=test queryid=q37f89d)
						@helm delete $(RELEASE) --purge
						@helm install --name $(RELEASE) --values=./chart/values/values-$(STACK_LOWER).yaml --set QUERY_ID=Q-$(QUERY_ID) ./chart

deploy:                 require_tag require_stack require_queryid ## Compiles, builds and deploys a stack for a tag. vars: tag, stack, queryid (e.g. tag=latest stack=test queryid=q37f89d)
						@CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix cgo -o bin/main .
						@docker build  -t $(IMAGE_NAME) .
						@docker tag $(REPOSITORY_NAME) $(ECR_REPOSITORY):$(tag)
						@docker push $(ECR_REPOSITORY):$(tag)
						@aws ecr get-login --region us-west-2 | sh -
						@helm install --name $(RELEASE) --values=./chart/values/values-$(STACK_LOWER).yaml --set QUERY_ID=Q-$(QUERY_ID) ./chart

redeploy:               require_tag require_stack require_queryid ## Compiles, builds and re-deploys a stack for a tag. vars: tag, stack, queryid (e.g. tag=latest stack=test queryid=q37f89d)
						@CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix cgo -o bin/main .
						@docker build  -t $(IMAGE_NAME) .
						@docker tag $(REPOSITORY_NAME) $(ECR_REPOSITORY):$(tag)
						@aws ecr get-login --region us-west-2 | sh -
						@docker push $(ECR_REPOSITORY):$(tag)
						@helm delete $(RELEASE) --purge
						@helm install --name $(RELEASE) --values=./chart/values/values-$(STACK_LOWER).yaml --set QUERY_ID=Q-$(QUERY_ID) ./chart

help:                   ## this helps
						@awk 'BEGIN {FS = ":.*?## "} /^[\/a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36mmetrics-downsampling-job \033[0m%-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)
