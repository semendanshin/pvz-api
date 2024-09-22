package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/abstractions"
)

func getReturnsCmd(pvzOrderUseCase abstractions.IPVZOrderUseCase) *cobra.Command {
	command := &cobra.Command{
		Use:     "get_returns",
		Short:   "Get returns",
		Args:    cobra.NoArgs,
		Example: "hw1 get_returns",
		RunE: func(cmd *cobra.Command, args []string) error {
			page, _ := cmd.Flags().GetInt("page")
			pageSize, _ := cmd.Flags().GetInt("pageSize")

			data, err := pvzOrderUseCase.GetReturns(
				abstractions.WithPage(page),
				abstractions.WithPageSize(pageSize),
			)
			if err != nil {
				return err
			}

			cmd.Println("Returns:")
			for _, order := range data {
				cmd.Println(order)
			}

			return nil
		},
	}

	command.Flags().Int("page", 0, "page")
	command.Flags().Int("pageSize", 10, "page size")

	return command
}
