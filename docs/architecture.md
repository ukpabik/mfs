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

### Side note
`FileHandler` uses an internal RW mutex to guard operations.