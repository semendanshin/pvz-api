package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/infrastructure/handlers/bubbletea"
	"homework/internal/infrastructure/repositories/pvzorder"
	"homework/internal/usecases"
)

var rootCmd = &cobra.Command{
	Use:   "hw1",
	Short: "Homework 1",
	Run: func(cmd *cobra.Command, args []string) {
		pvzOrderRepository := pvzorder.NewJSONRepository(cmd.Flag("orders").Value.String())

		pvzOrderUseCase := usecases.NewPVZOrderUseCase(pvzOrderRepository, cmd.Flag("pvz").Value.String())

		handler := bubbletea.NewHandler(pvzOrderUseCase)

		err := handler.Run()
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("orders", "o", "orders.json", "orders file")
	rootCmd.PersistentFlags().StringP("pvz", "p", "1", "pvz id")

	rootCmd.AddCommand(acceptDeliveryCmd)
	rootCmd.AddCommand(acceptReturnCmd)
	rootCmd.AddCommand(getOrdersCmd)
	rootCmd.AddCommand(getReturnsCmd)
	rootCmd.AddCommand(giveOrderToClientCmd)
	rootCmd.AddCommand(returnOrderDeliveryCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
