import { ParticipantView, useCallStateHooks } from "@stream-io/video-react-sdk";
import "../style/component/streamer-ui.css";

const WatcherUI = () => {
    const {
        useParticipantCount,
        useParticipants,
        useScreenShareState
    } = useCallStateHooks();

    const participantCount = useParticipantCount();
    const { isEnabled: isScreenEnabled } = useScreenShareState();
    const [firstParticipant] = useParticipants();

    return (
        <div className="streamer-ui-container">
            <div className="streamer-ui-status live">
                {`Live: ${participantCount}`}
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
        </div>
    );
};

export default WatcherUI;
