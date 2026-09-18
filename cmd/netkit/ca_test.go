package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/biancarosa/netkit/internal/proxy"
	"github.com/stretchr/testify/require"
)

func TestCAGenerationDoesNotOverwrite(t *testing.T) {
	dir := t.TempDir()
	cert, key := filepath.Join(dir, "ca.pem"), filepath.Join(dir, "key.pem")
	args := []string{"--cert", cert, "--key", key}
	require.NoError(t, runCA(args))
	_, err := proxy.LoadCA(cert, key)
	require.NoError(t, err)
	original, err := os.ReadFile(key)
	require.NoError(t, err)
	st, err := os.Stat(key)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0600), st.Mode().Perm())
	require.Error(t, runCA(args))
	unchanged, err := os.ReadFile(key)
	require.NoError(t, err)
	require.Equal(t, original, unchanged)
	otherKey := filepath.Join(dir, "other-key.pem")
	require.Error(t, runCA([]string{"--cert", cert, "--key", otherKey}))
	_, err = os.Stat(otherKey)
	require.True(t, os.IsNotExist(err))
}
