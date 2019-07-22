package containers

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/testcontainers/testcontainers-go"
)

// TerminateFn signals that given resource stops using the container.
type TerminateFn func()

type containerSync struct {
	once      sync.Once
	wg        sync.WaitGroup
	container testcontainers.Container
}

var c containerSync

// Start starts given container once and terminates it once all TerminateFn are called.
//nolint
func Start(t *testing.T, ctx context.Context,
	image string, ports []string, waitingFor wait.Strategy) (testcontainers.Container, TerminateFn) {
	t.Helper()

	c.wg.Add(1)
	c.once.Do(func() {
		startTime := time.Now()
		container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        image,
				ExposedPorts: ports,
				WaitingFor:   waitingFor,
			},
			Started: true,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Container startup took: %s", time.Since(startTime))

		go func() {
			c.wg.Wait()
			err := container.Terminate(ctx)
			if err != nil {
				panic(err)
			}
		}()

		c.container = container
	})

	return c.container, func() {
		c.wg.Done()
	}
}

