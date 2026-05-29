package engine

import (
	"testing"
	"time"
)

func TestNodeSetDownLevelLogonNameDoesNotDeadlockWithDataSource(t *testing.T) {
	node := NewNode()
	node.Set(DataSource, NV("HOST01"))

	done := make(chan struct{})
	go func() {
		node.Set(DownLevelLogonName, NV("HOST01\\alice"))
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("setting DownLevelLogonName deadlocked")
	}
}
