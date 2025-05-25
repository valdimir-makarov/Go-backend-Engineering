import React from 'react';

const MessageBubble = ({ message, sender, timestamp }) => {
  const isMe = sender === 'me';
  return (
    <div className={`flex ${isMe ? 'justify-end' : 'justify-start'} mb-2`}>
      <div className={`max-w-xs px-4 py-2 rounded-2xl shadow text-white ${
        isMe ? 'bg-blue-500' : 'bg-gray-600'
      }`}>
        <p className="text-sm">{message}</p>
        <span className="block text-[10px] text-right opacity-70 mt-1">{timestamp}</span>
      </div>
    </div>
  );
};

export default MessageBubble;