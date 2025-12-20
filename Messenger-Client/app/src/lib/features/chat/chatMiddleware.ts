import { Middleware } from '@reduxjs/toolkit';
import { receiveMessage, setConnectionStatus, sendMessage } from '@/lib/features/chat/chatSlice';

let socket: WebSocket | null = null;

export const chatMiddleware: Middleware = (store) => (next) => (action: any) => {
    if (action.type === 'chat/connect') {
        const { userId, token } = action.payload;
        if (socket) {
            socket.close();
        }

        socket = new WebSocket(`ws://localhost:3001/ws?user_id=${userId}&token=${token}`);

        socket.onopen = () => {
            store.dispatch(setConnectionStatus(true));
            console.log('WebSocket connected');
        };

        socket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                // Map backend message format to frontend format
                const message = {
                    id: data.id,
                    message: data.content,
                    sender: 'other', // Logic to determine sender name needed
                    timestamp: new Date(data.created_at || Date.now()).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
                    senderId: data.sender_id,
                    receiverId: data.receiver_id,
                };
                store.dispatch(receiveMessage(message));
            } catch (error) {
                console.error('Error parsing WebSocket message:', error);
            }
        };

        socket.onclose = () => {
            store.dispatch(setConnectionStatus(false));
            console.log('WebSocket disconnected');
        };

        socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    } else if (action.type === 'chat/disconnect') {
        if (socket) {
            socket.close();
            socket = null;
        }
        store.dispatch(setConnectionStatus(false));
    } else if (sendMessage.match(action)) {
        // Intercept sendMessage action to send to WebSocket
        if (socket && socket.readyState === WebSocket.OPEN) {
            const { chatId, message } = action.payload;
            const state = store.getState() as any;
            const payload = {
                sender_id: state.chat.currentUser?.id,
                receiver_id: parseInt(chatId), // Assuming chatId is numeric user ID for 1-on-1
                content: message,
            };
            socket.send(JSON.stringify(payload));
        }
    }

    return next(action);
};
