package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/abstractions"
)

func returnOrderDeliveryCmd(pvzOrderUseCase abstractions.IPVZOrderUseCase) *cobra.Command {
	command := &cobra.Command{
		Use:     "return_delivery",
		Short:   "Return order delivery",
		Args:    cobra.ExactArgs(1),
		Example: "hw1 return_delivery <order_id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			orderID := args[0]

			err := pvzOrderUseCase.ReturnOrderDelivery(cmd.Context(), orderID)
			if err != nil {
				return err
			}

			cmd.Println("Delivery returned")

			return nil
		},
	}

	return command
}
