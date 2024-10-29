package testcomponents

import (
	"context"

	"github.com/go-kit/log"
	"github.com/grafana/alloy/internal/component"
	"github.com/grafana/alloy/internal/featuregate"
	"github.com/grafana/alloy/internal/runtime/logging/level"
	"go.uber.org/atomic"
)

func init() {
	component.Register(component.Registration{
		Name:      "testcomponents.summation1",
		Stability: featuregate.StabilityPublicPreview,
		Args:      SummationConfig_1{},
		Exports:   SummationExports_1{},

		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			return NewSummation(opts, args.(SummationConfig))
		},
	})
}

type SummationConfig_1 struct {
	Input int `alloy:"input,attr"`
	//TODO: What should the type be?
	ForwardTo []string `alloy:"forward_to,attr"`
}

type SummationExports_1 struct {
	Sum       int `alloy:"sum,attr"`
	LastAdded int `alloy:"last_added,attr"`
}

type Summation_1 struct {
	opts component.Options
	log  log.Logger
	sum  atomic.Int32
}

// NewSummation creates a new summation component.
func NewSummation_1(o component.Options, cfg SummationConfig) (*Summation, error) {
	t := &Summation{opts: o, log: o.Logger}
	if err := t.Update(cfg); err != nil {
		return nil, err
	}
	return t, nil
}

var (
	_ component.Component = (*Summation)(nil)
)

// Run implements Component.
func (t *Summation_1) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Update implements Component.
func (t *Summation_1) Update(args component.Arguments) error {
	c := args.(SummationConfig)
	newSum := int(t.sum.Add(int32(c.Input)))

	level.Info(t.log).Log("msg", "updated sum", "value", newSum, "input", c.Input)
	t.opts.OnStateChange(SummationExports{Sum: newSum, LastAdded: c.Input})

	//TODO: Send the data over to all the receivers listed in forward_to

	return nil
}
