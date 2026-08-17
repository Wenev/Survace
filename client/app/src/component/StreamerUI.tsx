import { ParticipantView, useCall, useCallStateHooks } from "@stream-io/video-react-sdk";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import "../style/component/streamer-ui.css";

const StreamerUI = () => {
    const {
        useCameraState,
        useMicrophoneState,
        useParticipantCount,
        useIsCallLive,
        useParticipants,
        useScreenShareState
    } = useCallStateHooks();

    const { camera: cam, isEnabled: isCamEnabled } = useCameraState();
    const { microphone: mic, isEnabled: isMicEnabled } = useMicrophoneState();
    const { isEnabled: isScreenEnabled, screenShare: screen } = useScreenShareState()

    const participantCount = useParticipantCount();
    const isLive = useIsCallLive();
    const call = useCall();
    const { stream } = useGrpc();
    const { user } = useAuth();

    const [firstParticipant] = useParticipants();


    const handleGoLive = async () => {
        if (!user?.userId) return;
        try {
            await stream.goLive({ userId: user.userId.toString(), callId: call?.id || "" });
            await call.goLive();
        } catch (err) {
            console.error("Failed to go live:", err);
        }
    };

    const handleStopLive = async () => {
        if (!user?.userId) return;
        try {
            await stream.stopLive({ userId: user.userId.toString() });
            await call.stopLive();
        } catch (err) {
            console.error("Failed to stop live:", err);
        }
    };

    return (
        <div className="streamer-ui-container">
            <div className={`streamer-ui-status ${isLive ? "live" : "backstage"}`}>
                {isLive ? `Live: ${participantCount}` : ""}
            </div>
            <div className="streamer-ui-participant">
                {firstParticipant ? (
                    <>
                    <ParticipantView participant={firstParticipant} trackType="screenShareTrack" />
                    <ParticipantView participant={firstParticipant} trackType="videoTrack" />
                    </>
                ) : (
                    <div style={{ color: "#aaa" }}>The host hasn't joined yet</div>
                )}
            </div>
            <div className="streamer-ui-buttons">
                <button
                    className={`streamer-ui-btn ${isLive ? "stop-live" : "go-live"}`}
                    onClick={isLive ? handleStopLive : handleGoLive}
                >
                    {isLive ? "Stop Live" : "Go Live"}
                </button>
                <button
                    className={`streamer-ui-btn ${isCamEnabled ? "enabled" : ""}`}
                    onClick={() => cam.toggle()}
                >
                    {isCamEnabled ? "Disable camera" : "Enable camera"}
                </button>
                <button
                    className={`streamer-ui-btn ${isScreenEnabled ? "enabled" : ""}`}
                    onClick={() => screen.toggle()}
                >
                    {isMicEnabled ? "Start Screen Share" : "Stop Screen Share"}
                </button>
                <button
                    className={`streamer-ui-btn ${isMicEnabled ? "enabled" : ""}`}
                    onClick={() => mic.toggle()}
                >
                    {isMicEnabled ? "Mute Mic" : "Unmute Mic"}
                </button>
            </div>
        </div>
    );
};

export default StreamerUI;
