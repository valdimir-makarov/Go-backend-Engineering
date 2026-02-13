# Chat Service & Messenger Client - Design Patterns & Refactoring Guide

This document outlines the architectural improvements and design patterns identified for the `chat-service` and `Messenger-Client` to enhance code quality, maintainability, and scalability.

## 1. Message Handling Strategy (Strategy Pattern)

**Current Issue:**
The `chat_handler.go` file currently uses separate methods (`handleIncomingMessage`, `handleGroupMessage`) and potentially `switch` statements to handle different message types. This violates the Open/Closed Principle as adding new message types requires modifying the core handler logic.

**Proposed Solution:**
Implement the **Strategy Pattern** to encapsulate message processing logic.

*   **Interface:** Define a `MessageHandler` interface.
    ```go
    type MessageHandler interface {
        Handle(msg models.Message, session *Session) error
    }
    ```
*   **Concrete Strategies:**
    *   `DirectMessageHandler`: Handles 1-on-1 messages.
    *   `GroupMessageHandler`: Handles group chat messages.
    *   `FileMessageHandler`: Handles file upload notifications.
*   **Context:** The `WebSocketHandler` will act as the context, selecting the appropriate strategy based on the message type (e.g., `msg.Type` or context).

**Benefits:**
*   Easily add new message types without modifying existing code.
*   Isolates logic for different message types.

## 2. Message Delivery System (Observer Pattern)

**Current Issue:**
The `WebSocketHandler` manually manages a map of user channels (`userschannel`) and pushes messages directly. This couples connection management with message routing.

**Proposed Solution:**
Formalize the **Observer Pattern**.

*   **Subject:** A `ConnectionManager` or `Hub` that maintains the list of active connections.
*   **Observer:** Each user's WebSocket connection (wrapped in a struct) implements an `Update(msg models.Message)` method.
*   **Flow:** When a message needs to be delivered, the `Hub` notifies the specific `Observer` (user connection).

**Benefits:**
*   Decouples the logic of *how* a message is sent (WebSocket write) from *when* it is sent.
*   Simplifies broadcasting (notify all observers) vs. unicasting.

## 3. Service Initialization (Builder Pattern)

**Current Issue:**
In `cmd/main.go`, the server setup involves manually initializing DBs, repositories, services, producers, and handlers in a specific order. This "main bloat" is error-prone and hard to test.

**Proposed Solution:**
Use the **Builder Pattern** for server construction.

*   **Builder:** `ServerBuilder` struct.
*   **Methods:**
    *   `.WithDatabase(config)`
    *   `.WithKafka(brokers)`
    *   `.WithChatService()`
    *   `.Build()`
*   **Usage:**
    ```go
    server := NewServerBuilder().
        WithDatabase(dbConfig).
        WithKafka(kafkaConfig).
        Build()
    server.Run()
    ```

**Benefits:**
*   Clean, readable `main` function.
*   Easy to inject mock dependencies for integration testing.

## 4. Repository Access (Abstract Factory Pattern)

**Current Issue:**
Repositories are instantiated directly (e.g., `repository.NewWebSocketRepo()`). This binds the application to a specific implementation (e.g., GORM/Postgres).

**Proposed Solution:**
Use the **Abstract Factory Pattern**.

*   **Interface:** `RepositoryFactory`.
*   **Concrete Factory:** `PostgresRepositoryFactory`.
*   **Usage:** The service accepts a `RepositoryFactory` and asks it for the `ChatRepository` or `UserRepository`.

**Benefits:**
*   Allows seamless switching of storage backends (e.g., to MongoDB for chat history).
*   Simplifies testing by providing a `MockRepositoryFactory`.

## 5. Cross-Cutting Concerns (Decorator Pattern)

**Current Issue:**
Logging, metrics, and caching logic are mixed with business logic inside services (e.g., `utils.Info` calls inside `SendMessages`).

**Proposed Solution:**
Use the **Decorator Pattern**.

*   **Component:** `ChatService` interface.
*   **Concrete Component:** `CoreChatService` (contains only business logic).
*   **Decorator:** `LoggingChatServiceWrapper`.
    ```go
    type LoggingChatServiceWrapper struct {
        next ChatService
    }
    func (l *LoggingChatServiceWrapper) SendMessages(...) {
        log.Info("Sending message...")
        l.next.SendMessages(...)
        log.Info("Message sent.")
    }
    ```

**Benefits:**
*   Keeps core business logic pure and testable.
*   "Plug-and-play" behavior for logging, caching, or tracing.

## 6. Kafka Event Creation (Factory Method)

**Current Issue:**
Kafka messages are manually constructed in `kafka_prod.go`. This can lead to inconsistent event schemas.

**Proposed Solution:**
Use a **Factory Method** for event creation.

*   **Factory:** `EventFactory`.
*   **Method:** `CreateMessageEvent(msg models.Message) kafka.Message`.

**Benefits:**
*   Enforces a consistent schema for all events.
*   Centralizes event versioning logic.

---

## Implementation Roadmap

1.  **Refactor `chat_handler.go`:** Implement the **Strategy Pattern** first to clean up the growing message handling logic.
2.  **Refactor `main.go`:** Implement the **Builder Pattern** to simplify startup and enable better testing.
3.  **Logging Wrapper:** Apply the **Decorator Pattern** to remove logging noise from the core `ChatService`.
