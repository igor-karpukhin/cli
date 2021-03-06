package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStep(t *testing.T) {
	c := Command{} // uninitialized command

	// test uninitialized command
	require.Panics(t, func() { c.NewStep("Oh noes...") }, "NewStep on uninitialized command should panic.")

	c.Options = NewOptions() // properly initialize command

	// test current step update when creating a new step
	s := c.NewStep("test-step")

	require.Equal(t, s, c.CurrentStep, "Command's current step should be the newly created step.")

}

func TestKubectl(t *testing.T) {
	c := Command{} // uninitialized command

	// test uninitialized command
	require.Panics(t, func() { c.Kubectl() }, "Kubectl wrapper getter on uninitialized command should panic.")

	c.Options = NewOptions() // properly initialize command

	// test lazy init of kubectl wrapper
	k := c.Kubectl()

	require.NotNil(t, k, "Kubectl wrapper should be initialized on demand when getter is called.")

}
