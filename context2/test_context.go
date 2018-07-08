package context2

import (
	"context"
	"testing"

	"github.com/localhots/gobelt/log"
)

// TestContext returns a new context with test name added.
func TestContext(t *testing.T) context.Context {
	return log.ContextWithFields(context.Background(), log.F{"test": t.Name()})
}
