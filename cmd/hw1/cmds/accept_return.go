package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/abstractions"
)

func acceptReturnCmd(pvzOrderUseCase abstractions.IPVZOrderUseCase) *cobra.Command {
	command := &cobra.Command{
		Use:     "accept_return",
		Short:   "Accept return",
		Args:    cobra.ExactArgs(2),
		Example: "hw1 accept_return <recipient_id> <order_id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			recipientID := args[0]
			orderID := args[1]

			err := pvzOrderUseCase.AcceptReturn(cmd.Context(), recipientID, orderID)
			if err != nil {
				return err
			}

			cmd.Println("Return accepted")

			return nil
		},
	}

	return command
}
