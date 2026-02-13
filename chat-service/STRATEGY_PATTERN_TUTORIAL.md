# Tutorial: Applying Strategy Pattern to Chat Message Handling

This tutorial explains how to refactor the current message handling logic in `chat_handler.go` using the **Strategy Pattern**. This pattern allows us to define a family of algorithms (message handling logics), encapsulate each one, and make them interchangeable.

---

### Step 1: Define the Strategy Interface

First, we define an interface that all message handling strategies must implement. This interface will have a `Handle` method.

```go
type MessageStrategy interface {
    Handle(h *WebSocketHandler, msg models.Message) error
}
```

---

### Step 2: Implement Concrete Strategies

Now, we create concrete implementations for different types of messages: Private and Group.

#### Private Message Strategy

This strategy handles individual messages between two users.

```go
type PrivateMessageStrategy struct{}

func (s *PrivateMessageStrategy) Handle(h *WebSocketHandler, msg models.Message) error {
    if msg.ReceiverID <= 0 {
        return fmt.Errorf("invalid receiver ID: %d", msg.ReceiverID)
    }

    if msg.ID == uuid.Nil {
        msg.ID = uuid.New()
    }

    // Save to DB
    if h.handler != nil {
        h.handler.SendMessages(msg.SenderID, msg.ReceiverID, msg.Content, msg.ID)
    }

    // Publish to Kafka
    if h.kafkaProducer != nil {
        h.kafkaProducer.PublishMessage(msg)
    }

    // Send to specific user channel
    lock.RLock()
    ch, ok := userschannel[msg.ReceiverID]
    lock.RUnlock()

    if ok && ch != nil {
        select {
        case ch <- msg:
            utils.Info("Private message queued", zap.Int("receiver_id", msg.ReceiverID))
        default:
            utils.Info("Queue full for receiver", zap.Int("receiver_id", msg.ReceiverID))
        }
    }
    return nil
}
```

#### Group Message Strategy

This strategy handles messages sent to a group.

```go
type GroupMessageStrategy struct{}

func (s *GroupMessageStrategy) Handle(h *WebSocketHandler, msg models.Message) error {
    if msg.GroupID == nil {
        return fmt.Errorf("GroupID missing in group message")
    }

    memberIDs, err := h.handler.GetGroupMemberIDs(*msg.GroupID)
    if err != nil {
        return err
    }

    for _, uid := range memberIDs {
        if uid == msg.SenderID {
            continue
        }
        lock.RLock()
        ch, ok := userschannel[uid]
        lock.RUnlock()
        if ok {
            ch <- msg
        }
    }
    return nil
}
```

---

### Step 3: Integrate with WebSocketHandler

Add a `strategies` map to your `WebSocketHandler` struct and initialize it.

```go
type WebSocketHandler struct {
    handler           *service.Service
    kafkaProducer     *kafka.KafkaProducer
    kafkaProducerAuth *authkafka.KafkaProducer
    strategies        map[string]MessageStrategy // New field
}

func NewWebSocketHandler(svc *service.Service, producer *kafka.KafkaProducer, authProducer *authkafka.KafkaProducer) *WebSocketHandler {
    h := &WebSocketHandler{
        handler:           svc,
        kafkaProducer:     producer,
        kafkaProducerAuth: authProducer,
        strategies:        make(map[string]MessageStrategy),
    }

    // Register strategies
    h.strategies["private"] = &PrivateMessageStrategy{}
    h.strategies["group"]   = &GroupMessageStrategy{}

    return h
}
```

---

### Step 4: Use the Strategy

Finally, refactor the message listener to use the registered strategies.

```go
func (h *WebSocketHandler) processMessage(msg models.Message) {
    var strategyKey string
    if msg.GroupID != nil {
        strategyKey = "group"
    } else {
        strategyKey = "private"
    }

    if strategy, ok := h.strategies[strategyKey]; ok {
        err := strategy.Handle(h, msg)
        if err != nil {
            utils.Error("Strategy execution failed", zap.Error(err))
        }
    } else {
        utils.Error("No strategy found for message type")
    }
}
```

---

### Benefits of this Approach

1.  **Open/Closed Principle**: You can add new message types (e.g., "broadcast", "system", "image") by creating new strategy classes without modifying the core `WebSocketHandler` logic.
2.  **Cleaner Code**: The `WebSocketHandler` doesn't need giant `if/else` or `switch` blocks to decide how to process each message.
3.  **Testability**: Individual strategies can be tested in isolation.
