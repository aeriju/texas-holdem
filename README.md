# Texas Hold'em Lab

A Texas Hold'em evaluator and odds calculator built with **Go** backend and **Flutter** frontend.  
Designed for containerization and deployment on **Google Kubernetes Engine (GKE)** with AMD64 nodes, for the course 'Distributed Systems'.



## Tech Stack

| Layer     | Technology                   |
|----------|-------------------------------|
| Frontend | Flutter / Dart                |
| Backend  | Go (REST API)                 |
| Runtime  | Docker, GKE (AMD64)           |

## Card Format

Cards are 2-character strings:
- **First char**: Suit – `H` (Hearts), `D` (Diamonds), `C` (Clubs), `S` (Spades)
- **Second char**: Rank – `A` (Ace), `K` (King), `Q` (Queen), `J` (Jack), `T` (Ten), `9`–`2`

Examples: `HA` (Heart Ace), `S7` (Spade 7), `CT` (Club Ten)


## Project Structure

```
.
├── assets/
│   └── test_cases/             # CSV test cases for comparisons
├── backend/
│   ├── cmd/api/                # Go REST server
│   └── internal/
│       ├── api/                # HTTP handlers
│       └── poker/              # Hand evaluation + Monte Carlo
├── docs/
│   └── PROJECT_GUIDE.md        # Docker + GKE deployment steps
├── frontend/
│   ├── lib/
│   │    ├── main.dart          # App entry + theme
│   │    ├── models/            # Shared enums/types
│   │    ├── screens/           # UI screens
│   │    ├── services/          # REST client
│   │    └── widgets/           # Reusable UI widgets
│   └── pubspec.yaml
├── k8s/
│   ├── backend-deployment.yaml
│   ├── frontend-deployment.yaml
│   ├── ingress.yaml
│   └── namespace.yaml
├── Makefile
└── README.md
```


## Quick Start

### Local Development

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

Backend URL is injected at build/run time:

```bash
flutter run --dart-define=REST_BACKEND_URL=http://localhost:8080
```

### Deployment

See `docs/PROJECT_GUIDE.md`.

### Makefile Commands

```bash
make backend       # build Go backend
make backend-run   # run backend binary
make frontend      # build Flutter web app
make docker-build  # build amd64 Docker images of both backend and frontend
make loadtest      # run k6 load tests
make test          # run Go tests
make clean         # remove build artifacts
```

### API Endpoints
| Method | Endpoint            | Description                                           |
|--------|---------------------|-------------------------------------------------------|
| POST   | `/api/v1/best-hand` | Best hand from 2 hole + 5 community cards             |
| POST   | `/api/v1/heads-up`  | Compare two hands, return winner                      |
| POST   | `/api/v1/odds`      | Win probability via Monte Carlo simulation            |


## References
- [Texas Hold'em (Wikipedia)](https://en.wikipedia.org/wiki/Texas_hold_%27em)
- [Poker Hand Rankings](https://en.wikipedia.org/wiki/List_of_poker_hands)
