import type { ChatMessage } from "../generated/dto/social";

const MessageList = ({ messages, loading, selfId }: {
  messages: ChatMessage[];
  loading: boolean;
  selfId: number;
}) => (
  <div className="message-list">
    {loading ? (
      Array.from({ length: 8 }).map((_, i) => (
        <div className="message-skeleton" key={i} />
      ))
    ) : (
      messages.map((msg, i) => (
        <div
          className={`message${msg.senderId === selfId ? " self" : ""}`}
          key={i}
        >
          {msg.content}
        </div>
      ))
    )}
  </div>
);

export default MessageList;

