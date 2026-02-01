# WireGuard Manager API Documentation

This document describes the API endpoints provided by the `wg-manager` backend.

## Base URL

Default: `http://localhost:8080`

## Endpoints

### 1. List All Peers

Returns a list of all WireGuard peers, combining real-time interface data with persistent metadata (names).

- **URL**: `/peers`
- **Method**: `GET`
- **Response Body**: `[]Peer`
  ```json
  [
  	{
  		"id": "publicKey...",
  		"publicKey": "publicKey...",
  		"name": "Peer Name",
  		"endpoint": "1.2.3.4:51820",
  		"allowedIPs": ["10.0.0.2/32"],
  		"lastHandshake": "2026-01-31 12:00:00",
  		"receiveBytes": 1024,
  		"transmitBytes": 2048
  	}
  ]
  ```

### 2. Add/Configure Peer

Adds a new peer to the WireGuard interface and persists its metadata. If no public key is provided, a new key pair will be generated.

- **URL**: `/peers`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
  	"name": "New Peer",
  	"publicKey": "optionalPublicKey",
  	"allowedIPs": ["10.0.0.3/32"]
  }
  ```
- **Response Body (201 Created)**: `PeerResponse`
  ```json
  {
  	"id": "publicKey...",
  	"publicKey": "publicKey...",
  	"name": "New Peer",
  	"allowedIPs": ["10.0.0.3/32"],
  	"config": "[Interface]\nPrivateKey = ...\n\n[Peer]\nPublicKey = ...",
  	"privateKey": "generatedPrivateKey (if applicable)"
  }
  ```
- **Error Responses (400 Bad Request)**:
  - `Name is required`: If the `name` field is empty or whitespace-only.
  - `At least one AllowedIP is required`: If the `allowedIPs` array is empty.
  - `Invalid AllowedIP CIDR: <value>`: If any item in `allowedIPs` is not a valid CIDR notation.

### 3. Remove Peer

Removes a peer from the WireGuard interface and deletes its persistent metadata.

- **URL**: `/peers/{id}`
- **Method**: `DELETE`
- **Path Parameters**:
  - `id`: The Public Key (ID) of the peer to remove.
- **Response**: `204 No Content`

### 4. Update Peer

Allows partial updates to a peer's metadata (name) and configuration (AllowedIPs).

- **URL**: `/peers/{id}`
- **Method**: `PATCH`
- **Path Parameters**:
  - `id`: The Public Key (ID) of the peer to update.
- **Request Body**:
  ```json
  {
  	"name": "Updated Name",
  	"allowedIPs": ["10.0.0.4/32"]
  }
  ```
- **Response Body (200 OK)**: `Peer` (the updated peer object)
- **Error Responses (400 Bad Request)**:
  - `Invalid AllowedIP CIDR: <value>`: If any item in `allowedIPs` is not a valid CIDR notation.
  - `Invalid request body`: If the JSON body is malformed.

### 5. Regenerate Keys

Generates a new WireGuard keypair for an existing peer while preserving its name and allowed IPs.

- **URL**: `/peers/{id}/regenerate-keys`
- **Method**: `POST`
- **Path Parameters**:
  - `id`: The Public Key (ID) of the peer to regenerate keys for.
- **Response Body (200 OK)**: `PeerResponse` (contains the new keypair and config)
  ```json
  {
  	"id": "newPublicKey...",
  	"publicKey": "newPublicKey...",
  	"name": "Same Name",
  	"allowedIPs": ["10.0.0.2/32"],
  	"config": "[Interface]\nPrivateKey = ...\n...",
  	"privateKey": "newPrivateKey..."
  }
  ```

### 6. Interface Statistics

Returns aggregated real-time statistics for the WireGuard interface.

- **URL**: `/stats`
- **Method**: `GET`
- **Response Body**: `Stats`
  ```json
  {
  	"interfaceName": "wg0",
  	"publicKey": "serverPublicKey...",
  	"listenPort": 51820,
  	"subnet": "10.0.0.0/24",
  	"peerCount": 5,
  	"totalRx": 1048576,
  	"totalTx": 2097152
  }
  ```

## Configuration

The backend uses a hybrid configuration system (Twelve-Factor App). It loads defaults from `backend/internal/config/config.json` and supports overrides via a `.env` file or environment variables.

### Environment Variables

| Variable               | Description                         | Default (JSON)      |
| :--------------------- | :---------------------------------- | :------------------ |
| `WG_SERVER_PORT`       | Port for the HTTP server            | `:8080`             |
| `WG_INTERFACE_NAME`    | Name of the WireGuard interface     | `wg0`               |
| `WG_STORAGE_PATH`      | Path to persistent peer metadata    | `./data/peers.json` |
| `WG_SERVER_ENDPOINT`   | Public IP/Domain:Port of the server | `1.2.3.4:51820`     |
| `WG_SERVER_PUBKEY`     | Public Key of the server interface  | (None)              |
| `WG_VPN_SUBNET`       | VPN subnet CIDR                     | `10.0.0.0/24`       |
| `CORS_ALLOWED_ORIGINS` | Comma-separated list of origins     | (Reflective/Dev)    |

## Middleware

- **Logging**: All requests are logged in structured JSON format via `slog`.
- **Graceful Shutdown**: Intercepts `SIGINT`/`SIGTERM` to drained connections and close `wgctrl` safely.
- **CORS**: Configurable origins; defaults to reflecting `Origin` header in development.

## Running the Server

The server requires `CAP_NET_ADMIN` to interact with native WireGuard interfaces.

```bash
sudo go run ./cmd/server/main.go
```

If permissions are missing, it will automatically fall back to a **Mock Mode** for development.
