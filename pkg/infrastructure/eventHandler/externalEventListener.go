package eventHandler

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"subscriptions-service/pkg/domain"
)

type ExternalEventListener struct {
	messageChannel <-chan amqp.Delivery
	balanceRepo    domain.BalanceRepo
}

func (e *ExternalEventListener) ActivateExternalEventListener() {
	go func() {
		log.Infof("Consumer ready, PID: %d", os.Getpid())
		for d := range e.messageChannel {
			log.Infof("Received a message: %s", d.Body)
			eventInfo := &textTranslatedInfo{}
			err := json.Unmarshal(d.Body, eventInfo)
			if err != nil {
				log.Errorf("Error decoding JSON: %s", err)
			}
			userID, err := uuid.Parse(eventInfo.UserID)
			if err != nil {
				log.Errorf("Error decoding JSON: %s", err)
			}
			if err := d.Ack(false); err != nil {
				log.Errorf("Error acknowledging message : %s", err)
			} else {
				log.Info("Acknowledged message")
			}
			err = domain.NewBalanceService(e.balanceRepo).WriteOffFromBalance(userID, eventInfo.Score)
			if err != nil {
				log.Errorf("Error decoding JSON: %s", err)
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
	UserID string `json:"userID"`
	Score  int    `json:"score"`
}
