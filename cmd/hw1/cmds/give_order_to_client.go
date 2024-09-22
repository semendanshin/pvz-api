package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/abstractions"
)

func giveOrderToClientCmd(pvzOrderUseCase abstractions.IPVZOrderUseCase) *cobra.Command {
	command := &cobra.Command{
		Use:     "give_orders",
		Short:   "Give orders to client",
		Args:    cobra.MinimumNArgs(1),
		Example: "hw1 give_orders <order_id1> <order_id2> ...",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := pvzOrderUseCase.GiveOrderToClient(args)
			if err != nil {
				return err
			}

			cmd.Println("Orders given to client")

			return nil
		},
	}

	return command
}
