# Texas Hold'em Lab

A Texas Hold'em evaluator and odds calculator built with **Go** backend and **Flutter** frontend.  
Designed for containerization and deployment on **Google Kubernetes Engine (GKE)** with AMD64 nodes, for the course 'Distributed Systems'.


## Architecture

```
┌─────────────────┐     HTTP/JSON      ┌──────────────────────┐
│  Flutter Web    │ ─────────────────► │     Go Backend       │
│  (Frontend)     │  /api/v1/* (REST)  │     :8080 (REST)     │
└─────────────────┘                    └──────────────────────┘
```

- **Backend**: Go REST API on port 8080
- **Frontend**: Flutter web app calling the REST API
- **API Modes**:
  - Best hand from 7 cards
  - Heads-up comparison
  - Monte Carlo odds

## Prerequisites

- **Go** 1.22+
- **Flutter** (for web)
- **Docker** (for images)
- **kubectl** + **gcloud** (for GKE)

## Project Structure

```
.
├── assets/
│   └── test_cases/           # CSV test cases for comparisons
├── backend/
│   ├── cmd/api/              # Go REST server
│   ├── internal/poker/        # Hand evaluation + Monte Carlo
│   └── internal/api/          # HTTP handlers
├── frontend/
│   ├── lib/main.dart          # Flutter web UI
│   └── pubspec.yaml
├── Makefile
└── README.md
```

## Local Development

### Backend

```bash
cd backend
go run ./cmd/api
```

The server runs on `http://localhost:8080`.

### Frontend

If this is a fresh Flutter folder:

```bash
cd frontend
flutter create .
flutter pub get
flutter run
```

Set the API base URL in the UI to `http://localhost:8080`.

## Makefile Commands

```bash
make backend       # build Go backend
make backend-run   # run backend binary
make frontend      # build Flutter web app
make docker-build  # build amd64 Docker images
make test          # run Go tests
make clean         # remove build artifacts
```
