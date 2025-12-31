package files

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	writeStr = "hello world"
	testFile = "temp.txt"
	root     = "./data"
)

func TestFileHandler(t *testing.T) {
	handler := NewFileHandler(root)

	err := handler.Create(testFile)
	require.NoError(t, err)
	defer handler.Delete(testFile)
	defer handler.Clear()

	n, err := handler.Write(testFile, bytes.NewReader([]byte(writeStr)))
	require.NoError(t, err)
	require.Equal(t, n, len(writeStr))

	data, err := handler.Read(testFile, 0)
	require.NoError(t, err)
	require.Equal(t, "hello world", string(data))
}
