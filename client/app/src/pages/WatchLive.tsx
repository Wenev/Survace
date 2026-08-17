import { useEffect, useState } from "react";
import {
    StreamVideoClient,
    StreamVideo,
    StreamCall,
    Call,
} from "@stream-io/video-react-sdk";
import { useGrpc } from "../context/ClientContext";
import WatcherUI from "./WatcherUI";
import { useNavigate, useParams } from "react-router-dom";
import Navbar from "../component/Navbar";
import MobileNavbar from "../component/MobileNavbar";

export default function WatchLive() {
    const { stream } = useGrpc();
    const { callId } = useParams<{ callId: string }>();
    const navigate = useNavigate();
    const [client, setClient] = useState<StreamVideoClient | null>(null);
    const [call, setCall] = useState<Call | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [isMobile, setIsMobile] = useState(false);

    useEffect(() => {
        async function initializeViewer() {
            try {
                const userId = "-1";
                const tokenResponse = await stream.getStreamToken({
                    userId: userId
                });

                if (tokenResponse.response.code !== 0) {
                    setError(`Failed to get stream token: ${tokenResponse.response.message}`);
                    setLoading(false);
                    return;
                }

                const token = tokenResponse.response.token;
                const apiKey = import.meta.env.VITE_STREAM_API_KEY;

                if (!apiKey) {
                    setError("Stream API key is not configured. Please check your environment variables.");
                    setLoading(false);
                    return;
                }

                const streamUser = {
                    id: userId,
                    name: "Viewer"
                };

                const streamClient = new StreamVideoClient({
                    apiKey,
                    user: streamUser,
                    token
                });

                setClient(streamClient);

                // Use callId from params
                if (!callId) {
                    setError("No callId provided in URL.");
                    setLoading(false);
                    return;
                }
                const streamCall = streamClient.call("livestream", callId);
                await streamCall.join();
                setCall(streamCall);

                setLoading(false);

            } catch (error) {
                setError("Failed to join live stream. Please try again.");
                setLoading(false);
            }
        }

        if (callId) {
            initializeViewer();
        }

        return () => {
            if (call) {
                call.leave();
            }
            if (client) {
                client.disconnectUser();
            }
        };
    }, [callId, stream]);

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
            <MobileNavbar onToggleSidebar={() => setSidebarOpen((open) => !open)} isOpen={isMobile && sidebarOpen} />
            <Navbar isOpen={sidebarOpen} />
            <main className="main-content" style={{ width: "100%", minHeight: "100vh", display: "flex", justifyContent: "center", alignItems: "center", background: "var(--color-bg)" }}>
                <div style={{ width: "100%", maxWidth: 800, display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center" }}>
                    <button className="streamer-ui-btn stop-live" style={{ alignSelf: "flex-start", margin: "24px 0" }} onClick={() => navigate("/live")}>
                        ← Back to Streams
                    </button>
                    <h1 style={{ color: "#fe2c55", fontWeight: 700, marginBottom: 24, fontSize: 32 }}>
                        Watching Live Stream
                    </h1>
                    {loading && <div>Loading stream...</div>}
                    {error && <div className="error">{error}</div>}
                    {!loading && !error && client && call ? (
                        <StreamVideo client={client}>
                            <StreamCall call={call}>
                                <WatcherUI />
                            </StreamCall>
                        </StreamVideo>
                    ) : null}
                </div>
            </main>
        </div>
    );
}