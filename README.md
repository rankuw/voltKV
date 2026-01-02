# VoltKV: In-Memory Key-Value Store

VoltKV is a high-performance, in-memory key-value store written in Go. It mimics Redis by supporting the **RESP (Redis Serialization Protocol)** and ensures durability using an **Append-Only Log (AOL)** for disk persistence.

---

## 🚀 Project Roadmap

### Phase 1: Networking & Protocol (Current Status)
- [x] **TCP Server**: Accepts concurrent connections.
- [x] **RESP Parser**: Reads raw bytes and converts them into Go structures.
- [x] **Command Handler**: Executes commands (`SET`, `GET`, `PING`).

### Phase 2: Storage Engine
- [x] **In-Memory Store**: Thread-safe Hash Map (`sync.RWMutex`).
- [ ] **Key Expiry**: Logic to TTL (Time To Live) for keys.
g
### Phase 3: Persistence
- [ ] **AOF (Append-Only File)**: Log every write command to disk.
- [ ] **Recovery**: Replay AOF on startup to restore state.

---
