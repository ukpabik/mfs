# Magic File System Architecture

## Current state (single-node)

### Storage
Files are stored on the local filesystem under `FileHandler.RootDir`.

User input is treated as a **key** and mapped to an on-disk path with `TransformPath(key)`, which produces:

`<hashSeg1>/<hashSeg2>/<hashSeg3>/<baseName>`

Hash segments come from `sha256(key)` (hex-encoded). `baseName` is derived from `filepath.Base(key)`.

### File operations
Implemented by `internal/files.FileHandler`:

- `Create(key)` creates an empty file (and parent dirs).
- `Write(key, data io.Reader)` overwrites the file contents.
- `Read(key, size)` reads all bytes if `size == 0`, otherwise reads up to `size` bytes.
- `Delete(key)` removes the file.
- `Clear()` deletes the entire root directory.

### Networking (TCP transport)
Implemented by `internal/server` (transport layer).

- `TCPTransport` listens on a TCP address and accepts connections.
- Each connection is wrapped as a `TCPPeer`.
- Optional `HandshakeFunc` runs on connect.
- Incoming messages are decoded into an `RPC{From, Payload}` and pushed into `TCPTransport.rpcChannel`.
- Consumers read inbound RPCs via `Consume() <-chan RPC`.