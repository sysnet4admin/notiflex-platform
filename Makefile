# Makefile for notiflex-platform

# --- Variables ---
# Image details
IMAGE_NAME ?= api
GCP_PROJECT ?= project-75fce205-dfa5-4975-a56
GCP_REGION ?= asia-northeast3
REPO_NAME ?= notiflex
IMAGE_REGISTRY = $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT)/$(REPO_NAME)

# Use the short git hash as the default tag.
TAG ?= $(shell git rev-parse --short HEAD)
IMAGE_TAG = $(TAG)
IMAGE_URL = $(IMAGE_REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

# Kubernetes
K8S_DIR = k8s/smb
K8S_NAMESPACE = notiflex

# Go
APP_DIR = app

# --- Targets ---

.PHONY: all build push deploy

all: build push deploy
	@echo "Build, push, and deploy complete for image: $(IMAGE_URL)"

# Build the Docker image
build:
	@echo "Building Docker image: $(IMAGE_URL)"
	docker build -t $(IMAGE_URL) -f $(APP_DIR)/Dockerfile $(APP_DIR)

# Push the Docker image to the registry
push:
	@echo "Pushing Docker image: $(IMAGE_URL)"
	docker push $(IMAGE_URL)

# Deploy the application to Kubernetes
deploy:
	@echo "Deploying to Kubernetes in namespace $(K8S_NAMESPACE)"
	@# Update the image tag in the deployment yaml
	sed -i.bak "s|image: .*/api:.*|image: $(IMAGE_URL)|" $(K8S_DIR)/deployment.yaml
	kubectl apply -f $(K8S_DIR)/
	@echo "Deployment updated with image: $(IMAGE_URL)"
	@# Optional: clean up the backup file
	rm $(K8S_DIR)/deployment.yaml.bak

# Target to print the image URL
show-image-url:
	@echo $(IMAGE_URL)
