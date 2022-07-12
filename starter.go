package lite

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/moeryomenko/healing"
	"golang.org/x/sync/errgroup"
)

// defaultContextGracePeriod is default grace period.
// see: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination
const defaultContextGracePeriod = 30 * time.Second

// Lite is compact squad implementation.
type Lite struct {
	group errgroup.Group
	ctx   context.Context

	healthController *healing.Health
}

// New returns new instance of Lite.
func New(healthPort int, opts ...healing.Option) *Lite {
	return &Lite{
		ctx:              context.Background(),
		group:            errgroup.Group{},
		healthController: healing.New(healthPort, opts...),
	}
}

// AddLiveCheck adds subsystem check for checks liveness.
func (l *Lite) AddLiveCheck(subsystem string, check func(context.Context) healing.CheckResult) {
	l.healthController.AddLiveChecker(subsystem, check)
}

// AddReadyCheck adds subsystem check for checks readiness.
func (l *Lite) AddReadyCheck(subsystem string, check func(context.Context) healing.CheckResult) {
	l.healthController.AddReadyChecker(subsystem, check)
}

// Run launch service and health controller.
func (l *Lite) Run(service func(context.Context) error) error {
	l.group.Go(signalListener(l.ctx))

	ctx := WithDelay(l.ctx, defaultContextGracePeriod)

	l.group.Go(func() error {
		return service(ctx)
	})
	l.group.Go(func() error {
		return l.healthController.Heartbeat(ctx)
	})

	err := l.group.Wait()
	_ = l.healthController.Stop(l.ctx)
	return err
}

func signalListener(ctx context.Context) func() error {
	return func() error {
		ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		<-ctx.Done()
		return ctx.Err()
	}
}
