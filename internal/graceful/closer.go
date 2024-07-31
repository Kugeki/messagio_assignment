package graceful

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"messagioassignment/internal/logger"
	"sync"
	"time"
)

type Func func(ctx context.Context) error

// Closer performs graceful shutdown via Closer.Shutdown
// with close functions that can added with Closer.Add
type Closer struct {
	mu        sync.Mutex
	functions []Func
	log       *slog.Logger
}

func NewCloser(log *slog.Logger) *Closer {
	return &Closer{
		mu:        sync.Mutex{},
		functions: make([]Func, 0),
		log:       log.With(slog.String("component", "graceful/closer")),
	}
}

func (c *Closer) Add(shutdownFunc Func) {
	c.mu.Lock()
	c.functions = append(c.functions, shutdownFunc)
	c.mu.Unlock()
}

func (c *Closer) Shutdown(shutdownTimeout time.Duration, goLimit int) {
	c.log.Info("shutting down app gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := c.Close(shutdownCtx, goLimit); err != nil {
		c.log.Error("graceful close", logger.Err(err))
		return
	}

	c.log.Info("app was successfully shut down gracefully!")
}

func (c *Closer) Close(ctx context.Context, goLimit int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		errs   = make([]error, 0, len(c.functions))
		errsMu = sync.Mutex{}
		done   = make(chan struct{}, 1)
	)

	g := errgroup.Group{}
	if goLimit > 0 {
		g.SetLimit(goLimit)
	}

	go func() {
		for _, f := range c.functions {
			g.Go(func() error {
				err := f(ctx)
				if err != nil {
					errsMu.Lock()
					errs = append(errs, err)
					errsMu.Unlock()
				}
				return nil
			})
		}

		err := g.Wait()
		if err != nil {
			errsMu.Lock()
			errs = append(errs, err)
			errsMu.Unlock()
		}
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	case <-done:
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown finished with errors: %w",
			errors.Join(errs...))
	}

	return nil
}
