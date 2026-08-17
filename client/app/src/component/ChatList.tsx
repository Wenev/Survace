import { useEffect, useState, useRef } from "react";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import type {Follow, Friend} from "../generated/dto/social.ts";
import defaultAvatar from "../assets/default.jpg";

const ChatList = ({ selectedId, onSelect }: {
    friendds: Friend[]
    selectedId: number | null;
    onSelect: (id: number) => void;
}) => {
    const { user } = useAuth();
    const { social, galactus } = useGrpc();
    const [following, setFollowing] = useState<Friend[]>([]);
    const [userProfiles, setUserProfiles] = useState<Record<number, { username: string; avatarUrl: string }>>({});
    const [loading, setLoading] = useState(true);
    const [search, setSearch] = useState("");
    const [debouncedSearch, setDebouncedSearch] = useState("");
    const debounceTimeout = useRef<NodeJS.Timeout | null>(null);

    useEffect(() => {
        if (debounceTimeout.current) clearTimeout(debounceTimeout.current);
        debounceTimeout.current = setTimeout(() => {
            setDebouncedSearch(search);
        }, 300);
        return () => {
            if (debounceTimeout.current) clearTimeout(debounceTimeout.current);
        };
    }, [search]);

    useEffect(() => {
        const fetchChatPartners = async () => {
            if (!user) {
                setFollowing([]);
                setLoading(false);
                return;
            }
            setLoading(true);
            try {
                const res = await social.listFriends({ userId: user.userId });

                setFollowing(res.response.data);
            } finally {
                setLoading(false);
            }
        };
        fetchChatPartners();
    }, [user, social, debouncedSearch]);

    return (
        <div className="chat-list" style={{ height: 400, overflowY: "auto" }}>
            <div className="chat-list-header">Chat List</div>
            <input
                type="text"
                placeholder="Search user handle..."
                value={search}
                onChange={e => setSearch(e.target.value)}
                style={{ width: "90%", margin: "8px auto", display: "block", padding: 6, borderRadius: 6, border: "1px solid #ccc" }}
            />
            {loading ? (
                <>
                    {Array.from({ length: 6 }).map((_, i) => (
                        <div className="message-skeleton" key={i} style={{ width: "80%", margin: "16px auto", height: 40, borderRadius: 8, background: "#f0f0f0" }} />
                    ))}
                </>
            ) : following.length === 0 ? (
                <div className="no-friends-skeleton">
                    <div className="no-friends-icon">👤</div>
                    <div className="no-friends-text">No chat found.</div>
                </div>
            ) : (
                following.map((f: Friend) => (
                    <div
                        className={`chat-list-item${selectedId === f.friendId? " selected" : ""}`}
                        key={f.friendId}
                        onClick={() => onSelect(f.friendId)}
                    >
                        <img
                            src={userProfiles[f.friendId]?.avatarUrl || defaultAvatar}
                            alt={userProfiles[f.friendId]?.username || "Unknown"}
                            style={{ width: 28, height: 28, borderRadius: '50%', objectFit: 'cover', marginRight: 8 }}
                        />
                        {userProfiles[f.friendId]?.username || `User #${f.friendId}`}
                    </div>
                ))
            )}
        </div>
    );
};

export default ChatList;