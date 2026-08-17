"use client"

import { useState, useEffect } from "react"
import Navbar from "../component/Navbar"
import MobileNavbar from "../component/MobileNavbar"
import FriendFollowingFeed from "../component/FriendFollowingFeed"
import ErrorToast from "../component/ErrorToast"
import "../style/home.css"
import "../style/component/friend-feed.css";

export default function FriendPage() {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [isMobile, setIsMobile] = useState(false)
    const [searchQuery, setSearchQuery] = useState("")
    const [error, setError] = useState<string | null>(null)

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
    }, [])

    const toggleSidebar = () => {
        if (isMobile) setSidebarOpen((open) => !open);
    }
    const closeOverlay = () => {
        setSidebarOpen(false)
    }

    return (
        <div className="app">
            <MobileNavbar onToggleSidebar={toggleSidebar} isOpen={isMobile && sidebarOpen} searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
            <Navbar isOpen={sidebarOpen} searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
            <ErrorToast message={error || ""} onClose={() => setError(null)} />
            <main className="main-content">
                <FriendFollowingFeed mode="friends" />
            </main>
            {(sidebarOpen) && <div className="overlay active" onClick={closeOverlay} />}
        </div>
    )
}