# Project Guide

This guide covers containerization and deployment to Google Kubernetes Engine (GKE).

## 1) Build and Push Backend Image (AMD64)

From repo root:

```bash
docker build --platform linux/amd64 -t holdem-backend:latest -f backend/Dockerfile .
```

Tag and push to Artifact Registry:

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

# Tag backend image
BACKEND_IMAGE=$REGION-docker.pkg.dev/$PROJECT_ID/holdem/holdem-backend:latest

docker tag holdem-backend:latest $BACKEND_IMAGE

# Push
docker push $BACKEND_IMAGE
```

## 2) Create GKE Cluster

```bash
gcloud container clusters create holdem-cluster \
  --zone us-central1-a \
  --num-nodes 2 \
  --machine-type e2-standard-2

gcloud container clusters get-credentials holdem-cluster --zone us-central1-a
```

## 3) Deploy Backend and Get Backend IP

Edit the backend image in `k8s/backend-deployment.yaml`:

```yaml
image: REGISTRY/holdem-backend:latest
```

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/backend-deployment.yaml
kubectl apply -f k8s/ingress.yaml
```

Check external IPs:

```bash
kubectl get svc -n holdem
```

Copy the backend LoadBalancer IP (service `holdem-backend-lb`).

## 4) Build and Push Frontend Image (with backend IP)

```bash
BACKEND_IP=http://YOUR_BACKEND_LB_IP
docker build --platform linux/amd64 \
  --build-arg REST_BACKEND_URL=$BACKEND_IP \
  -t holdem-frontend:latest -f frontend/Dockerfile frontend
```

Tag and push:

```bash
FRONTEND_IMAGE=$REGION-docker.pkg.dev/$PROJECT_ID/holdem/holdem-frontend:latest
docker tag holdem-frontend:latest $FRONTEND_IMAGE
docker push $FRONTEND_IMAGE
```

## 5) Deploy Frontend, then Update CORS

Edit the frontend image in `k8s/frontend-deployment.yaml`:

```yaml
image: REGISTRY/holdem-frontend:latest
```

Apply:

```bash
kubectl apply -f k8s/frontend-deployment.yaml
```

Get the frontend external IP:

```bash
kubectl get svc -n holdem
```

Update CORS in `k8s/backend-deployment.yaml`:

```yaml
- name: CORS_ORIGINS
  value: "http://YOUR_FRONTEND_LB_IP"
```

Apply again:

```bash
kubectl apply -f k8s/backend-deployment.yaml
kubectl rollout restart deployment/holdem-backend -n holdem
```

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
