package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/infrastructure/repositories/pvzorder"
	"homework/internal/usecases"
)

var acceptReturnCmd = &cobra.Command{
	Use:     "accept_return",
	Short:   "Accept return",
	Args:    cobra.ExactArgs(2),
	Example: "hw1 accept_return <recipient_id> <order_id>",
	RunE: func(cmd *cobra.Command, args []string) error {
		pvzOrderRepository := pvzorder.NewJSONRepository(cmd.Flag("config").Value.String())

		pvzOrderUseCase := usecases.NewPVZOrderUseCase(pvzOrderRepository, cmd.Flag("pvz").Value.String())

		recipientID := args[0]
		orderID := args[1]

		err := pvzOrderUseCase.AcceptReturn(recipientID, orderID)
		if err != nil {
			return err
		}

		cmd.Println("Return accepted")

		return nil
	},
}
