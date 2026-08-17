"use client"

import { useState, useEffect } from "react"
import Navbar from "../component/Navbar"
import MobileNavbar from "../component/MobileNavbar"
import Feed from "../component/Feed"
import ErrorToast from "../component/ErrorToast"
import "../style/home.css"

export default function HomePage() {
    const [sidebarOpen, setSidebarOpen] = useState(false)
    const [activeVideoId, setActiveVideoId] = useState<string | null>(null)
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
                <Feed activeVideoId={activeVideoId} onVideoChange={setActiveVideoId} searchQuery={searchQuery} setError={setError} />
            </main>
            {(sidebarOpen) && <div className="overlay active" onClick={closeOverlay} />}
        </div>
    )
}
