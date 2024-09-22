package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"homework/internal/abstractions"
	"homework/internal/domain"
	"strconv"
	"time"
)

func acceptDeliveryCmd(pvzOrderUseCase abstractions.IPVZOrderUseCase) *cobra.Command {
	command := &cobra.Command{
		Use:     "accept_delivery",
		Short:   "Accept delivery",
		Args:    cobra.ExactArgs(6),
		Example: "hw1 accept_delivery <order_id> <recipient_id> <storage_time: 1h30m> <cost> <weight> <packaging>",
		RunE: func(cmd *cobra.Command, args []string) error {
			orderID := args[0]

			recipientID := args[1]

			storageTime, err := time.ParseDuration(args[2])
			if err != nil {
				return err
			}

			if storageTime < 0 {
				return fmt.Errorf("storageTime is negative")
			}

			cost, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			weight, err := strconv.Atoi(args[4])
			if err != nil {
				return err
			}

			packaging, err := domain.NewPackagingType(args[5])
			if err != nil {
				return err
			}

			additionalFilm, _ := cmd.Flags().GetBool("additional_film")

			err = pvzOrderUseCase.AcceptOrderDelivery(
				cmd.Context(),
				orderID,
				recipientID,
				storageTime,
				cost,
				weight,
				packaging,
				additionalFilm,
			)
			if err != nil {
				return err
			}

			cmd.Println("Delivery accepted")

			return nil
		},
	}

	command.Flags().Bool("additional_film", false, "additional film")

	return command
}
