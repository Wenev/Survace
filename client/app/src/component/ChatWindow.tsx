import { useEffect, useState, useRef } from "react";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import defaultAvatar from "../assets/default.jpg";
import type { ChatMessage } from "../generated/dto/social.ts";

interface UserProfile {
    id: number;
    username: string;
    avatarUrl: string;
}

export default function ChatWindow({ friendId }: { friendId: number | null }){
    const { user } = useAuth();
    const { social, galactus } = useGrpc();
    const [messages, setMessages] = useState<ChatMessage[]>([]);
    const [loading, setLoading] = useState(false);
    const [input, setInput] = useState("");
    const [sending, setSending] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [userProfiles, setUserProfiles] = useState<Record<number, UserProfile>>({});
    const [typingUserId, setTypingUserId] = useState<number | null>(null);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const typingTimeoutRef = useRef<NodeJS.Timeout | null>(null);

    useEffect(() => {
        let streamController: AbortController | null = null;
        let cancelled = false;
        let reconnectTimeout: ReturnType<typeof setTimeout> | null = null;

        const fetchInitialMessages = async () => {
            if (!user || !friendId) return;
            setLoading(true);
            setError(null);
            try {
                const res = await social.listMessages({ userId: user.userId, otherUserId: friendId });
                const msgs = res.response.data || [];
                setMessages(msgs);
                const userIds = Array.from(new Set(msgs.map((m: ChatMessage) => m.senderId)));
                const missingIds = userIds.filter(id => !(id in userProfiles));
                if (missingIds.length > 0) {
                    const profiles: Record<number, UserProfile> = { ...userProfiles };
                    for (const id of missingIds) {
                        try {
                            const res = await galactus.findByUserId({ id: id });
                            const u = res.response?.user;
                            if (u) {
                                profiles[id] = {
                                    id: u.userId,
                                    username: u.username,
                                    avatarUrl: u.avatarUrl || defaultAvatar,
                                };
                            }
                        } catch (e) {
                            profiles[id] = {
                                id,
                                username: "Unknown",
                                avatarUrl: defaultAvatar,
                            };
                        }
                    }
                    setUserProfiles(profiles);
                }
            } catch (error) {
                setError("Failed to load messages");
            } finally {
                setLoading(false);
            }
        };

        fetchInitialMessages();

        const connectToChatStream = () => {
            if (!user || !friendId || cancelled) return;
            streamController = new AbortController();
            try {
                const stream = social.subscribeChatEvents(
                    { userId: user.userId, otherUserId: friendId },
                    { timeout: 60000 }
                );
                (async () => {
                    try {
                        for await (const event of stream.responses) {
                            if (cancelled) break;
                            if (event && event.event && event.event.oneofKind === "message" && event.event.message) {
                                const temp = event.event.message;
                                setMessages(prev => {
                                    // Remove typing bubble if present before pushing new message
                                    return prev.filter(m => m.id !== -9999).concat(event.event.message);
                                });
                                const senderId = event.event.message.senderId;
                                if (!(senderId in userProfiles)) {
                                    try {
                                        const res = await galactus.findByUserId({ id: senderId });
                                        const u = res.response?.user;
                                        if (u) {
                                            setUserProfiles(prev => ({
                                                ...prev,
                                                [senderId]: {
                                                    id: u.userId,
                                                    username: u.username,
                                                    avatarUrl: u.avatarUrl || defaultAvatar,
                                                }
                                            }));
                                        }
                                    } catch (e) {
                                        setUserProfiles(prev => ({
                                            ...prev,
                                            [senderId]: {
                                                id: senderId,
                                                username: "Unknown",
                                                avatarUrl: defaultAvatar,
                                            }
                                        }));
                                    }
                                }
                            }
                            else if (event && event.event && event.event.oneofKind === "typing" && event.event.typing) {
                                const typingId = event.event.typing.senderId;
                                if (typingId !== user?.userId) {
                                    setTypingUserId(typingId);
                                    // Push typing bubble if not present
                                    setMessages(prev => {
                                        // Only add if not already present
                                        if (!prev.some(m => m.id === -9999)) {
                                            return [
                                                ...prev,
                                                {
                                                    id: -9999,
                                                    senderId: typingId,
                                                    receiverId: user?.userId ?? 0,
                                                    content: "...",
                                                    createdAt: Date.now(),
                                                }
                                            ];
                                        }
                                        return prev;
                                    });
                                    if (typingTimeoutRef.current) clearTimeout(typingTimeoutRef.current);
                                    typingTimeoutRef.current = setTimeout(() => {
                                        setTypingUserId(null);
                                        // Remove typing bubble
                                        setMessages(prev => prev.filter(m => m.id !== -9999));
                                    }, 2000);
                                }
                            }
                            else if (event && event.event && event.event.oneofKind === "unsend" && event.event.unsend) {
                                const messageId = event.event.unsend.messageId;
                                setMessages(prev => prev.filter(m => m.id !== messageId));
                            }
                        }
                    } catch (err) {
                        if (!cancelled) {
                            // Attempt reconnection after delay
                            if (reconnectTimeout) clearTimeout(reconnectTimeout);
                            reconnectTimeout = setTimeout(() => {
                                if (!cancelled) {
                                    connectToChatStream();
                                }
                            }, 2000);
                        }
                    }
                })();
            } catch (err) {
                if (reconnectTimeout) clearTimeout(reconnectTimeout);
                reconnectTimeout = setTimeout(() => {
                    if (!cancelled) connectToChatStream();
                }, 2000);
            }
        };
        connectToChatStream();
        return () => {
            cancelled = true;
            if (reconnectTimeout) clearTimeout(reconnectTimeout);
            if (streamController) {
                streamController.abort();
            }
        };
    }, [friendId, user, social, userProfiles]);

    const sendTypingEvent = () => {
        if (!user || !friendId) return;
        social.sendChatEvent({
            event: {
                oneofKind: "typing",
                typing: {
                    senderId: user.userId,
                    receiverId: friendId,
                    isTyping: true,
                }
            }
        });
    };

    const sendMessage = async () => {
        if (!input.trim() || !friendId || !user) return;
        setSending(true);
        setError(null);
        try {
            const res = await social.sendChatEvent({
                event: {
                    oneofKind: "message",
                    message: {
                        id: 129092,
                        senderId: user.userId,
                        receiverId: friendId,
                        content: input,
                    }
                }
            });
            if (res?.response?.code !== 0) {
                setError(res?.response?.message || "Failed to send message");
            } else {
                setInput("");
            }
        } catch (error: any) {
            setError(error?.message || "Failed to send message");
        } finally {
            setSending(false);
        }
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setInput(e.target.value);
        sendTypingEvent();
    };

    const handleUnsend = async (messageId: number) => {
        if (!user || !friendId) return;
        try {
            // Use unary unsendMessage instead of sending an unsend event
            const res = await social.unsendMessage({
                messageId,
                userId: user.userId,
            });
            console.log(res.response)
            if (res?.response?.code === 0) {
                setMessages(prev => prev.filter(m => m.id !== messageId));
            }
        } catch (e) {}
    };

    function formatTimestamp(ts?: string | number | Date) {
        if (!ts) return "";
        let date: Date;
        if (typeof ts === "object" && ts !== null && 'seconds' in ts) {
            date = new Date(Number((ts as any).seconds) * 1000);
        } else if (typeof ts === "string") {
            const parsed = Date.parse(ts);
            if (!isNaN(parsed)) {
                date = new Date(parsed);
            } else if (!isNaN(Number(ts))) {
                date = new Date(Number(ts));
            } else {
                return ts;
            }
        } else if (typeof ts === "number") {
            date = new Date(ts);
        } else {
            date = ts as Date;
        }
        return date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
    }

    if (!user || !friendId) {
        return <div className="chat-window" style={{flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', color: '#888'}}>Select a chat to start messaging.</div>;
    }

    return (
        <div className="chat-window">
            {error && <div style={{ color: 'red', textAlign: 'center', marginBottom: 8 }}>{error}</div>}
            <div className="message-list" style={{ height: 400, overflowY: "auto" }}>
                {loading ? (
                    Array.from({ length: 8 }).map((_, i) => (
                        <div className="message-skeleton" key={i} />
                    ))
                ) : (
                    messages.map((msg, i) => {
                        // Typing bubble: id === -9999
                        const isTypingBubble = msg.id === -9999;
                        const profile = userProfiles[msg.senderId] || { username: "Unknown", avatarUrl: defaultAvatar };
                        const isSelf = msg.senderId === user?.userId;
                        const timestamp = msg.createdAt;
                        if (isTypingBubble) {
                            return (
                                <div
                                    className="message-row"
                                    key={"typing-bubble"}
                                    style={{
                                        display: 'flex',
                                        flexDirection: 'row',
                                        alignItems: 'flex-end',
                                        gap: 8,
                                        marginBottom: 12
                                    }}
                                >
                                    <img
                                        src={profile.avatarUrl || defaultAvatar}
                                        alt={profile.username || "Someone"}
                                        style={{ width: 32, height: 32, borderRadius: '50%', objectFit: 'cover' }}
                                    />
                                    <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'flex-start', maxWidth: 320 }}>
                                        <div style={{ fontWeight: 500, fontSize: 13, marginBottom: 2 }}>
                                            {profile.username || 'Someone'}
                                        </div>
                                        <div
                                            className="chat-bubble"
                                            style={{
                                                background: '#f5f5f5',
                                                color: '#888',
                                                borderRadius: 16,
                                                padding: '8px 14px',
                                                fontSize: 15,
                                                boxShadow: '0 1px 2px rgba(0,0,0,0.04)',
                                                wordBreak: 'break-word',
                                                minWidth: 40,
                                                maxWidth: 320,
                                                position: "relative",
                                                fontStyle: "italic"
                                            }}
                                        >
                                            ...
                                        </div>
                                    </div>
                                </div>
                            );
                        }
                        return (
                            <div
                                className={`message-row${isSelf ? " self" : ""}`}
                                key={i}
                                style={{ display: 'flex', flexDirection: isSelf ? 'row-reverse' : 'row', alignItems: 'flex-end', gap: 8, marginBottom: 12 }}
                            >
                                <img src={profile.avatarUrl} alt={profile.username} style={{ width: 32, height: 32, borderRadius: '50%', objectFit: 'cover' }} />
                                <div style={{ display: 'flex', flexDirection: 'column', alignItems: isSelf ? 'flex-end' : 'flex-start', maxWidth: 320 }}>
                                    <div style={{ fontWeight: 500, fontSize: 13, marginBottom: 2 }}>{profile.username}</div>
                                    <div
                                        className="chat-bubble"
                                        style={{
                                            background: isSelf ? '#e0f7fa' : '#f5f5f5',
                                            color: '#222',
                                            borderRadius: 16,
                                            padding: '8px 14px',
                                            fontSize: 15,
                                            boxShadow: '0 1px 2px rgba(0,0,0,0.04)',
                                            wordBreak: 'break-word',
                                            minWidth: 40,
                                            maxWidth: 320,
                                            position: "relative"
                                        }}
                                    >
                                        {msg.content}
                                        <span
                                            style={{
                                                display: "block",
                                                fontSize: 10,
                                                color: "#aaa",
                                                marginTop: 4,
                                                textAlign: isSelf ? "right" : "left"
                                            }}
                                        >
                                            {formatTimestamp(timestamp)}
                                        </span>
                                    </div>
                                    {isSelf && (
                                        <span
                                            style={{
                                                marginTop: 2,
                                                fontSize: 11,
                                                color: '#888',
                                                cursor: 'pointer',
                                                userSelect: 'none',
                                                padding: '1px 4px',
                                                borderRadius: 4,
                                                alignSelf: 'flex-end',
                                                transition: 'background 0.2s',
                                                lineHeight: 1.5,
                                            }}
                                            onClick={() => handleUnsend(msg.id)}
                                            title="Unsend"
                                            onMouseOver={e => (e.currentTarget.style.background = '#f0f0f0')}
                                            onMouseOut={e => (e.currentTarget.style.background = 'transparent')}
                                        >
                                            unsend
                                        </span>
                                    )}
                                </div>
                            </div>
                        );
                    })
                )}
            </div>
            <div ref={messagesEndRef} />
            <div className="chat-input-row">
                <input
                    className="chat-input"
                    value={input}
                    onChange={handleInputChange}
                    placeholder="Type a message..."
                    disabled={sending}
                    onKeyDown={e => { if (e.key === "Enter") sendMessage(); }}
                />
                <button className="send-btn" onClick={sendMessage} disabled={sending || !input.trim()}>
                    Send
                </button>
            </div>
        </div>
    );
}