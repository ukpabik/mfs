package manager

import (
	"fmt"
	"log"
	"sync"

	"github.com/ukpabik/mfs/internal/storage"
	"github.com/ukpabik/mfs/internal/transport"
)

// MetadataManager is the coordinator for a distributed file system.
// It receives client requests, replicates them across storage nodes,
// and returns results based on quorum consensus.
type MetadataManager struct {
	ID        string
	Transport transport.Transport

	StorageNodes []*storage.StorageNode
	stopCh       chan struct{}
}

func NewMetadataManager(id string, transport transport.Transport, nodes []*storage.StorageNode) *MetadataManager {
	return &MetadataManager{
		ID:           id,
		Transport:    transport,
		StorageNodes: nodes,
		stopCh:       make(chan struct{}),
	}
}

// loop receives client RPCs, broadcasts to replicas, and sends back results.
func (mm *MetadataManager) loop() {
	for {
		select {
		case rpc := <-mm.Transport.Consume():
			message, err := storage.ParseMessage(&rpc)
			if err != nil {
				_ = rpc.Peer.Send([]byte("ERROR: " + err.Error()))
				continue
			}

			log.Printf("mm[%s] client asked for %v on %s", mm.ID, message.Action, message.FilePath)

			resp := mm.broadcast(message)

			_ = rpc.Peer.Send(resp)
		case <-mm.stopCh:
			return
		}
	}
}

type BroadcastResult struct {
	Data []byte
	Err  error
}

// broadcast sends the operation to all storage nodes in parallel
// and returns success only if quorum (majority) of nodes succeed.
func (mm *MetadataManager) broadcast(msg storage.StorageNodeMessage) []byte {

	var wg sync.WaitGroup
	results := make(chan BroadcastResult, len(mm.StorageNodes))

	for _, node := range mm.StorageNodes {
		wg.Add(1)
		go func(n *storage.StorageNode) {
			defer wg.Done()
			data, err := n.HandleMessage(msg)
			results <- BroadcastResult{Data: data, Err: err}
		}(node)
	}

	wg.Wait()
	close(results)

	successCount := 0
	var responseData []byte
	for res := range results {
		if res.Err == nil {
			successCount++
			if responseData == nil {
				responseData = res.Data
			}
		}
	}

	quorumSize := (len(mm.StorageNodes) / 2) + 1
	if successCount >= quorumSize {
		return responseData
	}

	return fmt.Appendf(nil, "ERROR: only %d/%d nodes succeeded", successCount, len(mm.StorageNodes))
}

func (mm *MetadataManager) Start() error {
	go mm.loop()

	if err := mm.Transport.ListenAndAccept(); err != nil {
		return err
	}

	return nil
}
func (mm *MetadataManager) Close() error {
	close(mm.stopCh)
	return mm.Transport.Close()
}
