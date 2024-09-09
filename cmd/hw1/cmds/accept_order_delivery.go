package cmds

import (
	"github.com/spf13/cobra"
	"homework/internal/infrastructure/repositories/pvzorder"
	"homework/internal/usecases"
	"time"
)

var acceptDeliveryCmd = &cobra.Command{
	Use:     "accept_delivery",
	Short:   "Accept delivery",
	Args:    cobra.ExactArgs(3),
	Example: "hw1 accept_delivery <order_id> <recipient_id> <storage_time: 1h30m>",
	RunE: func(cmd *cobra.Command, args []string) error {
		pvzOrderRepository := pvzorder.NewJSONRepository(cmd.Flag("orders").Value.String())

		pvzOrderUseCase := usecases.NewPVZOrderUseCase(pvzOrderRepository, cmd.Flag("pvz").Value.String())

		orderID := args[0]
		recipientID := args[1]
		storageTime, err := time.ParseDuration(args[2])
		if err != nil {
			return err
		}

		err = pvzOrderUseCase.AcceptOrderDelivery(orderID, recipientID, storageTime)
		if err != nil {
			return err
		}

		cmd.Println("Delivery accepted")

		return nil
	},
}
