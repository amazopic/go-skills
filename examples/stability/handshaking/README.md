# Handshaking Example

Client and server negotiate capacity before work begins. The server exposes
a `/capacity` endpoint advertising available slots and a lease TTL; the client
caches the response for the lease duration and refuses to send when `Slots == 0`.

## Structure

| File | Purpose |
|---|---|
| `handshaking.go` | `Server` (capacity endpoint + enforced limit), `Client` (lease-cached handshake) |
| `handshaking_test.go` | Tests covering slot advertisement, lease caching, re-fetch on expiry, and integration |

## Run

```bash
go test -race ./stability/handshaking/
```

## Key points

- The server enforces the limit server-side too — the client handshake reduces wasted RPCs but is not the only guard.
- Lease caching avoids a `/capacity` round-trip on every request while keeping the advertisement fresh.
- `LeaseTTL=0` disables caching — suitable when capacity changes rapidly (e.g. auction-style allocation).
- The canonical production example is HTTP/2 `SETTINGS_MAX_CONCURRENT_STREAMS`, which uses the same push-based advertisement mechanism.
