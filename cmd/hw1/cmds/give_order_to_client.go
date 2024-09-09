package cmds

import (
	"github.com/spf13/cobra"
)

var giveOrderToClientCmd = &cobra.Command{
	Use:     "give_orders",
	Short:   "Give orders to client",
	Args:    cobra.MinimumNArgs(1),
	Example: "hw1 give_orders <order_id1> <order_id2> ...",
	RunE: func(cmd *cobra.Command, args []string) error {
		ordersFile, _ := cmd.Flags().GetString("orders")
		pvzID, _ := cmd.Flags().GetString("pvz")

		pvzOrderUseCase := InitUseCase(ordersFile, pvzID)

		err := pvzOrderUseCase.GiveOrderToClient(args)
		if err != nil {
			return err
		}

		cmd.Println("Orders given to client")

		return nil
	},
}
