import React from 'react';

interface ConversationSummaryProps {
    isOpen: boolean;
    onClose: () => void;
}

const ConversationSummary: React.FC<ConversationSummaryProps> = ({ isOpen, onClose }) => {
    if (!isOpen) return null;

    return (
        <div className="w-80 bg-white border-l border-gray-200 flex flex-col h-full shadow-xl transition-all duration-300 ease-in-out">
            <div className="p-4 border-b border-gray-200 flex justify-between items-center bg-gray-50">
                <h2 className="text-lg font-bold text-gray-800">Conversation Insights</h2>
                <button
                    onClick={onClose}
                    className="p-2 text-gray-500 hover:bg-gray-200 rounded-full transition-colors"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                        <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                    </svg>
                </button>
            </div>

            <div className="flex-1 overflow-y-auto p-4 space-y-6">
                {/* Tone Analysis Section */}
                <div className="bg-blue-50 p-4 rounded-xl border border-blue-100">
                    <div className="flex items-center mb-2">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-blue-600 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14.828 14.828a4 4 0 01-5.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <h3 className="font-semibold text-blue-900">Tone Analysis</h3>
                    </div>
                    <div className="flex items-center space-x-2 mb-3">
                        <span className="px-3 py-1 bg-white text-blue-700 text-xs font-medium rounded-full shadow-sm border border-blue-100">Friendly</span>
                        <span className="px-3 py-1 bg-white text-blue-700 text-xs font-medium rounded-full shadow-sm border border-blue-100">Casual</span>
                    </div>
                    <p className="text-sm text-blue-800 leading-relaxed">
                        The conversation is light-hearted and positive. Both parties seem engaged and cooperative.
                    </p>
                </div>

                {/* Main Points Section */}
                <div>
                    <div className="flex items-center mb-3">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-gray-700 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
                        </svg>
                        <h3 className="font-bold text-gray-800">Key Takeaways</h3>
                    </div>
                    <ul className="space-y-3">
                        <li className="flex items-start">
                            <span className="flex-shrink-0 h-1.5 w-1.5 rounded-full bg-indigo-500 mt-2 mr-2"></span>
                            <span className="text-sm text-gray-600">Discussed project timeline updates for Q4.</span>
                        </li>
                        <li className="flex items-start">
                            <span className="flex-shrink-0 h-1.5 w-1.5 rounded-full bg-indigo-500 mt-2 mr-2"></span>
                            <span className="text-sm text-gray-600">Agreed to meet next Tuesday at 10 AM.</span>
                        </li>
                        <li className="flex items-start">
                            <span className="flex-shrink-0 h-1.5 w-1.5 rounded-full bg-indigo-500 mt-2 mr-2"></span>
                            <span className="text-sm text-gray-600">Action item: Send the updated design assets.</span>
                        </li>
                    </ul>
                </div>

                {/* Action Button */}
                <div className="pt-4">
                    <button className="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2.5 px-4 rounded-lg shadow-md transition-all duration-200 flex items-center justify-center group">
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2 group-hover:animate-pulse" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                        </svg>
                        Refresh Insights
                    </button>
                    <p className="text-xs text-center text-gray-400 mt-2">
                        Powered by AI Analysis
                    </p>
                </div>
            </div>
        </div>
    );
};

export default ConversationSummary;
