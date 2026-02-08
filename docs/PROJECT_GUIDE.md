# Project Guide

This guide covers containerization and deployment to Google Kubernetes Engine (GKE).

## 1) Build Docker Images (AMD64)

From repo root:

```bash
# Backend
docker build --platform linux/amd64 -t holdem-backend:latest -f backend/Dockerfile .

# Frontend (inject backend URL)
docker build --platform linux/amd64 \
  --build-arg REST_BACKEND_URL=http://YOUR_BACKEND_URL \
  -t holdem-frontend:latest -f frontend/Dockerfile frontend
```

## 2) Push Images to Google Artifact Registry (recommended)

```bash
# Set your project and region
PROJECT_ID=your-gcp-project
REGION=us-central1

# Create repo once
gcloud artifacts repositories create holdem \
  --repository-format=docker \
  --location=$REGION

# Configure auth
 gcloud auth configure-docker $REGION-docker.pkg.dev

# Tag images
BACKEND_IMAGE=$REGION-docker.pkg.dev/$PROJECT_ID/holdem/holdem-backend:latest
FRONTEND_IMAGE=$REGION-docker.pkg.dev/$PROJECT_ID/holdem/holdem-frontend:latest

docker tag holdem-backend:latest $BACKEND_IMAGE
docker tag holdem-frontend:latest $FRONTEND_IMAGE

# Push
 docker push $BACKEND_IMAGE
 docker push $FRONTEND_IMAGE
```

## 3) Create GKE Cluster

```bash
gcloud container clusters create holdem-cluster \
  --zone us-central1-a \
  --num-nodes 2 \
  --machine-type e2-standard-2

gcloud container clusters get-credentials holdem-cluster --zone us-central1-a
```

## 4) Update K8s Manifests

Edit the image names in:
- `k8s/backend-deployment.yaml`
- `k8s/frontend-deployment.yaml`

Also set `CORS_ORIGINS` in `k8s/backend-deployment.yaml` to your frontend domain.

## 5) Deploy to Kubernetes

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/backend-deployment.yaml
kubectl apply -f k8s/frontend-deployment.yaml

# Expose services (LoadBalancer)
kubectl apply -f k8s/ingress.yaml
```

Check external IPs:

```bash
kubectl get svc -n holdem
```

Use the frontend external IP in the browser. Use the backend external IP when building the frontend image.

## 6) Verify

```bash
# Health check
curl http://BACKEND_IP/healthz
```

## 7) Load Testing (k6)

```bash
# Local backend
BASE_URL=http://localhost:8080 k6 run loadtest/loadtest.js

# Against LoadBalancer backend
BASE_URL=http://BACKEND_IP k6 run loadtest/loadtest.js
```
