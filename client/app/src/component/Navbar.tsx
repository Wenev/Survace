"use client"

import type React from "react"
import defaultAvatar from "../assets/default.jpg";
import { useState, useEffect } from "react"
import "../style/component/navbar.css"
import {
    Home,
    Compass,
    Users,
    UserPlus,
    Upload,
    Zap,
    MessageCircle,
    Radio,
    User,
    Settings,
    LogOut,
    Moon,
    Sun,
} from "lucide-react"
import {useNavigate} from "react-router-dom";
import { useAuth } from "../context/AuthContext";
import { useGrpc } from "../context/ClientContext";
import type {FollowingVideoRequest, FollowingVideoResponse} from "../generated/dto/video.ts";
import type { Follow} from "../generated/dto/social.ts";

interface SidebarProps {
    isOpen: boolean
    searchQuery: string
    setSearchQuery: (q: string) => void
}

interface NavItem {
    icon: JSX.Element
    label: string
    href: string
    active?: boolean
}

const navItems: NavItem[] = [
    { icon: <Home size={18} />, label: "For You", href: "/", active: true },
    { icon: <Users size={18} />, label: "Following", href: "/following" },
    { icon: <UserPlus size={18} />, label: "Friends", href: "/friends" },
    { icon: <Upload size={18} />, label: "Upload", href: "/upload" },
    { icon: <MessageCircle size={18} />, label: "Chat", href: "/Chat" },
    { icon: <Radio size={18} />, label: "Live", href: "/live" },
    { icon: <User size={18} />, label: "Profile", href: "/profile" },
    { icon: <Settings size={18} />, label: "Settings", href: "/settings" },
    { icon: <LogOut size={18} />, label: "Log Out", href: "/logout" },
]

export default function Sidebar({ isOpen, searchQuery, setSearchQuery }: SidebarProps) {
    const navigate = useNavigate()
    const { user, logout } = useAuth();
    const { social, galactus } = useGrpc();
    const [darkMode, setDarkMode] = useState(false)
    const [following, setFollowing] = useState<any[]>([]);

    useEffect(() => {
        const savedTheme = localStorage.getItem("theme")
        const systemPrefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches

        const shouldUseDarkMode = savedTheme === "dark" || (!savedTheme && systemPrefersDark)
        setDarkMode(shouldUseDarkMode)

        if (shouldUseDarkMode) {
            document.body.classList.add("dark-mode")
        } else {
            document.body.classList.remove("dark-mode")
        }
    }, [])

    useEffect(() => {
        async function fetchFollowing() {
            if (user && user.userId && social) {
                try {
                    const res = await social.listFollowing({ userId: user.userId });
                    const follows: Follow[] = res.response.data || [];

                    const enriched = await Promise.all(
                        follows.map(async (f) => {
                            try {
                                const userRes = await galactus.findByUserId({ id: f.followeeId });
                                const userData = userRes.response.user;
                                return {
                                    ...f,
                                    userId: userData?.userId,
                                    username: userData?.username,
                                    email: userData?.email,
                                    avatarUrl: userData?.avatarUrl,
                                };
                            } catch {
                                return f;
                            }
                        })
                    );
                    setFollowing(enriched);
                } catch {
                    setFollowing([]);
                }
            } else {
                setFollowing([]);
            }
        }
        fetchFollowing();
    }, [user, social]);

    const toggleTheme = () => {
        const newDarkMode = !darkMode
        setDarkMode(newDarkMode)

        if (newDarkMode) {
            document.body.classList.add("dark-mode")
            localStorage.setItem("theme", "dark")
        } else {
            document.body.classList.remove("dark-mode")
            localStorage.setItem("theme", "light")
        }
    }

    const filteredNavItems = navItems.filter(item => {
        if (!user) {
            return ["For You"].includes(item.label);
        }
        if (user && ["Log Out"].includes(item.label)) return true;
        if (user) return item.label !== "Register" && item.label !== "Login";
        return item.label !== "Log Out";
    });

    return (
        <nav className={`sidebar ${isOpen ? "active" : ""}`}>
            <div className="sidebar-header">
                <div className="logo">
                    <span className="tiktok-icon">🎵</span>
                    TikTok
                </div>
                <button aria-label="Toggle dark mode" className="theme-toggle-btn" onClick={toggleTheme}>
                    {darkMode ? <Sun size={20} /> : <Moon size={20} />}
                </button>
            </div>

            <div className="sidebar-content">
                <form className="search-container">
                    <input
                        type="text"
                        placeholder="Search"
                        className="search-input"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                    />
                    <span className="search-icon">🔍</span>
                </form>

                <ul className="nav-menu">
                    {filteredNavItems.map((item, index) => (
                        <li key={index} className={`nav-item ${item.active ? "active" : ""}`}>
                            {item.label === "Log Out" ? (
                                <button
                                    className="nav-link"
                                    style={{ width: "100%", background: "none", border: "none", textAlign: "left", padding: 0, cursor: "pointer" }}
                                    onClick={() => {
                                        logout();
                                        navigate("/login");
                                    }}
                                >
                                    <span className="nav-icon">{item.icon}</span>
                                    {item.label}
                                </button>
                            ) : (
                                <button
                                    className="nav-link"
                                    style={{ width: "100%", background: "none", border: "none", textAlign: "left", padding: 0, cursor: "pointer" }}
                                    onClick={() => navigate(`${item.href}`)}
                                >
                                    <span className="nav-icon">{item.icon}</span>
                                    {item.label}
                                </button>
                            )}
                        </li>
                    ))}
                </ul>

                <div className="sidebar-footer">
                    {!user && (
                        <>
                            <button type="button" className="login-btn" onClick={() => {
                                navigate("/signup")
                            }}>
                                Register
                            </button>
                            <button type="button" className="login-btn" onClick={() => {
                                navigate("/login")
                            }}>
                                Login
                            </button>
                        </>
                    )}
                    {user && (
                        <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
                            <div style={{ display: "flex", alignItems: "center", gap: 10 }}>
                                <img
                                    src={user.avatarUrl || defaultAvatar}
                                    style={{
                                        width: 32,
                                        height: 32,
                                        borderRadius: "50%",
                                        objectFit: "cover",
                                        border: "2px solid var(--color-border)"
                                    }}
                                    loading="lazy"
                                />
                                <span className="sidebar-user">
                                    Logged in as <b>{user.username || user.email}</b>
                                </span>
                            </div>
                            {following.length > 0 && (
                                <div className="sidebar-following-section" style={{ marginTop: 12 }}>
                                    <div style={{ fontWeight: 700, marginBottom: 8, fontSize: 15 }}>Following</div>
                                    <ul style={{ listStyle: "none", padding: 0, margin: 0 }}>
                                        {following.map((f) => (
                                            <li key={f.userId} style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 8, cursor: "pointer" }}
                                                onClick={() => navigate(`/profile/${f.userId}`)}>
                                                <img
                                                    src={f.avatarUrl || defaultAvatar}
                                                    alt={f.username}
                                                    style={{
                                                        width: 28,
                                                        height: 28,
                                                        borderRadius: "50%",
                                                        objectFit: "cover",
                                                        border: "1px solid var(--color-border)"
                                                    }}
                                                    loading="lazy"
                                                />
                                                <span style={{ fontSize: 14 }}>{f.username || f.email}</span>
                                            </li>
                                        ))}
                                    </ul>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            </div>
        </nav>
    )
}