import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Message {
    id: number | string;
    message: string;
    sender: string;
    timestamp: string;
    senderId?: number;
    receiverId?: number;
}

export interface ChatState {
    messages: Record<string, Message[]>;
    activeChatId: string | null;
    isConnected: boolean;
    currentUser: { id: number; name: string } | null; // Mock current user
}

const initialState: ChatState = {
    messages: {
        '1': [
            { id: 1, message: 'Hello!', sender: 'other', timestamp: '10:00 AM', senderId: 1, receiverId: 2 },
            { id: 2, message: 'Hey there!', sender: 'me', timestamp: '10:01 AM', senderId: 2, receiverId: 1 },
        ],
        'alice': [
            { id: 1, message: 'Hi Alice!', sender: 'me', timestamp: '09:00 AM', senderId: 2, receiverId: 3 }, // Assuming Alice is 3
            { id: 2, message: 'Can we meet?', sender: 'other', timestamp: '09:15 AM', senderId: 3, receiverId: 2 },
        ],
        'general': [
            { id: 1, message: 'Welcome everyone!', sender: 'other', timestamp: 'Yesterday', senderId: 0, receiverId: 0 },
        ]
    },
    activeChatId: null,
    isConnected: false,
    currentUser: { id: 2, name: 'Me' }, // Hardcoded for now, assume we are user 2
};

export const chatSlice = createSlice({
    name: 'chat',
    initialState,
    reducers: {
        setActiveChat: (state, action: PayloadAction<string>) => {
            state.activeChatId = action.payload;
        },
        sendMessage: (state, action: PayloadAction<{ chatId: string; message: string }>) => {
            const { chatId, message } = action.payload;
            if (!state.messages[chatId]) {
                state.messages[chatId] = [];
            }

            const newMessage: Message = {
                id: Date.now(),
                message,
                sender: 'me',
                timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
                senderId: state.currentUser?.id,
                receiverId: parseInt(chatId) || 0, // Simple parsing for demo
            };

            state.messages[chatId].push(newMessage);
        },
        receiveMessage: (state, action: PayloadAction<Message>) => {
            const msg = action.payload;
            // Determine chat ID based on sender. If I am the receiver, the chat is with the sender.
            // If it's a group message (not handled deeply here), logic might differ.
            const chatId = msg.senderId?.toString() || 'general';

            if (!state.messages[chatId]) {
                state.messages[chatId] = [];
            }
            state.messages[chatId].push(msg);
        },
        setConnectionStatus: (state, action: PayloadAction<boolean>) => {
            state.isConnected = action.payload;
        },
        setCurrentUser: (state, action: PayloadAction<{ id: number; name: string }>) => {
            state.currentUser = action.payload;
        }
    },
});

export const { setActiveChat, sendMessage, receiveMessage, setConnectionStatus, setCurrentUser } = chatSlice.actions;
export default chatSlice.reducer;
