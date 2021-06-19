package infrastructure

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"log"
	"os"
	"subscriptions-service/pkg/domain"
)

type ExternalEventListener struct {
	messageChannel <-chan amqp.Delivery
	balanceRepo    domain.BalanceRepo
}

func (e *ExternalEventListener) activateTextTranslatedHandler() {
	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range e.messageChannel {
			log.Printf("Received a message: %s", d.Body)
			eventInfo := &textTranslatedInfo{}
			err := json.Unmarshal(d.Body, eventInfo)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
			userID, err := uuid.Parse(eventInfo.UserID)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
			err = domain.NewBalanceService(e.balanceRepo).WriteOffFromBalance(userID, eventInfo.AmountOfSymbols)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}
		}
	}()
}

func NewExternalEventListener(rabbitMqChannel *amqp.Channel, balanceRepo domain.BalanceRepo) (*ExternalEventListener, error) {
	queue, err := rabbitMqChannel.QueueDeclare("textTranslated", true, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	err = rabbitMqChannel.Qos(1, 0, false)
	if err != nil {
		return nil, err
	}
	messageChannel, err := rabbitMqChannel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	return &ExternalEventListener{messageChannel: messageChannel, balanceRepo: balanceRepo}, nil
}

type textTranslatedInfo struct {
	UserID          string `json:"userID"`
	AmountOfSymbols int    `json:"amountOfSymbols"`
}
