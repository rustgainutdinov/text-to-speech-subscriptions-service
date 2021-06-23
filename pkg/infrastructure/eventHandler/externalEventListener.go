package eventHandler

import (
	"encoding/json"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"os"
	"subscriptions-service/pkg/domain"
)

const (
	rabbitMqMessageType = iota
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
			eventInfo := &rabbitMqMessage{}
			err := json.Unmarshal(d.Body, eventInfo)
			if err != nil {
				log.Errorf("Error decoding JSON: %s", err)
				return
			}
			switch eventInfo.Type {
			case rabbitMqMessageType:
				e.textTranslatedEventHandler(d.Body, d)
			default:
				log.Errorf("Unknown event type: %v", eventInfo.Type)
				return
			}
		}
	}()
}

func (e *ExternalEventListener) textTranslatedEventHandler(messageBody []byte, d amqp.Delivery) {
	textTranslatedInfo := &textTranslatedRabbitMqMessage{}
	err := json.Unmarshal(messageBody, textTranslatedInfo)
	if err != nil {
		log.Errorf("Error decoding JSON: %s", err)
		return
	}
	userID, err := uuid.Parse(textTranslatedInfo.Data.UserID)
	if err != nil {
		log.Errorf("Error decoding JSON: %s", err)
		return
	}
	if err := d.Ack(false); err != nil {
		log.Errorf("Error acknowledging message : %s", err)
		return
	} else {
		log.Info("Acknowledged message")
	}
	err = domain.NewBalanceService(e.balanceRepo).WriteOffFromBalance(userID, textTranslatedInfo.Data.Score)
	if err != nil {
		log.Errorf("Error decoding JSON: %s", err)
	}
}

func NewExternalEventListener(rabbitMqChannel *amqp.Channel, balanceRepo domain.BalanceRepo) (*ExternalEventListener, error) {
	queue, err := rabbitMqChannel.QueueDeclare("translationQueue", true, false, false, false, nil)
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

type rabbitMqMessage struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

type textTranslatedInfo struct {
	UserID string `json:"userID"`
	Score  int    `json:"score"`
}

type textTranslatedRabbitMqMessage struct {
	Type int                `json:"type"`
	Data textTranslatedInfo `json:"data"`
}
