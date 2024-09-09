package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/abstractions"
	"homework/internal/infrastructure/repositories/pvzorder"
	"homework/internal/usecases"
)

var getOrdersCmd = &cobra.Command{
	Use:     "get_orders",
	Short:   "Get orders",
	Args:    cobra.ExactArgs(1),
	Example: "hw1 get_orders <user_id>",
	RunE: func(cmd *cobra.Command, args []string) error {
		pvzOrderRepository := pvzorder.NewJSONRepository(cmd.Flag("orders").Value.String())

		pvzOrderUseCase := usecases.NewPVZOrderUseCase(pvzOrderRepository, cmd.Flag("pvz").Value.String())

		userID := args[0]

		opts := make([]abstractions.GetOrdersOptFunc, 0)
		if lastN, _ := cmd.Flags().GetInt("lastN"); lastN > 0 {
			opts = append(opts, abstractions.WithLastNOrders(lastN))
		} else {
			page, _ := cmd.Flags().GetInt("page")
			pageSize, _ := cmd.Flags().GetInt("pageSize")

			paginationOpts, err := abstractions.NewPaginationOptions(
				abstractions.WithPage(page),
				abstractions.WithPageSize(pageSize),
			)
			if err != nil {
				return err
			}

			opts = append(opts, abstractions.WithPaginationOptions(paginationOpts))
		}

		if samePVZ, _ := cmd.Flags().GetBool("samePVZ"); samePVZ {
			opts = append(opts, abstractions.WithPVZID(cmd.Flag("pvz").Value.String()))
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

func init() {
	getOrdersCmd.Flags().Int("page", 0, "page")
	getOrdersCmd.Flags().Int("pageSize", 10, "page size")
	getOrdersCmd.Flags().Int("lastN", 0, "last N")
	getOrdersCmd.Flags().Bool("samePVZ", false, "same PVZ")
}
