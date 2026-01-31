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
- **Response Body**: `PeerResponse`
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

### 3. Remove Peer
Removes a peer from the WireGuard interface and deletes its persistent metadata.

- **URL**: `/peers/{id}`
- **Method**: `DELETE`
- **Path Parameters**:
    - `id`: The Public Key (ID) of the peer to remove.
- **Response**: `204 No Content`

### 4. Interface Statistics
Returns aggregated real-time statistics for the WireGuard interface.

- **URL**: `/stats`
- **Method**: `GET`
- **Response Body**: `Stats`
    ```json
    {
      "interfaceName": "wg0",
      "peerCount": 5,
      "totalRx": 1048576,
      "totalTx": 2097152
    }
    ```

## Middleware
- **Logging**: All requests are logged in structured JSON format.
- **CORS**: Enabled for all origins (`*`) with methods `GET, POST, DELETE, OPTIONS`.

## Running the Server
The server requires `CAP_NET_ADMIN` to interact with native WireGuard interfaces.
```bash
sudo go run ./cmd/server/main.go
```
If permissions are missing, it will automatically fall back to a **Mock Mode** for development.
