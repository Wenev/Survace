import { useEffect, useState } from "react";
import {
    StreamVideoClient,
    StreamVideo,
    StreamCall,
    Call,
    useCall,
    useCallStateHooks
} from "@stream-io/video-react-sdk";
import { useAuth } from "../context/AuthContext";
import { useGrpc } from "../context/ClientContext";
import StreamerUI from "../component/StreamerUI.tsx";
import { useNavigate } from "react-router-dom";
import Navbar from "../component/Navbar";
import MobileNavbar from "../component/MobileNavbar";

export default function GoLivePage() {
    const { user } = useAuth();
    const { stream } = useGrpc();
    const navigate = useNavigate();
    const [client, setClient] = useState<StreamVideoClient | null>(null);
    const [call, setCall] = useState<Call | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);


    const [screenShareError, setScreenShareError] = useState<string | null>(null);

    const [isScreenSharing, setIsScreenSharing] = useState(false);

    useEffect(() => {
        if (!user) return;
        async function initializeStream() {
            try {
                console.log("Fetching stream token for user:", user.userId);

                const tokenResponse = await stream.getStreamToken({
                    userId: (user.userId).toString()
                });

                console.log("Token response:", tokenResponse);

                if (tokenResponse.response.code !== 0) {
                    setError(`Failed to get stream token: ${tokenResponse.response.message}`);
                    setLoading(false);
                    return;
                }

                const token = tokenResponse.response.token;
                const callId = tokenResponse.response.callId;

                const apiKey = import.meta.env.VITE_STREAM_API_KEY;

                console.log(apiKey)

                if (!apiKey) {
                    setError("Stream API key is not configured. Please check your environment variables.");
                    setLoading(false);
                    return;
                }

                console.log("Using Stream API key:", apiKey);
                console.log("Using call ID:", callId);

                const userId = user.userId.toString();
                const streamUser = {
                    id: userId,
                    name: user.username || "Anonymous"
                };

                const streamClient = new StreamVideoClient({
                    apiKey,
                    user: streamUser,
                    token
                });

                setClient(streamClient);

                const streamCall = streamClient.call("livestream", callId);
                await streamCall.join({ create: true });
                setCall(streamCall);

                console.log("Successfully joined stream");
                setLoading(false);

            } catch (error) {
                console.error("Error setting up stream:", error);
                setError("Failed to initialize stream. Please try again.");
                setLoading(false);
            }
        }

        initializeStream();

        return () => {
            if (call) {
                console.log("Leaving call");
                call.leave();
            }
            if (client) {
                console.log("Disconnecting user");
                client.disconnectUser();
            }
        };
    }, [user, stream]);

    return (
        <div className="app">
            <MobileNavbar />
            <Navbar />
            <main className="main-content" style={{ width: "100%", minHeight: "100vh", display: "flex", justifyContent: "center", alignItems: "center", background: "var(--color-bg)" }}>
                <div style={{ width: "100%", maxWidth: 800, display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center" }}>
                    {loading && <div>Loading stream...</div>}
                    {error && <div className="error">{error}</div>}
                    {!loading && !error && client && call ? (
                        <StreamVideo client={client}>
                            <StreamCall call={call}>
                                <StreamerUI />
                            </StreamCall>
                        </StreamVideo>
                    ) : null}
                </div>
            </main>
        </div>
    );
}