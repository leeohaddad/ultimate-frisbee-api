################################################################################
## Deployment targets
################################################################################

ECR_URL ?= to-be-defined.dkr.ecr.us-east-1.amazonaws.com
ECR_REPO ?= ${ECR_URL}/${APP_NAME}

IMAGE_REPOSITORY ?= ${ECR_REPO}
IMAGE_VERSION ?= $(shell git rev-parse --short HEAD)
IMAGE_TAG ?= $(IMAGE_VERSION)
IMAGE_NAME ?= ${IMAGE_REPOSITORY}:${IMAGE_TAG}

.PHONY: aws/login
aws/login:
	aws-mfa
	aws --region us-east-1 ecr get-login-password | docker login -u AWS --password-stdin ${ECR_URL}

.PHONY: docker/build
docker/build:
	docker pull ${IMAGE_REPOSITORY}:latest || true
	docker build --cache-from ${IMAGE_REPOSITORY}:latest -t ${IMAGE_REPOSITORY} .
	docker tag ${IMAGE_REPOSITORY} ${IMAGE_REPOSITORY}:${IMAGE_TAG}
	docker tag ${IMAGE_REPOSITORY} ${IMAGE_REPOSITORY}:latest

.PHONY: docker/push
docker/push:
	docker push ${IMAGE_REPOSITORY}:${IMAGE_TAG}
	docker push ${IMAGE_REPOSITORY}:latest

.PHONY: deploy/stag
deploy/stag:
	kubectl --context='kube-stag-us-cilium.tfgco.com' -n ultimate-frisbee-manager apply -f .kubernetes/ultimate-frisbee-api-stag.yaml
	kubectl --context='kube-stag-us-cilium.tfgco.com' -n ultimate-frisbee-manager rollout restart deploy ultimate-frisbee-api
