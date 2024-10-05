package cmds

import (
	"context"
	"github.com/spf13/cobra"
	"homework/internal/abstractions"
	"homework/internal/infrastructure/handlers/bubbletea"
)

func rootCMD(pvzOrderUseCase abstractions.IPVZOrderUseCase) *cobra.Command {
	command := &cobra.Command{
		Use:   "hw1",
		Short: "Homework 1",
		Run: func(cmd *cobra.Command, args []string) {
			handler := bubbletea.NewHandler(pvzOrderUseCase)

			err := handler.Run(cmd.Context())
			if err != nil {
				panic(err)
			}
		},
	}

	return command
}

func setup(ctx context.Context, pvzOrderUseCase abstractions.IPVZOrderUseCase) *cobra.Command {
	rootCmd := rootCMD(pvzOrderUseCase)
	rootCmd.SetContext(ctx)

	rootCmd.AddCommand(acceptDeliveryCmd(pvzOrderUseCase))
	rootCmd.AddCommand(acceptReturnCmd(pvzOrderUseCase))
	rootCmd.AddCommand(getOrdersCmd(pvzOrderUseCase))
	rootCmd.AddCommand(getReturnsCmd(pvzOrderUseCase))
	rootCmd.AddCommand(giveOrderToClientCmd(pvzOrderUseCase))
	rootCmd.AddCommand(returnOrderDeliveryCmd(pvzOrderUseCase))

	return rootCmd
}

// Execute executes the root command.
func Execute(ctx context.Context, pvzOrderUseCase abstractions.IPVZOrderUseCase) error {
	return setup(ctx, pvzOrderUseCase).Execute()
}
