import { StreamVideo, StreamCall, useCallStateHooks, ParticipantView, useCall } from "@stream-io/video-react-sdk";
import { useGrpc } from "../context/ClientContext";
import { useEffect, useState } from "react";

interface LiveStream {
  callId: string;
  title?: string;
  hostName?: string;
}

export default function LiveStreamCard({ stream, onWatch }: { stream: LiveStream; onWatch: (callId: string) => void }) {
    const { client } = useGrpc();
    const [call, setCall] = useState<any>(null);

    useEffect(() => {
        if (!client || !stream.callId) return;
        const callInstance = client.call("default", stream.callId);
        callInstance.join({ mode: "replay" })
            .then(() => setCall(callInstance))
            .catch(() => setCall(null));
        return () => {
            callInstance.leave();
        };
    }, [client, stream.callId]);

    return (
        <div className="streamer-ui-container" style={{ maxWidth: 400, margin: 16 }}>
            <div className="streamer-ui-status live">🔴 {stream.title || `Live by ${stream.hostName || stream.callId}`}</div>
            <div className="streamer-ui-participant" style={{ minHeight: 180, justifyContent: 'center', alignItems: 'center' }}>
                {call ? (
                    <StreamVideo client={client}>
                        <StreamCall call={call}>
                            {(() => {
                                const { useParticipants, useScreenShareState } = useCallStateHooks();
                                const [firstParticipant] = useParticipants();
                                const { isEnabled: isScreenEnabled } = useScreenShareState();
                                return isScreenEnabled ? (
                                    <ParticipantView participant={firstParticipant} trackType="screenShareTrack" />
                                ) : firstParticipant ? (
                                    <ParticipantView participant={firstParticipant} trackType="videoTrack" />
                                ) : (
                                    <span style={{ color: '#fe2c55', fontWeight: 600, fontSize: 18 }}>Live Now</span>
                                );
                            })()}
                        </StreamCall>
                    </StreamVideo>
                ) : (
                    <span style={{ color: '#fe2c55', fontWeight: 600, fontSize: 18 }}>Live Now</span>
                )}
            </div>
            <div className="streamer-ui-buttons" style={{ position: 'static', marginTop: 12 }}>
                <button className="streamer-ui-btn go-live" onClick={() => onWatch(stream.callId)}>
                    Watch
                </button>
            </div>
        </div>
    );
}
