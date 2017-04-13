package commands

import (
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
	"github.com/urfave/cli"
)

func TestSetup(t *testing.T) {
	g := Goblin(t)

	// Gomega fail handler
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Setup command", func() {

		setup := &Setup

		g.Describe("Structure", func() {
			g.It("must have a name `setup`", func() {
				Expect(setup.Name).To(Equal("setup"))
			})

			g.It("must have a `force` flag", func() {
				Expect(setup.Flags).To(HaveLen(1))
				flag := setup.Flags[0]

				boolFlag, ok := flag.(cli.BoolFlag)
				if ok {
					Expect(boolFlag.Name).To(Equal("force, f"))
				} else {
					g.Fail("force flag must be a bool type flag")
				}
			})
		})
	})

}
