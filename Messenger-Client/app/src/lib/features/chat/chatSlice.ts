import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export interface Message {
    id: number;
    message: string;
    sender: string;
    timestamp: string;
}

export interface ChatState {
    messages: Record<string, Message[]>;
    activeChatId: string | null;
}

const initialState: ChatState = {
    messages: {
        '1': [
            { id: 1, message: 'Hello!', sender: 'other', timestamp: '10:00 AM' },
            { id: 2, message: 'Hey there!', sender: 'me', timestamp: '10:01 AM' },
        ],
        'alice': [
            { id: 1, message: 'Hi Alice!', sender: 'me', timestamp: '09:00 AM' },
            { id: 2, message: 'Can we meet?', sender: 'other', timestamp: '09:15 AM' },
        ],
        'general': [
            { id: 1, message: 'Welcome everyone!', sender: 'other', timestamp: 'Yesterday' },
        ]
    },
    activeChatId: null,
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
            };

            state.messages[chatId].push(newMessage);
        },
    },
});

export const { setActiveChat, sendMessage } = chatSlice.actions;
export default chatSlice.reducer;
