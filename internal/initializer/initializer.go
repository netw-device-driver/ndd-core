package initializer

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// New returns a new *Initializer.
func New(kube client.Client, steps ...Step) *Initializer {
	return &Initializer{kube: kube, steps: steps}
}

// Step is a blocking step of the initialization process.
type Step interface {
	Run(ctx context.Context, kube client.Client) error
}

// Initializer makes sure the objects ndd reconciles are ready to go before
// starting core ndd routines.
type Initializer struct {
	steps []Step
	kube  client.Client
}

// Init does all operations necessary for controllers and webhooks to work.
func (c *Initializer) Init(ctx context.Context) error {
	for _, s := range c.steps {
		if err := s.Run(ctx, c.kube); err != nil {
			return err
		}
	}
	return nil
}
