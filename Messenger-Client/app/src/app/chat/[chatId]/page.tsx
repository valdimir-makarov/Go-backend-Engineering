'use client';
import React, { useEffect } from 'react';
import { useParams, useSearchParams } from 'next/navigation';
import { useState } from 'react';
import MessageBubble from '@/components/MessageBubble';
import MessageInput from '@/components/MessageInput';
import ChatSidebar from '@/components/ChatSidebar';
import ChatHeader from '@/components/ChatHeader';
import ConversationSummary from '@/components/ConversationSummary';
import { useAppDispatch, useAppSelector } from '@/lib/hooks';
import { sendMessage, setActiveChat } from '@/lib/features/chat/chatSlice';

const ChatPage = () => {
  const params = useParams();
  const searchParams = useSearchParams();
  const chatId = params?.chatId as string;
  const dispatch = useAppDispatch();
  const [isSummaryOpen, setIsSummaryOpen] = useState(false);

  // Get current user ID from query param or default to 2
  const currentUserId = parseInt(searchParams.get('uid') || '2');

  // Select messages for the current chat from Redux store
  const messages = useAppSelector((state) => state.chat.messages[chatId] || []);

  useEffect(() => {
    if (chatId) {
      dispatch(setActiveChat(chatId));
    }

    // Connect to WebSocket with dynamic user ID
    dispatch({
      type: 'chat/connect',
      payload: {
        userId: currentUserId,
        token: 'mock-token'
      }
    });

    // Also update current user in store so sendMessage knows who is sending
    dispatch({
      type: 'chat/setCurrentUser',
      payload: { id: currentUserId, name: `User ${currentUserId}` }
    });

    return () => {
      // Optional: disconnect on unmount if desired, but usually we keep it open
      // dispatch({ type: 'chat/disconnect' });
    };
  }, [chatId, dispatch, currentUserId]);

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

      <div className="flex-1 flex flex-col h-full relative min-w-0">
        <ChatHeader
          chatName={chatId === '1' ? 'User 1' : chatId === 'alice' ? 'Alice' : chatId === 'general' ? 'General' : 'Chat'}
          avatar={`https://ui-avatars.com/api/?name=${chatId}&background=random`}
          onToggleSummary={() => setIsSummaryOpen(!isSummaryOpen)}
        />

        <div className="flex-1 flex overflow-hidden">
          <div className="flex-1 flex flex-col relative">
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

          <ConversationSummary isOpen={isSummaryOpen} onClose={() => setIsSummaryOpen(false)} />
        </div>
      </div>
    </div>
  );
};

export default ChatPage;
