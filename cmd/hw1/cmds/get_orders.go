package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/abstractions"
)

func getOrdersCmd(pvzOrderUseCase abstractions.IPVZOrderUseCase) *cobra.Command {
	command := &cobra.Command{
		Use:     "get_orders",
		Short:   "Get orders",
		Args:    cobra.ExactArgs(1),
		Example: "hw1 get_orders <user_id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			userID := args[0]

			opts := make([]abstractions.GetOrdersOptFunc, 0)
			if lastN, _ := cmd.Flags().GetInt("lastN"); lastN > 0 {
				opts = append(opts, abstractions.WithLastNOrders(lastN))
			}

			if samePVZ, _ := cmd.Flags().GetBool("samePVZ"); samePVZ {
				opts = append(opts, abstractions.WithPVZID(cmd.Flag("pvz").Value.String()))
			}

			if cursorID, _ := cmd.Flags().GetString("cursorID"); cursorID != "" {
				opts = append(opts, abstractions.WithCursorID(cursorID))
			}

			if limit, _ := cmd.Flags().GetInt("limit"); limit > 0 {
				opts = append(opts, abstractions.WithLimit(limit))
			}

			data, err := pvzOrderUseCase.GetOrders(userID, opts...)
			if err != nil {
				return err
			}

			cmd.Println("Orders:")
			for _, order := range data {
				cmd.Println(order)
			}

			return nil
		},
	}

	command.Flags().Int("lastN", 0, "last N")
	command.Flags().Bool("samePVZ", false, "same PVZ")
	command.Flags().String("cursorID", "", "cursor ID")
	command.Flags().Int("limit", 10, "limit")

	return command
}
