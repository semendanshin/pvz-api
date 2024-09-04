package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/infrastructure/repositories/pvzorder"
	"homework/internal/usecases"
)

var returnOrderDeliveryCmd = &cobra.Command{
	Use:     "return_delivery",
	Short:   "Return order delivery",
	Args:    cobra.ExactArgs(1),
	Example: "hw1 return_delivery <order_id>",
	RunE: func(cmd *cobra.Command, args []string) error {
		pvzOrderRepository := pvzorder.NewJSONRepository(cmd.Flag("config").Value.String())

		pvzOrderUseCase := usecases.NewPVZOrderUseCase(pvzOrderRepository, cmd.Flag("pvz").Value.String())

		orderID := args[0]

		err := pvzOrderUseCase.ReturnOrderDelivery(orderID)
		if err != nil {
			return err
		}

		cmd.Println("Delivery returned")

		return nil
	},
}
