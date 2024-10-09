package e2e

import (
	"github.com/rogpeppe/go-internal/testscript"
	"homework/cmd/cli"
	"os"
	"testing"
)

func Main() int {
	err := main.Run()
	if err != nil {
		return 1
	}
	return 0
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(
		m,
		map[string]func() int{
			"hw": Main,
		},
	))
}

func TestAcceptOrderDelivery(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "scripts/accept_order_delivery",
	})
}

func TestGetOrders(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir: "scripts/get_orders",
	})
}
