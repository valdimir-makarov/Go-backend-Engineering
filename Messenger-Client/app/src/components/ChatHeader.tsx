import React from 'react';

interface ChatHeaderProps {
    chatName: string;
    status?: string;
    avatar?: string;
}

const ChatHeader: React.FC<ChatHeaderProps> = ({ chatName, status = 'Online', avatar }) => {
    return (
        <div className="h-16 border-b border-gray-200 bg-white flex items-center px-6 justify-between shadow-sm z-10">
            <div className="flex items-center">
                {avatar && (
                    <img src={avatar} alt={chatName} className="w-10 h-10 rounded-full object-cover mr-3" />
                )}
                <div>
                    <h2 className="font-bold text-gray-800 text-lg">{chatName}</h2>
                    <div className="flex items-center">
                        <span className="w-2 h-2 bg-green-500 rounded-full mr-1.5"></span>
                        <span className="text-xs text-gray-500">{status}</span>
                    </div>
                </div>
            </div>

            <div className="flex items-center space-x-4 text-gray-500">
                <button className="hover:text-blue-600 transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                    </svg>
                </button>
                <button className="hover:text-blue-600 transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                </button>
                <button className="hover:text-blue-600 transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                </button>
            </div>
        </div>
    );
};

export default ChatHeader;
