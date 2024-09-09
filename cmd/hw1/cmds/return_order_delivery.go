package cmds

import (
	"github.com/spf13/cobra"
)

var returnOrderDeliveryCmd = &cobra.Command{
	Use:     "return_delivery",
	Short:   "Return order delivery",
	Args:    cobra.ExactArgs(1),
	Example: "hw1 return_delivery <order_id>",
	RunE: func(cmd *cobra.Command, args []string) error {
		ordersFile, _ := cmd.Flags().GetString("orders")
		pvzID, _ := cmd.Flags().GetString("pvz")

		pvzOrderUseCase := InitUseCase(ordersFile, pvzID)

		orderID := args[0]

		err := pvzOrderUseCase.ReturnOrderDelivery(orderID)
		if err != nil {
			return err
		}

		cmd.Println("Delivery returned")

		return nil
	},
}
