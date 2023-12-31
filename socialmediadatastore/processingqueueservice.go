package socialmediadatastore

import (
	ctx "context"
	"encoding/json"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus"
)

var queue_client *azservicebus.Client

const (
	NEW          = "0_new"
	INTERESTING  = "1_interesting"
	SHORT_LISTED = "2_short_listed"
	USER_ACTION  = "3_action_suggested"
	ACTION_TAKEN = "4_action_taken"
	IGNORE       = "9_ignore"
)

func getQueueName(item_type string) string {
	switch item_type {
	case NEW:
		return getNewItemsQueue()
	case INTERESTING:
		return getInterestingItemsQueue()
	case SHORT_LISTED:
		return getShortListedItemsQueue()
	case USER_ACTION:
		return getUserActionsQueue()
	default:
		return ""
	}
}

func createSbMsg(item *UserActionData) azservicebus.Message {
	body, _ := json.Marshal(item) //standard json blob, no error expected here
	msg := azservicebus.Message{Body: body}
	return msg
}

func InitializeProcessingQueues() *azservicebus.Client {
	if queue_client != nil {
		// since connection already exists there is no need to create new connection
		return queue_client
	}

	client, err := azservicebus.NewClientFromConnectionString(getServiceBusConnection(), nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	queue_client = client
	return client
}

func Que(process_stage string, data *UserActionData) bool {
	sender, _ := queue_client.NewSender(getQueueName(process_stage), nil)
	defer sender.Close(ctx.Background())

	msg := createSbMsg(data)
	if err := sender.SendMessage(ctx.Background(), &msg, nil); err != nil {
		//just log sending the message. precision is not a target here
		log.Println("Failed sending message: ", err)
		return false
	}
	return true
}

// adds an entire array in the que
func BatchQue(process_stage string, items []UserActionData) {
	if len(items) > 0 {

		sender, _ := queue_client.NewSender(getQueueName(process_stage), nil)
		defer sender.Close(ctx.Background())

		msg_batch, _ := sender.NewMessageBatch(ctx.Background(), nil)
		for i := range items {
			msg := createSbMsg(&items[i])
			msg_batch.AddMessage(&msg, nil)
		}

		if err := sender.SendMessageBatch(ctx.Background(), msg_batch, nil); err != nil {
			//just log sending the message. precision is not a target here
			log.Println("Failed sending message batch: ", err)
		}
	}
}

// returns an array of items in the queue based on the process stage
// if there is no item in the queue it will block until there is at least one item
// if there are items in the queue then it will return at most MAX_BATCH_SIZE number of items at a time
// if the number of item <= MAX_BATCH_SIZE it will return all the items in the queue
func Deque(process_stage string) []UserActionData {

	rcvr_ctx, cancel := ctx.WithTimeout(ctx.Background(), getMaxWaitTime())
	defer cancel()

	rcvr, _ := queue_client.NewReceiverForQueue(getQueueName(process_stage), nil)
	defer rcvr.Close(rcvr_ctx)

	// lesson learned: you dont need a timeout context if you are willing to wait for at least one message
	// its doesnt matter if the current queue has less that MAX_BATCH_SIZE, the queue will return however many items there are
	// as long as there is at least 1 item and will cap the return to MAX_BATCH_SIZE
	messages, _ := rcvr.ReceiveMessages(rcvr_ctx, getMaxBatchSize(), nil)
	resp := make([]UserActionData, len(messages))
	for i, msg := range messages {
		json.Unmarshal(msg.Body, &resp[i])
		rcvr.CompleteMessage(rcvr_ctx, msg, nil)
	}

	return resp
}
