import Link from "next/link";
import MessageBubble from "@/components/MessageBubble";

export default function Home() {
  return (
    <div className="p-8 font-sans">
      <h1 className="text-2xl font-bold mb-6">Messenger App</h1>

      <div className="mb-8">
        <h2 className="text-xl mb-4">Component Preview:</h2>
        <MessageBubble message="Hello! This is a preview bubble." sender="me" timestamp="10:00 AM" />
      </div>

      <div>
        <h2 className="text-xl mb-4">Active Chats:</h2>
        <ul className="list-disc pl-5 space-y-2">
          <li>
            <Link href="/chat/1" className="text-blue-500 hover:underline">
              Chat with User 1
            </Link>
          </li>
          <li>
            <Link href="/chat/alice" className="text-blue-500 hover:underline">
              Chat with Alice
            </Link>
          </li>
          <li>
            <Link href="/chat/general" className="text-blue-500 hover:underline">
              General Channel
            </Link>
          </li>
        </ul>
      </div>
    </div>
  );
}