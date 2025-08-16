import express from "express";
import dotenv from "dotenv";
import { WebSocketServer } from "ws";
import { WebSocket } from "ws";
dotenv.config();
const app = express();
const PORT = process.env.PORT || 3000;
// Middleware
app.use(express.json());
const wss = new WebSocketServer({ port: 3005 });
wss.on("connection", (socket) => {
    socket.on(("message"), (msg) => {
        console.log("Received message:", msg);
        const message = JSON.parse(msg);
        socket.send(JSON.stringify(message));
    });
    socket.on("close", () => {
        console.log("Client disconnected");
    });
});
// Routes
app.get("/", (req, res) => {
    res.send("Hello from Express + TypeScript!");
});
// Start server
app.listen(PORT, () => {
    console.log(`🚀 Server running at http://localhost:${PORT}`);
});
//# sourceMappingURL=server.js.map