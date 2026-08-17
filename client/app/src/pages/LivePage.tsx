import { useEffect, useState } from "react";
import Navbar from "../component/Navbar";
import MobileNavbar from "../component/MobileNavbar";
import ErrorToast from "../component/ErrorToast";
import "../style/component/streamer-ui.css";
import LiveStreamCard from "../component/LiveStreamCard";
import type {LiveStreamInfo} from "../generated/dto/stream.ts";
import { useNavigate } from "react-router-dom";
import {useAuth} from "../context/AuthContext.tsx";
import { useGrpc } from "../context/ClientContext";

export default function LivePage() {
    const { user } = useAuth();
    const { stream } = useGrpc();
    const [liveStreams, setLiveStreams] = useState<LiveStreamInfo[]>([]);
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [isMobile, setIsMobile] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate();

    useEffect(() => {
        async function fetchLive() {
            try {
                const res = await stream.listLiveStreams({});
                setLiveStreams(res.response.streams || []);
            } catch (e: any) {
                setLiveStreams([]);
                setError(e?.message || "Failed to load live streams");
            }
        }
        fetchLive();
    }, [stream]);

    const handleWatch = (callId: string) => {
        navigate(`/live/${callId}`);
    };

    const toggleSidebar = () => {
        if (isMobile) setSidebarOpen((open) => !open);
    };

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

    return (
        <div className="app">
            <MobileNavbar onToggleSidebar={toggleSidebar} isOpen={isMobile && sidebarOpen} />
            <Navbar isOpen={sidebarOpen} />
            <ErrorToast message={error || ""} onClose={() => setError(null)} />
            <main className="main-content" style={{ width: "100%", minHeight: "100vh", display: "flex", flexDirection: "column", alignItems: "center", background: "var(--color-bg)" }}>
                <h1 style={{ color: "#fe2c55", fontWeight: 700, margin: "32px 0 24px 0", fontSize: 32 }}>Browse Live Streams</h1>
                <button
                    className="streamer-ui-btn go-live"
                    style={{ marginBottom: 24, fontWeight: 600, fontSize: 18, padding: "10px 32px" }}
                    onClick={() => navigate("/live/streaming")}
                >
                    Go Live
                </button>
                <div style={{ display: "flex", flexWrap: "wrap", justifyContent: "center", gap: 32, width: "100%", maxWidth: 1200 }}>
                    {liveStreams.length === 0 ? (
                        <div style={{ color: '#a8a8b3', fontSize: 20, marginTop: 40 }}>No live streams right now.</div>
                    ) : (
                        liveStreams.map((stream) => (
                            <LiveStreamCard
                                key={stream.callId}
                                stream={{ callId: stream.callId, title: `Live by ${stream.userId}` }}
                                onWatch={() => handleWatch(stream.callId)}
                            />
                        ))
                    )}
                </div>
            </main>
        </div>
    );
}