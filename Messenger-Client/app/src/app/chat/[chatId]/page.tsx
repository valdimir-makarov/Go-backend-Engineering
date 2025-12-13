'use client';
import React, { useEffect } from 'react';
import { useParams } from 'next/navigation';
import MessageBubble from '@/components/MessageBubble';
import MessageInput from '@/components/MessageInput';
import ChatSidebar from '@/components/ChatSidebar';
import ChatHeader from '@/components/ChatHeader';
import { useAppDispatch, useAppSelector } from '@/lib/hooks';
import { sendMessage, setActiveChat } from '@/lib/features/chat/chatSlice';

const ChatPage = () => {
  const params = useParams();
  const chatId = params?.chatId as string;
  const dispatch = useAppDispatch();

  // Select messages for the current chat from Redux store
  const messages = useAppSelector((state) => state.chat.messages[chatId] || []);

  useEffect(() => {
    if (chatId) {
      dispatch(setActiveChat(chatId));
    }
  }, [chatId, dispatch]);

  const handleSend = (text: string) => {
    dispatch(sendMessage({ chatId, message: text }));
  };

  // Scroll to bottom on new message
  useEffect(() => {
    const chatContainer = document.getElementById('chat-container');
    if (chatContainer) {
      chatContainer.scrollTop = chatContainer.scrollHeight;
    }
  }, [messages]);

  return (
    <div className="flex h-screen bg-gray-50">
      <ChatSidebar activeChatId={chatId} />

      <div className="flex-1 flex flex-col h-full relative">
        <ChatHeader
          chatName={chatId === '1' ? 'User 1' : chatId === 'alice' ? 'Alice' : chatId === 'general' ? 'General' : 'Chat'}
          avatar={`https://ui-avatars.com/api/?name=${chatId}&background=random`}
        />

        <div
          id="chat-container"
          className="flex-grow p-6 overflow-y-auto bg-[#f0f2f5]"
          style={{ backgroundImage: 'url("https://user-images.githubusercontent.com/15075759/28719144-86dc0f70-73b1-11e7-911d-60d70fcded21.png")', backgroundBlendMode: 'soft-light' }}
        >
          {messages.length > 0 ? (
            messages.map((msg) => (
              <MessageBubble
                key={msg.id}
                message={msg.message}
                sender={msg.sender}
                timestamp={msg.timestamp}
              />
            ))
          ) : (
            <div className="flex justify-center items-center h-full text-gray-400">
              <p>No messages yet. Start the conversation!</p>
            </div>
          )}
        </div>

        <MessageInput onSend={handleSend} />
      </div>
    </div>
  );
};

export default ChatPage;
