# Blockchain Demo Project

A minimal, multi-node blockchain demo built with Go (backend) and React (frontend). This project demonstrates core blockchain concepts, wallet management, mining, and multi-node networking, with a modern UI for interacting with the network.

---

## Project Structure

```
blockchain/
├── app/                  # React frontend (Vite + TypeScript)
│   ├── src/
│   │   ├── components/   # React components (Dashboard, Wallet, Mining, etc.)
│   │   ├── lib/          # Shared types, constants, node list
│   │   ├── services/     # API service functions (fetch chain, wallet, tx, etc.)
│   │   ├── pages/        # (Optional) Page-level components
│   │   └── ...
│   ├── public/           # Static assets
│   ├── package.json      # Frontend dependencies
│   └── ...
├── block/                # Go blockchain core (block, transaction, chain logic)
├── blockchain_server/    # Go Gin server, API handlers, node logic
├── wallet/               # Go wallet generation, signing, etc.
├── utils/                # Go utility functions (ECDSA, config, JSON)
├── Dockerfile            # (Optional) Containerization
├── config.env            # (Optional) Environment config
└── README.md             # This file
```

---

## How to Run

### Prerequisites

-   Go 1.20+
-   Node.js 18+
-   (Optional) Docker

### 1. Start Blockchain Nodes (Backend)

You can run multiple nodes on different ports. Example:

```bash
cd blockchain
# Start node 1
go run blockchain_server/*.go -port 5000
# Start node 2
go run blockchain_server/*.go -port 5001
# Start node 3
go run blockchain_server/*.go -port 5002
```

Just Remember to add this into your config.

### 2. Start the React Frontend

```bash
cd blockchain/app
npm install
npm run dev
```

-   The frontend will run on `localhost:5173` (or as shown in the terminal)
-   By default, it connects to the first node(`gateway_endpoint`), but you can switch nodes in the UI

### 3. Using Docker (optional)

You can build and run both backend and frontend with Docker Compose (add your own `docker-compose.yml` if needed).

---

## Features

-   **Multi-node blockchain:** Run several nodes, each with its own chain and transaction pool
-   **Wallet management:** Create, view, and switch between user and miner wallets
-   **Send crypto:** Transfer coins between wallets
-   **Mining:** Mine new blocks and see rewards in the miner wallet
-   **Chain explorer:** View blocks, transactions, and mempool for any node
-   **Node switching:** Instantly switch between nodes in the frontend
-   **Auto-refresh:** Balances, chain, and transactions pool auto-update

---

## Things To Work On

-   **Tests:** Add unit and integration tests for both backend (Go) and frontend (React)
-   **Dynamic Discovery of Nodes / Registry:** Implement a node registry or peer discovery so nodes can find each other automatically
-   **Improve Error Handling:** Make error messages more user-friendly and robust across the stack
-   **Restructure:** Refactor code for better modularity, maintainability, and scalability

---

## Technical Approach

### Backend (Go)

-   **Gin** for HTTP API
-   **Easily extensible:** Add new endpoints or consensus logic as needed
-   **Node networking:** Each node exposes its own API and can resolve conflicts

### Frontend (React + Vite)

-   **TypeScript** for type safety
-   **TanStack Query** for robust data fetching and caching
-   **Component-driven:** Dashboard, Wallet, Mining, Chain, and Transaction views
-   **Minimal, modern UI:** Clean, responsive, and demo ready
