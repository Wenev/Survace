import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import "../style/profile.css";
import Navbar from "../component/Navbar";
import MobileNavbar from "../component/MobileNavbar";
import ErrorToast from "../component/ErrorToast";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import defaultAvatar from "../assets/default.jpg";
import { Settings, Trash2, Eye, EyeOff, Bell, Lock, Unlock, MessageCircle, MessageCircleOff, ArrowLeft } from "lucide-react";

export default function SettingPage() {
    const { user: currentUser, logout } = useAuth();
    const { galactus } = useGrpc();
    const navigate = useNavigate();

    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [isMobile, setIsMobile] = useState(false);

    useEffect(() => {
        const checkMobile = () => {
            const mobile = window.innerWidth <= 768;
            setIsMobile(mobile);
            if (!mobile) {
                setSidebarOpen(false);
            }
        };
        checkMobile();
        window.addEventListener("resize", checkMobile);
        return () => window.removeEventListener("resize", checkMobile);
    }, []);

    const toggleSidebar = () => {
        if (isMobile) setSidebarOpen((open) => !open);
    };
    const closeOverlay = () => {
        setSidebarOpen(false);
    };

    const [username, setUsername] = useState("");
    const [bio, setBio] = useState("");
    const [avatarUrl, setAvatarUrl] = useState("");
    const [enableNotif, setEnableNotif] = useState(true);
    const [visibility, setVisibility] = useState<"public" | "private">("public");
    const [privateAccount, setPrivateAccount] = useState(false);
    const [editChatRestriction, setEditChatRestriction] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);

    const [notifFollower, setNotifFollower] = useState(true);
    const [notifMessage, setNotifMessage] = useState(true);
    const [notifMentions, setNotifMentions] = useState(true);

    useEffect(() => {
        if (!currentUser) {
            navigate("/login");
            return;
        }
        const fetchUserAndSetting = async () => {
            setIsLoading(true);
            try {
                const userRes = await galactus.findByUserId({ id: currentUser.userId });
                const user = userRes.response?.user || userRes.user;
                setUsername(user?.username || "");
                setBio("");
                setAvatarUrl(user?.avatarUrl || defaultAvatar);

                const settingRes = await galactus.findUserSetting({ id: currentUser.userId });
                const setting = settingRes.response?.setting;
                if (setting) {
                    setNotifFollower(!!setting.newFollowerNotification);
                    setNotifMessage(!!setting.messageNotification);
                    setNotifMentions(!!setting.mentionsNotification);
                    setEnableNotif(
                        !!setting.newFollowerNotification ||
                        !!setting.messageNotification ||
                        !!setting.mentionsNotification
                    );
                    setVisibility(setting.likeTabVisibility ? "public" : "private");
                    setPrivateAccount(!!setting.private);
                    setEditChatRestriction(setting.chatRestriction === "none");
                }
            } catch (e: any) {
                setError(e?.message || "Failed to load settings");
            }
            setIsLoading(false);
        };
        fetchUserAndSetting();
    }, [currentUser, galactus, navigate]);

    // Only update username and avatar via updateUser
    const handleUpdate = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setSuccess(null);
        try {
            await galactus.updateUser({
                userId: currentUser.userId,
                username,
                avatarUrl,
            });
            setSuccess("Profile updated successfully.");
        } catch (e: any) {
            setError(e?.message || "Failed to update profile.");
        }
    }
    const handlePrivateAccount = async (checked: boolean) => {
        setPrivateAccount(checked);
        try {
            await galactus.enablePrivateAccount({
                userId: currentUser.userId,
                enable: checked,
            });
        } catch (e: any) {
            setError(e?.message || "Failed to update private account setting.");
        }
    };

    const handleEnableNotif = async (checked: boolean) => {
        setEnableNotif(checked);
        try {
            await galactus.enableNewFollowerNotification({
                userId: currentUser.userId,
                enable: checked,
            });
            await galactus.enableMessageNotification({
                userId: currentUser.userId,
                enable: checked,
            });
            await galactus.enableMentionsNotification({
                userId: currentUser.userId,
                enable: checked,
            });
        } catch (e: any) {
            setError(e?.message || "Failed to update notification setting.");
        }
    };

    const handleVisibility = async (value: "public" | "private") => {
        setVisibility(value);
        try {
            await galactus.enableLikeTabVisibility({
                userId: currentUser.userId,
                enable: value === "public",
            });
        } catch (e: any) {
            setError(e?.message || "Failed to update visibility.");
        }
    };

    const handleEditChatRestriction = async (checked: boolean) => {
        setEditChatRestriction(checked);
        try {
            await galactus.editChatRestriction({
                userId: currentUser.userId,
                chatRestriction: checked ? "none" : "everyone",
            });
        } catch (e: any) {
            setError(e?.message || "Failed to update chat restriction.");
        }
    };

    // Handle delete account
    const handleDeleteAccount = async () => {
        setError(null);
        setSuccess(null);
        try {
            await galactus.deleteAccount({ userId: currentUser.userId });
            logout();
            navigate("/login");
        } catch (e: any) {
            setError(e?.message || "Failed to delete account.");
        }
    };

    // Back button handler
    const handleBack = () => {
        navigate("/");
    };

    // Individual notification toggles
    const handleNotifFollower = async (checked: boolean) => {
        setNotifFollower(checked);
        try {
            await galactus.enableNewFollowerNotification({
                userId: currentUser.userId,
                enable: checked,
            });
        } catch (e: any) {
            setError(e?.message || "Failed to update follower notification setting.");
        }
    };
    const handleNotifMessage = async (checked: boolean) => {
        setNotifMessage(checked);
        try {
            await galactus.enableMessageNotification({
                userId: currentUser.userId,
                enable: checked,
            });
        } catch (e: any) {
            setError(e?.message || "Failed to update message notification setting.");
        }
    };
    const handleNotifMentions = async (checked: boolean) => {
        setNotifMentions(checked);
        try {
            await galactus.enableMentionsNotification({
                userId: currentUser.userId,
                enable: checked,
            });
        } catch (e: any) {
            setError(e?.message || "Failed to update mentions notification setting.");
        }
    };

    return (
        <div className="app">
            <MobileNavbar onToggleSidebar={toggleSidebar} isOpen={isMobile && sidebarOpen} />
            <Navbar isOpen={sidebarOpen} />
            <ErrorToast message={error || ""} onClose={() => setError(null)} />
            <main className="main-content profile-main-content">
                <div className="profile-container">
                    <div className="profile-header">
                        <button
                            className="back-button"
                            style={{
                                background: "none",
                                border: "none",
                                cursor: "pointer",
                                marginRight: 16,
                                display: "flex",
                                alignItems: "center"
                            }}
                            onClick={handleBack}
                            aria-label="Back to Home"
                        >
                            <ArrowLeft size={24} />
                        </button>
                        <div className="profile-avatar-container">
                            <img src={avatarUrl || defaultAvatar} alt={username} className="profile-avatar" />
                        </div>
                        <div className="profile-info">
                            <div className="profile-username">
                                <h1>Settings</h1>
                            </div>
                        </div>
                    </div>
                    <div className="profile-content">
                        {isLoading ? (
                            <div style={{ textAlign: "center", padding: "2rem" }}>Loading...</div>
                        ) : (
                        <form className="settings-form" onSubmit={handleUpdate}>
                            {/* Username */}
                            <div className="form-group">
                                <label>Username</label>
                                <input
                                    type="text"
                                    value={username}
                                    onChange={e => setUsername(e.target.value)}
                                    required
                                    disabled={isLoading}
                                />
                            </div>
                            {/* Bio */}
                            <div className="form-group">
                                <label>Bio</label>
                                <textarea
                                    value={bio}
                                    onChange={e => setBio(e.target.value)}
                                    rows={3}
                                    disabled={isLoading}
                                />
                            </div>
                            {/* Avatar URL */}
                            <div className="form-group">
                                <label>Avatar URL</label>
                                <input
                                    type="text"
                                    value={avatarUrl}
                                    onChange={e => setAvatarUrl(e.target.value)}
                                    disabled={isLoading}
                                />
                            </div>
                            {/* Enable Notifications - Follower */}
                            <div className="form-group switch-group">
                                <label>
                                    <Bell size={16} style={{ marginRight: 4 }} />
                                    Follower Notifications
                                </label>
                                <input
                                    type="checkbox"
                                    checked={notifFollower}
                                    onChange={e => handleNotifFollower(e.target.checked)}
                                    disabled={isLoading}
                                />
                            </div>
                            {/* Enable Notifications - Message */}
                            <div className="form-group switch-group">
                                <label>
                                    <MessageCircle size={16} style={{ marginRight: 4 }} />
                                    Message Notifications
                                </label>
                                <input
                                    type="checkbox"
                                    checked={notifMessage}
                                    onChange={e => handleNotifMessage(e.target.checked)}
                                    disabled={isLoading}
                                />
                            </div>
                            {/* Enable Notifications - Mentions */}
                            <div className="form-group switch-group">
                                <label>
                                    <Bell size={16} style={{ marginRight: 4 }} />
                                    Mentions Notifications
                                </label>
                                <input
                                    type="checkbox"
                                    checked={notifMentions}
                                    onChange={e => handleNotifMentions(e.target.checked)}
                                    disabled={isLoading}
                                />
                            </div>
                            {/* Visibility */}
                            <div className="form-group switch-group">
                                <label>
                                    {visibility === "public" ? <Eye size={16} /> : <EyeOff size={16} />}
                                    Visibility: {visibility.charAt(0).toUpperCase() + visibility.slice(1)}
                                </label>
                                <select
                                    value={visibility}
                                    onChange={e => handleVisibility(e.target.value as "public" | "private")}
                                    disabled={isLoading}
                                >
                                    <option value="public">Public</option>
                                    <option value="private">Private</option>
                                </select>
                            </div>
                            {/* Private Account */}
                            <div className="form-group switch-group">
                                <label>
                                    {privateAccount ? <Lock size={16} /> : <Unlock size={16} />}
                                    Private Account
                                </label>
                                <input
                                    type="checkbox"
                                    checked={privateAccount}
                                    onChange={e => handlePrivateAccount(e.target.checked)}
                                    disabled={isLoading}
                                />
                            </div>
                            {/* Edit Chat Restriction */}
                            <div className="form-group switch-group">
                                <label>
                                    {editChatRestriction ? <MessageCircleOff size={16} /> : <MessageCircle size={16} />}
                                    Restrict Chat Editing
                                </label>
                                <input
                                    type="checkbox"
                                    checked={editChatRestriction}
                                    onChange={e => handleEditChatRestriction(e.target.checked)}
                                    disabled={isLoading}
                                />
                            </div>
                            {/* Save Button */}
                            <div className="form-actions">
                                <button type="submit" className="edit-profile-button" disabled={isLoading}>
                                    <Settings size={16} /> Save Changes
                                </button>
                            </div>
                            {success && <div className="success-message">{success}</div>}
                            {error && <div className="profile-error">{error}</div>}
                        </form>
                        )}
                        {/* Delete Account */}
                        <div className="delete-account-section">
                            <button
                                className="delete-account-button"
                                onClick={() => setShowDeleteConfirm(true)}
                                style={{ color: "#e74c3c", marginTop: 24 }}
                                disabled={isLoading}
                            >
                                <Trash2 size={16} /> Delete Account
                            </button>
                            {showDeleteConfirm && (
                                <div className="delete-confirm-modal">
                                    <div className="modal-content">
                                        <p>Are you sure you want to delete your account? This action cannot be undone.</p>
                                        <div className="modal-actions">
                                            <button onClick={handleDeleteAccount} className="confirm-delete-button" disabled={isLoading}>Yes, Delete</button>
                                            <button onClick={() => setShowDeleteConfirm(false)} disabled={isLoading}>Cancel</button>
                                        </div>
                                    </div>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </main>
            {(sidebarOpen) && <div className="overlay active" onClick={closeOverlay} />}
        </div>
    );
}