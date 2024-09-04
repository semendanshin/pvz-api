package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/infrastructure/repositories/pvzorder"
	"homework/internal/usecases"
)

var giveOrderToClientCmd = &cobra.Command{
	Use:     "give_orders",
	Short:   "Give orders to client",
	Args:    cobra.MinimumNArgs(1),
	Example: "hw1 give_orders <order_id1> <order_id2> ...",
	RunE: func(cmd *cobra.Command, args []string) error {
		pvzOrderRepository := pvzorder.NewJSONRepository(cmd.Flag("config").Value.String())

		pvzOrderUseCase := usecases.NewPVZOrderUseCase(pvzOrderRepository, cmd.Flag("pvz").Value.String())

		err := pvzOrderUseCase.GiveOrderToClient(args)
		if err != nil {
			return err
		}

		cmd.Println("Orders given to client")

		return nil
	},
}
