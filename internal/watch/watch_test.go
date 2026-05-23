package watch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/watch"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, ".env")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestSnapshot_ReturnsStates(t *testing.T) {
	p := writeTempEnv(t, "KEY=val\n")
	states, err := watch.Snapshot([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(states) != 1 {
		t.Fatalf("expected 1 state, got %d", len(states))
	}
	if states[0].Path != p {
		t.Errorf("expected path %s, got %s", p, states[0].Path)
	}
}

func TestSnapshot_MissingFile(t *testing.T) {
	_, err := watch.Snapshot([]string{"/no/such/file.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestChanged_NoChange(t *testing.T) {
	p := writeTempEnv(t, "KEY=val\n")
	prev, _ := watch.Snapshot([]string{p})
	next, _ := watch.Snapshot([]string{p})
	if watch.Changed(prev, next) {
		t.Error("expected no change")
	}
}

func TestChanged_DetectsChange(t *testing.T) {
	p := writeTempEnv(t, "KEY=val\n")
	prev, _ := watch.Snapshot([]string{p})
	// Ensure mod time advances
	time.Sleep(10 * time.Millisecond)
	if err := os.WriteFile(p, []byte("KEY=new\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	next, _ := watch.Snapshot([]string{p})
	if !watch.Changed(prev, next) {
		t.Error("expected change to be detected")
	}
}

func TestWatch_CallsOnChange(t *testing.T) {
	p := writeTempEnv(t, "KEY=val\n")
	opts := watch.DefaultOptions()
	opts.Interval = 20 * time.Millisecond

	called := make(chan struct{}, 1)
	done := make(chan struct{})

	go func() {
		_ = watch.Watch([]string{p}, opts, func() {
			select {
			case called <- struct{}{}:
			default:
			}
		}, done)
	}()

	time.Sleep(30 * time.Millisecond)
	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case <-called:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Error("onChange was not called within timeout")
	}
	close(done)
}
