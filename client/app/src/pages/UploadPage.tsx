import "../style/upload.css"
import { useEffect, useRef, useState } from "react";
import { useGrpc } from "../context/ClientContext.tsx";
import { UploadVideoRequest } from "../generated/dto/video.ts";
import { useAuth } from "../context/AuthContext.tsx";
import { useNavigate } from "react-router-dom";
import MobileNavbar from "../component/MobileNavbar.tsx";
import Navbar from "../component/Navbar.tsx";
import { UploadCloud } from "lucide-react";
import { Timestamp } from "../generated/google/protobuf/timestamp.ts";

export const UploadPage = () => {
    const navigate = useNavigate();
    const { brainrot } = useGrpc();
    const { user } = useAuth();
    const inputRef = useRef<HTMLInputElement>(null);
    const [file, setFile] = useState<File | null>(null);
    const [message, setMessage] = useState("");
    const [title, setTitle] = useState("");
    const [enableComment, setEnableComment] = useState(true);
    const [visibility, setVisibility] = useState("public");
    const [thumbnail, setThumbnail] = useState<File | null>(null);
    const [hashtags, setHashtags] = useState<string>("");
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [isMobile, setIsMobile] = useState(false);
    const [showDraftModal, setShowDraftModal] = useState(false);
    const [scheduledDate, setScheduledDate] = useState<string>("");

    useEffect(() => {
        if (!user) {
            navigate("/login");
            return;
        }
    }, [user, navigate]);

    useEffect(() => {
        const checkMobile = () => {
            setIsMobile(window.innerWidth <= 768);
            if (window.innerWidth > 768) {
                setSidebarOpen(false);
            }
        };
        checkMobile();
        window.addEventListener("resize", checkMobile);
        return () => window.removeEventListener("resize", checkMobile);
    }, []);

    const toggleSidebar = () => setSidebarOpen(!sidebarOpen);

    const handleSelectClick = () => inputRef.current?.click();

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const selectedFile = e.target.files?.[0];
        if (selectedFile) setFile(selectedFile);
    };

    const handleThumbnailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const selectedFile = e.target.files?.[0];
        if (selectedFile) setThumbnail(selectedFile);
    };

    const handleUpload = async (asDraft: boolean = false, scheduled?: string) => {
        try {
            if (!user) {
                navigate("/login");
                return;
            }
            if (!file) {
                setMessage("No file uploaded");
                return;
            }
            if (!thumbnail) {
                setMessage("No thumbnail");
                return;
            }
            let postDateProto: Timestamp | undefined = undefined;
            if (asDraft && scheduled) {
                const dateObj = new Date(scheduled);
                postDateProto = {
                    seconds: Math.floor(dateObj.getTime() / 1000),
                    nanos: (dateObj.getTime() % 1000) * 1e6,
                };
            }
            const buff = await file.arrayBuffer();
            const thumbBuff = await thumbnail.arrayBuffer();
            const input: UploadVideoRequest = {
                file: new Uint8Array<ArrayBufferLike>(buff),
                userId: user.userId,
                title: title,
                enableComment: enableComment,
                visibility: visibility,
                thumbnail: new Uint8Array<ArrayBufferLike>(thumbBuff),
                isDraft: asDraft,
                postDate: postDateProto,
            };
            const response = await brainrot.uploadVideo(input);
            if (response?.response?.code === 401 || response?.response?.message?.toLowerCase().includes("user not found")) {
                setMessage("Session expired. Please log in again.");
                navigate("/login");
                return;
            }
            setMessage("Upload successful!");
            navigate("/");
        } catch (error: any) {
            if (error?.message?.toLowerCase().includes("user not found")) {
                setMessage("Session expired. Please log in again.");
                navigate("/login");
            } else {
                setMessage("Upload failed. Please try again.");
            }
        }
    };

    const handleDraftSubmit = () => {
        setShowDraftModal(false);
        handleUpload(true, scheduledDate);
    };

    return (
        <div className="app">
            <MobileNavbar onToggleSidebar={toggleSidebar} isOpen={sidebarOpen} />
            <Navbar isOpen={sidebarOpen} />
            <div className="container">
                <div className="upload-card">
                    <div className="upload-content">
                        <div className="upload-icon">
                            <UploadCloud size={32} strokeWidth={2} />
                        </div>
                        <h2 className="upload-title">Select video to upload</h2>
                        <p className="upload-subtitle">or drag and drop</p>
                        <button className="upload-button" onClick={handleSelectClick}>Select Video</button>
                        <input
                            ref={inputRef}
                            type="file"
                            accept="video/mp4"
                            onChange={handleFileChange}
                            style={{ display: 'none' }}
                        />
                        <div className="video-title-form">
                            <label htmlFor="title" className="title-label">Video Title</label>
                            <input
                                id="title"
                                type="text"
                                className="title-input"
                                placeholder="Enter your video title"
                                value={title}
                                onChange={(e) => setTitle(e.target.value)}
                            />
                        </div>
                        <div className="video-settings-form">
                            <label className="settings-label">
                                <input
                                    type="checkbox"
                                    checked={enableComment}
                                    onChange={e => setEnableComment(e.target.checked)}
                                />
                                Enable Comments
                            </label>
                            <label className="settings-label">
                                Visibility:
                                <select
                                    value={visibility}
                                    onChange={e => setVisibility(e.target.value)}
                                >
                                    <option value="public">Public</option>
                                    <option value="private">Private</option>
                                    <option value="unlisted">Unlisted</option>
                                </select>
                            </label>
                        </div>
                        <div className="video-thumbnail-form">
                            <label htmlFor="thumbnail" className="title-label">Upload Thumbnail</label>
                            <input
                                id="thumbnail"
                                type="file"
                                accept="image/*"
                                onChange={handleThumbnailChange}
                            />
                        </div>
                        <div className="video-hashtag-form">
                            <label htmlFor="hashtags" className="title-label">Hashtags (optional, e.g. #fun #music)</label>
                            <input
                                id="hashtags"
                                type="text"
                                className="title-input"
                                placeholder="#hashtag1 #hashtag2"
                                value={hashtags}
                                onChange={e => setHashtags(e.target.value)}
                            />
                        </div>
                        <div style={{ display: "flex", gap: "1rem", marginTop: "1rem" }}>
                            <button
                                type="button"
                                className="upload-button"
                                onClick={() => handleUpload(false)}
                            >
                                Upload Video
                            </button>
                            <button
                                type="button"
                                className="upload-button"
                                style={{ background: "#888" }}
                                onClick={() => setShowDraftModal(true)}
                            >
                                Save as Draft
                            </button>
                        </div>
                        {showDraftModal && (
                            <div className="modal-overlay">
                                <div className="modal-content">
                                    <h3>Schedule Draft</h3>
                                    <label>
                                        Scheduled Post Date:
                                        <input
                                            type="datetime-local"
                                            value={scheduledDate}
                                            onChange={e => setScheduledDate(e.target.value)}
                                            className="title-input"
                                        />
                                    </label>
                                    <div style={{ marginTop: "1rem", display: "flex", gap: "1rem" }}>
                                        <button
                                            type="button"
                                            onClick={handleDraftSubmit}
                                            className="upload-button"
                                            disabled={!scheduledDate}
                                        >
                                            Save Draft
                                        </button>
                                        <button
                                            type="button"
                                            onClick={() => setShowDraftModal(false)}
                                            className="upload-button"
                                            style={{ background: "#888" }}
                                        >
                                            Cancel
                                        </button>
                                    </div>
                                </div>
                            </div>
                        )}
                        {message && (
                            <div className="upload-message">{message}</div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};
