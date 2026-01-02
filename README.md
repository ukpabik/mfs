# 🪄🧙‍♂️ Magic File System

A MAGICAL distributed file system written in Go. 

## Architecture
- Runs a **MetadataManager** (client entrypoint) that replicates operations to **StorageNodes**
- Uses a simple **quorum rule** (majority) to decide success
- Stores files locally on each node under its own `./dataN` directory using a hashed path layout

## 💽 Scripts

```sh
make run      # start server + 3 nodes
make client   # run test client
make test     # run tests
```