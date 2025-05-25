'use client';
import React, { useState } from 'react';
import MessageBubble from '@/components/MessageBubble';
import MessageInput from '@/components/MessageInput';

const ChatPage = () => {
  const [messages, setMessages] = useState([
    { id: 1, message: 'Hello!', sender: 'other', timestamp: '10:00 AM' },
    { id: 2, message: 'Hey there!', sender: 'me', timestamp: '10:01 AM' },
  ]);

  const handleSend = (text) => {
    const newMessage = {
      id: Date.now(),
      message: text,
      sender: 'me',
      timestamp: new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }),
    };
    setMessages((prev) => [...prev, newMessage]);
  };

  return (
    <div className="flex flex-col h-screen">
      <div className="flex-grow p-4 overflow-y-auto bg-gray-100">
        {messages.map((msg) => (
          <MessageBubble
            key={msg.id}
            message={msg.message}
            sender={msg.sender}
            timestamp={msg.timestamp}
          />
        ))}
      </div>
      <MessageInput onSend={handleSend} />
    </div>
  );
};

export default ChatPage;
