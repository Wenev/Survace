import  { useEffect, useState, useRef } from "react";
import "../style/chat.css";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import type { Friend, ChatMessage, ListFriendsResponse, ListMessagesResponse } from "../generated/dto/social";
import Navbar from "../component/Navbar";
import MobileNavbar from "../component/MobileNavbar";
import ChatWindow from "../component/ChatWindow.tsx";
import ChatList from "../component/ChatList.tsx";
import ErrorToast from "../component/ErrorToast";

const ChatPage = () => {
    const { user } = useAuth();
    const { social } = useGrpc();
    const [friends, setFriends] = useState<Friend[]>([]);
    const [loadingFriends, setLoadingFriends] = useState(true);
    const [selectedFriend, setSelectedFriend] = useState<number | null>(null);
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [isMobile, setIsMobile] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const checkMobile = () => {
            setIsMobile(window.innerWidth <= 768)
            if (window.innerWidth > 768) {
                setSidebarOpen(false)
            }
        }
        checkMobile()
        window.addEventListener("resize", checkMobile)
        return () => window.removeEventListener("resize", checkMobile)
    }, [])

    useEffect(() => {
        if (!user) return;
        setLoadingFriends(true);
        social.listFriends({ userId: user.userId }).then((res: ListFriendsResponse) => {
            console.log(res.response.data)
            setFriends(res.response.data || []);
            setLoadingFriends(false);
            if (res.response.data && res.response.data.length > 0) setSelectedFriend(res.response.data[0].friendId);
        }).catch((e: any) => {
            setError(e?.message || "Failed to load friends");
            setLoadingFriends(false);
        });
    }, [user, social]);

    if (!user) return <div className="chat-page">Please login to use chat.</div>;

    return (
        <div className="app">
            <MobileNavbar onToggleSidebar={() => setSidebarOpen(!sidebarOpen)} isOpen={sidebarOpen} />
            <Navbar isOpen={sidebarOpen} />
            <ErrorToast message={error || ""} onClose={() => setError(null)} />
            <main className="main-content chat-main-content">
                <div className="chat-page" style={{ display: 'flex', height: '100%', minHeight: '70vh' }}>
                    <div className="chat-list-container" style={{ flex: '0 0 300px', maxWidth: 340, minWidth: 220, borderRight: '1px solid #eee', background: '#fafafa', height: '100%', overflowY: 'auto' }}>
                        <ChatList
                            friends={friends}
                            selectedId={selectedFriend}
                            onSelect={setSelectedFriend}
                            loading={loadingFriends}
                        />
                        {!loadingFriends && friends.length === 0 && (
                            <div className="no-friends-skeleton">
                                <div className="no-friends-icon">👥</div>
                                <div className="no-friends-text">You have no friends yet.<br/>Add friends to start chatting!</div>
                            </div>
                        )}
                    </div>
                    <div className="chat-window-container" style={{ flex: 1, minWidth: 0, display: 'flex', flexDirection: 'column', height: '100%' }}>
                        <ChatWindow friendId={selectedFriend} selfId={user.userId} />
                    </div>
                </div>
            </main>
        </div>
    );
};

export default ChatPage;
