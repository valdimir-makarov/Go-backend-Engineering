import React from 'react';
import Link from 'next/link';

interface ChatSidebarProps {
    activeChatId?: string;
}

const chats = [
    { id: '1', name: 'User 1', lastMessage: 'See you later!', time: '10:30 AM', avatar: 'https://i.pravatar.cc/150?u=1' },
    { id: 'alice', name: 'Alice', lastMessage: 'Can we meet?', time: '09:15 AM', avatar: 'https://i.pravatar.cc/150?u=alice' },
    { id: 'general', name: 'General', lastMessage: 'Welcome everyone!', time: 'Yesterday', avatar: 'https://ui-avatars.com/api/?name=General&background=random' },
    { id: 'bob', name: 'Bob', lastMessage: 'Thanks for the help.', time: 'Mon', avatar: 'https://i.pravatar.cc/150?u=bob' },
];

const ChatSidebar: React.FC<ChatSidebarProps> = ({ activeChatId }) => {
    return (
        <div className="w-80 bg-white border-r border-gray-200 flex flex-col h-full">
            <div className="p-4 border-b border-gray-200 flex justify-between items-center">
                <h1 className="text-xl font-bold text-gray-800">Messages</h1>
                <button className="p-2 text-gray-500 hover:bg-gray-100 rounded-full">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                </button>
            </div>

            <div className="p-3">
                <input
                    type="text"
                    placeholder="Search chats..."
                    className="w-full px-4 py-2 bg-gray-100 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                />
            </div>

            <div className="flex-1 overflow-y-auto">
                {chats.map((chat) => (
                    <Link href={`/chat/${chat.id}`} key={chat.id}>
                        <div className={`flex items-center p-3 hover:bg-gray-50 cursor-pointer transition-colors ${activeChatId === chat.id ? 'bg-blue-50' : ''}`}>
                            <img src={chat.avatar} alt={chat.name} className="w-12 h-12 rounded-full object-cover" />
                            <div className="ml-3 flex-1 overflow-hidden">
                                <div className="flex justify-between items-baseline">
                                    <h3 className="font-semibold text-gray-900 truncate">{chat.name}</h3>
                                    <span className="text-xs text-gray-500">{chat.time}</span>
                                </div>
                                <p className="text-sm text-gray-500 truncate">{chat.lastMessage}</p>
                            </div>
                        </div>
                    </Link>
                ))}
            </div>
        </div>
    );
};

export default ChatSidebar;
