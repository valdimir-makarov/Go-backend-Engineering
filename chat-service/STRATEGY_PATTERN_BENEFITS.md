# Why Use the Strategy Pattern?

The Strategy Pattern is about **extensibility** and **clean separation of concerns**. Here's why it's better than a standard "If-Else" approach.

### 1. The "Before" (If-Else Hell)

Without the strategy pattern, your message handler would look like this:

```go
func (h *WebSocketHandler) HandleIncomingMessage(msg models.Message) {
    if msg.GroupID != nil {
        // --- 20 lines of Group Logic ---
        memberIDs, _ := h.handler.GetGroupMemberIDs(*msg.GroupID)
        for _, uid := range memberIDs {
            // ... logic to send message to group members ...
        }
    } else if msg.IsImage {
         // --- 15 lines of Image processing logic ---
    } else if msg.IsSystemMessage {
         // --- 10 lines of System message logic ---
    } else {
        // --- 20 lines of Private Chat Logic ---
        h.handler.SendMessages(msg.SenderID, msg.ReceiverID, msg.Content, msg.ID)
        // ... kafka publishing ...
        // ... socket sending ...
    }
}
```

**Problems with this:**

- **Violates Open/Closed Principle**: Every time you add a new feature (like "Voice Messages" or "File Sharing"), you have to modify this core function.
- **Hard to Read**: As your chat app grows, this function becomes a "God Function" with hundreds of lines.
- **Tightly Coupled**: All the logic for Kafka, SQL, and WebSockets is mixed together for every message type.

---

### 2. The "After" (Strategy Pattern)

With the strategy pattern, the core handler only does **one thing**: Routing.

```go
func (h *WebSocketHandler) HandleIncomingMessage(msg models.Message) {
    // Find the right strategy based on message type
    strategy := h.getStrategy(msg)

    // Execute it
    strategy.Handle(h, msg)
}
```

**Benefits:**

1.  **Open for Extension**: Want to add "Voice Messages"? Just create a `VoiceMessageStrategy` struct. You **never** touch the `HandleIncomingMessage` function again.
2.  **Isolated Testing**: You can test `PrivateMessageStrategy` completely separately from `GroupMessageStrategy`.
3.  **Clean Code**: Each strategy file only contains logic relevant to that specific task.

---

### Example: Adding a "System Notify" Feature

If tomorrow you want to send "User X joined the group" messages:

**Without Strategy:** You hunt through `chat_handler.go`, find the `if` block, add another `else if`, and risk breaking private chat.

**With Strategy:**

1.  Create `SystemNotifyStrategy.go`.
2.  Implement `Handle(...)`.
3.  Register it in `NewWebSocketHandler`.
    **Result:** Zero risk of breaking existing features.
