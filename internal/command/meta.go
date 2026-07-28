package command

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/internal/backend"
	"github.com/hashicorp/terraform/internal/states/statemgr"
	"github.com/hashicorp/terraform/internal/views"
)

type stateLocker struct {
	backend backend.Enhanced
	view    views.StateLocker
	context context.Context
	lockID  string
}

func (s *stateLocker) Unlock() error {
	if s.backend == nil || s.lockID == "" {
		return nil
	}

	// Use a detached context with a timeout to ensure the unlock operation
	// can complete even if the main execution context has been canceled.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s.view.StateUnlockStart(s.lockID)
	err := s.backend.Unlock(ctx, s.lockID)
	if err != nil {
		s.view.StateUnlockFailed(s.lockID, err)
		return fmt.Errorf("failed to release state lock: %w", err)
	}
	s.view.StateUnlockSuccess(s.lockID)
	return nil
}
