package cmds

import (
	"fmt"
	"github.com/spf13/cobra"
	"homework/internal/domain"
	"strconv"
	"time"
)

var acceptDeliveryCmd = &cobra.Command{
	Use:     "accept_delivery",
	Short:   "Accept delivery",
	Args:    cobra.ExactArgs(6),
	Example: "hw1 accept_delivery <order_id> <recipient_id> <storage_time: 1h30m> <cost> <weight> <packaging>",
	RunE: func(cmd *cobra.Command, args []string) error {
		ordersFile, _ := cmd.Flags().GetString("orders")
		pvzID, _ := cmd.Flags().GetString("pvz")

		pvzOrderUseCase := InitUseCase(ordersFile, pvzID)

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

func init() {
	//acceptDeliveryCmd.Flags().Int("cost", 0, "cost")
	//acceptDeliveryCmd.Flags().Int("weight", 0, "weight")
	//acceptDeliveryCmd.Flags().String("packaging", "", "packaging")
	acceptDeliveryCmd.Flags().Bool("additional_film", false, "additional film")
}
