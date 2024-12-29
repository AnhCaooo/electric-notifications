package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AnhCaooo/electric-notifications/internal/constants"
	"github.com/AnhCaooo/electric-notifications/internal/db"
	"github.com/AnhCaooo/electric-notifications/internal/firebase"
	"github.com/AnhCaooo/electric-notifications/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

const (
	PUSH_NOTIFICATION_EXCHANGE string = "price_notifications"
	PUSH_NOTIFICATION_KEY      string = "price_notification_key"
	PUSH_NOTIFICATION_QUEUE    string = "push_notification_queue"
)

// Consumer represents a RabbitMQ consumer with necessary dependencies and configurations.
type Consumer struct {
	// The AMQP channel used for communication with RabbitMQ.
	channel *amqp.Channel
	// The AMQP connection to the RabbitMQ server.
	connection *amqp.Connection
	// The context for managing the consumer's lifecycle and cancellation.
	ctx context.Context
	// The name of the RabbitMQ exchange to bind the consumer to.
	exchange string
	//  The logger instance for logging consumer activities.
	logger *zap.Logger
	// The MongoDB instance for database operations.
	mongo *db.Mongo
	// The Firebase instance for sending push notifications.
	firebase *firebase.Firebase
	// The RabbitMQ queue to consume messages from.
	queue *amqp.Queue
	// The identifier for the worker handling the consumer.
	workerID int
}

// declareQueue declares a queue with the given name on the consumer's channel.
// It ensures the channel is not nil before attempting to declare the queue.
// If the queue declaration is successful, it assigns the declared queue to the consumer's queue field.
// Returns an error if the channel is nil or if the queue declaration fails.
func (c *Consumer) declareQueue(queueName string) error {
	if c.channel == nil {
		return fmt.Errorf("consumer channel is nil, ensure connection is established")
	}

	durable, autoDelete, exclusive, noWait := false, false, false, false
	queue, err := c.channel.QueueDeclare(
		queueName,  // queue name
		durable,    // durable
		autoDelete, // auto-delete when unused
		exclusive,  // exclusive
		noWait,     // no-wait
		nil,        // args
	)
	c.queue = &queue
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %s", err.Error())
	}
	return nil
}

// bindQueue binds the consumer's queue to the specified routing key on the exchange.
// It logs the binding action and returns an error if the binding fails.
func (c *Consumer) bindQueue(routingKey string) error {
	c.logger.Info(
		fmt.Sprintf("[worker_%d] binding queue to exchange with routing key", c.workerID),
		zap.String("queue_name", c.queue.Name),
		zap.String("exchange", c.exchange),
		zap.String("routing_key", routingKey),
	)

	if err := c.channel.QueueBind(
		c.queue.Name,
		routingKey,
		c.exchange,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("failed to bind a queue: %s", err.Error())
	}
	return nil
}

// Listen will start to read messages from the queue
func (c *Consumer) Listen(stopChan <-chan struct{}, errChan chan<- error) {
	consumer, autoAck, exclusive, noLocal, noWait := "", false, false, false, false
	msgs, err := c.channel.Consume(
		c.queue.Name, // queue
		consumer,     // consumer
		autoAck,      // auto-ack
		exclusive,    // exclusive
		noLocal,      // no-local
		noWait,       // no-wait
		nil,          // args
	)
	if err != nil {
		errMessage := fmt.Errorf("[worker_%d] failed to register a consumer: %s", c.workerID, err.Error())
		errChan <- errMessage
		return
	}

	c.logger.Info(fmt.Sprintf("[worker_%d] waiting for messages from %s...", c.workerID, c.queue.Name))

	// Make a channel to receive messages into infinite loop.
	for {
		select {
		case <-stopChan: // Respond to shutdown signal
			c.logger.Info(fmt.Sprintf("[worker_%d] stop listening for messages from %s...", c.workerID, c.queue.Name))
			return
		case msg, ok := <-msgs:
			if !ok {
				c.logger.Info(fmt.Sprintf("[worker_%d] message channel closed", c.workerID))
				return
			}

			if err := msg.Ack(false); err != nil {
				errMsg := fmt.Errorf("[worker_%d] error acknowledging message from %s: %s", c.workerID, c.queue.Name, err.Error())
				errChan <- errMsg
				return
			}

			// Process message
			switch msg.RoutingKey {
			case PUSH_NOTIFICATION_KEY:
				c.logger.Info(fmt.Sprintf("[worker_%d] received a message for pushing notification", c.workerID))
				var notificationMessage models.NotificationMessage
				json.Unmarshal(msg.Body, &notificationMessage)

				// retrieve all associated device tokens with given userId
				tokens, err := c.mongo.GetTokens(notificationMessage.UserId)
				if err != nil {
					errMsg := fmt.Errorf("[worker_%d] %s failed to get tokens: %s", c.workerID, constants.Server, err.Error())
					errChan <- errMsg
					return
				}
				err = c.firebase.SendToMultiTokens(tokens, notificationMessage.UserId, notificationMessage.Message)
				if err != nil {
					errMsg := fmt.Errorf("[worker_%d] %s failed to send multi tokens: %s", c.workerID, constants.Server, err.Error())
					errChan <- errMsg
					return
				}
				c.logger.Info(fmt.Sprintf("[worker_%d] send tokens successfully", c.workerID))
			default:
				c.logger.Info(fmt.Sprintf("[worker_%d] received an message from undefined routing key: '%s' with message: %v", c.workerID, msg.RoutingKey, msg.Body))
			}

		}
	}
}
