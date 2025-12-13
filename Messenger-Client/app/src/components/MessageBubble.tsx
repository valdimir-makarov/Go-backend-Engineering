import React from 'react';

interface MessageBubbleProps {
  message: string;
  sender: string;
  timestamp: string;
}

const MessageBubble: React.FC<MessageBubbleProps> = ({ message, sender, timestamp }) => {
  const isMe = sender === 'me';
  return (
    <div className={`flex ${isMe ? 'justify-end' : 'justify-start'} mb-4`}>
      {!isMe && (
        <div className="w-8 h-8 rounded-full bg-gray-300 flex-shrink-0 mr-2 overflow-hidden">
          <img src={`https://ui-avatars.com/api/?name=${sender}&background=random`} alt={sender} className="w-full h-full object-cover" />
        </div>
      )}
      <div className={`max-w-[70%] px-5 py-3 rounded-2xl shadow-sm ${isMe
          ? 'bg-blue-600 text-white rounded-br-none'
          : 'bg-white text-gray-800 border border-gray-100 rounded-bl-none'
        }`}>
        <p className="text-sm leading-relaxed">{message}</p>
        <span className={`block text-[10px] text-right mt-1 ${isMe ? 'text-blue-200' : 'text-gray-400'}`}>
          {timestamp}
        </span>
      </div>
    </div>
  );
};

export default MessageBubble;