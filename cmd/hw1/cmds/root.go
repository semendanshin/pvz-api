package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"homework/internal/infrastructure/handlers/bubbletea"
	"homework/internal/infrastructure/repositories/pvzorder"
	"homework/internal/usecases"
	"homework/internal/usecases/packager"
	"homework/internal/usecases/packager/strategies"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "hw1",
	Short: "Homework 1",
	Run: func(cmd *cobra.Command, args []string) {
		ordersFile, _ := cmd.Flags().GetString("orders")
		pvzID, _ := cmd.Flags().GetString("pvz")

		pvzOrderUseCase := InitUseCase(ordersFile, pvzID)

		handler := bubbletea.NewHandler(pvzOrderUseCase)

		err := handler.Run()
		if err != nil {
			panic(err)
		}
	},
}

func InitUseCase(ordersFile string, pvzID string) abstractions.IPVZOrderUseCase {
	pvzOrderRepository := pvzorder.NewJSONRepository(ordersFile)

	orderPackager := packager.NewOrderPackager(
		map[domain.PackagingType]packager.OrderPackagerStrategy{
			domain.PackagingTypeBox:  strategies.NewBoxPackager(),
			domain.PackagingTypeFilm: strategies.NewFilmPackager(),
			domain.PackagingTypeBag:  strategies.NewBagPackager(),
		},
	)

	return usecases.NewPVZOrderUseCase(
		pvzOrderRepository,
		orderPackager,
		pvzID,
	)

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

	rootCmd.SetOut(os.Stdout)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
