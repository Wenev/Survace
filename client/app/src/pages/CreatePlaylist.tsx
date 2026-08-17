import { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { useGrpc } from "../context/ClientContext.tsx";
import { useAuth } from "../context/AuthContext.tsx";
import Navbar from "../component/Navbar.tsx";
import MobileNavbar from "../component/MobileNavbar.tsx";
import ErrorToast from "../component/ErrorToast";
import "../style/upload.css";
import type { Video, VideoFeed } from "../generated/dto/video.ts";

// Add PlaylistVideoWithDetails type for local use
type PlaylistVideoWithDetails = {
    playlistId: number;
    videoId: number;
    order: number;
    video?: Video | VideoFeed;
    dragging?: boolean;
};

type Playlist = {
    id: number;
    userId: number;
    title: string;
    videos: PlaylistVideoWithDetails[];
};

export default function CreatePlaylist() {
    const { user } = useAuth();
    const { brainrot } = useGrpc();
    const navigate = useNavigate();

    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [isMobile, setIsMobile] = useState(false);

    // Playlists
    const [playlists, setPlaylists] = useState<any[]>([]);
    const [loadingPlaylists, setLoadingPlaylists] = useState(false);

    // Create modal
    const [showCreateModal, setShowCreateModal] = useState(false);
    const [createTitle, setCreateTitle] = useState("");
    const [createMessage, setCreateMessage] = useState("");
    const [createLoading, setCreateLoading] = useState(false);

    // Edit modal
    const [showEditModal, setShowEditModal] = useState(false);
    const [editPlaylist, setEditPlaylist] = useState<Playlist | null>(null);
    const [editTitle, setEditTitle] = useState("");
    const [editVideos, setEditVideos] = useState<PlaylistVideoWithDetails[]>([]);
    const [suggestedVideos, setSuggestedVideos] = useState<any[]>([]);
    const [editMessage, setEditMessage] = useState("");
    const [editLoading, setEditLoading] = useState(false);
    const [suggestLoading, setSuggestLoading] = useState(false);
    const [deleteLoading, setDeleteLoading] = useState(false);
    const [errorToast, setErrorToast] = useState<string>("");

    useEffect(() => {
        if (!user) {
            navigate("/login");
            return;
        }
        setLoadingPlaylists(true);
        brainrot.getUserPlaylists({ userId: user.userId })
            .then(res => {
                setPlaylists(res.response?.playlists || []);
            })
            .catch(() => setPlaylists([]))
            .finally(() => setLoadingPlaylists(false));
    }, [user, brainrot, navigate]);

    useEffect(() => {
        const checkMobile = () => {
            setIsMobile(window.innerWidth <= 768);
            if (window.innerWidth > 768) setSidebarOpen(false);
        };
        checkMobile();
        window.addEventListener("resize", checkMobile);
        return () => window.removeEventListener("resize", checkMobile);
    }, []);

    const toggleSidebar = () => setSidebarOpen(!sidebarOpen);

    // Helper to refetch playlists
    const refetchPlaylists = async () => {
        setLoadingPlaylists(true);
        try {
            const refreshed = await brainrot.getUserPlaylists({ userId: user.userId });
            setPlaylists(refreshed.response?.playlists || []);
        } catch {
            setPlaylists([]);
            setErrorToast("Failed to fetch playlists");
        }
        setLoadingPlaylists(false);
    };

    // Create playlist modal handlers
    const handleCreatePlaylist = async (e: React.FormEvent) => {
        e.preventDefault();
        setCreateMessage("");
        setCreateLoading(true);
        try {
            await brainrot.createPlaylist({
                userId: user.userId,
                title: createTitle,
            });
            setCreateMessage("Playlist created!");
            setShowCreateModal(false);
            setCreateTitle("");
            await refetchPlaylists();
        } catch {
            setCreateMessage("Failed to create playlist");
            setErrorToast("Failed to create playlist");
        }
        setCreateLoading(false);
    };

    // Edit playlist modal handlers
    const openEditModal = async (playlistId: number) => {
        setEditLoading(true);
        setEditMessage("");
        setSuggestedVideos([]);
        try {
            const res = await brainrot.getPlaylistById({ playlistId });
            const playlist: Playlist = res.response?.playlist;
            setEditPlaylist(playlist);
            setEditTitle(playlist?.title || "");
            let playlistVideos: PlaylistVideoWithDetails[] = playlist?.videos || [];

            // Always fetch video details by videoId for each playlist video
            const videosWithDetails: PlaylistVideoWithDetails[] = await Promise.all(
                playlistVideos.map(async (pv) => {
                    try {
                        const vidRes = await brainrot.getVideoByID({ id: pv.videoId });
                        return { ...pv, video: vidRes.response?.video };
                    } catch {
                        return { ...pv, video: undefined };
                    }
                })
            );
            setEditVideos(videosWithDetails);

            setShowEditModal(true);

            if (!playlist?.videos || playlist.videos.length === 0) {
                setSuggestLoading(true);
                const feedRes = await brainrot.getRandomFeed({ userId: user.userId, limit: 5, offset: 0 });
                setSuggestedVideos(feedRes.response?.videos || []);
                setSuggestLoading(false);
            }
        } catch {
            setEditMessage("Failed to load playlist");
            setErrorToast("Failed to load playlist");
        }
        setEditLoading(false);
    };

    // Add video to playlist handler
    const handleAddVideoToPlaylist = useCallback(async (videoId: number) => {
        if (!editPlaylist) return;
        setEditLoading(true);
        setEditMessage("");
        try {
            await brainrot.addVideoToPlaylist({
                playlistId: editPlaylist.id,
                videoId,
                order: editVideos.length,
            });
            // Refresh playlist videos
            const res = await brainrot.getPlaylistById({ playlistId: editPlaylist.id });
            const playlist: Playlist = res.response?.playlist;
            let playlistVideos: PlaylistVideoWithDetails[] = playlist?.videos || [];
            const videosWithDetails: PlaylistVideoWithDetails[] = await Promise.all(
                playlistVideos.map(async (pv) => {
                    try {
                        const vidRes = await brainrot.getVideoByID({ id: pv.videoId });
                        return { ...pv, video: vidRes.response?.video };
                    } catch {
                        return { ...pv, video: undefined };
                    }
                })
            );
            setEditVideos(videosWithDetails);
            setSuggestedVideos(suggestedVideos.filter(v => v.id !== videoId));
            await refetchPlaylists();
        } catch {
            setEditMessage("Failed to add video");
            setErrorToast("Failed to add video to playlist");
        }
        setEditLoading(false);
    }, [editPlaylist, editVideos, suggestedVideos, brainrot]);

    // Edit title change handler
    const handleEditTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setEditTitle(e.target.value);
    };

    // Drag-and-drop reorder logic
    const handleDragStart = (idx: number) => {
        setEditVideos(videos =>
            videos.map((v, i) => ({ ...v, dragging: i === idx }))
        );
    };
    const handleDragOver = (idx: number) => {
        // No-op, just needed for drop
    };
    const handleDrop = (fromIdx: number, toIdx: number) => {
        setEditVideos(prevVideos =>
            reorder(prevVideos, fromIdx, toIdx).map(v => ({ ...v, dragging: false }))
        );
    };

    // Save edited playlist handler
    const handleEditSave = async () => {
        setEditLoading(true);
        setEditMessage("");
        try {
            // Update playlist title if changed
            if (editPlaylist && editTitle !== editPlaylist.title) {
                await brainrot.updatePlaylistTitle({
                    playlistId: editPlaylist.id,
                    title: editTitle,
                });
            }
            // Reorder videos
            for (let i = 0; i < editVideos.length; i++) {
                const pv = editVideos[i];
                await brainrot.reorderVideo({
                    playlistId: editPlaylist!.id,
                    videoId: pv.videoId,
                    newOrder: i,
                });
            }
            setEditMessage("Playlist updated!");
            setShowEditModal(false);
            await refetchPlaylists();
        } catch {
            setEditMessage("Failed to update playlist");
            setErrorToast("Failed to update playlist");
        }
        setEditLoading(false);
    };

    // Delete playlist handler
    const handleDeletePlaylist = async () => {
        if (!editPlaylist) return;
        setDeleteLoading(true);
        setEditMessage("");
        try {
            await brainrot.deletePlaylist({
                playlistId: editPlaylist.id,
                userId: user.userId,
            });
            setEditMessage("Playlist deleted!");
            setShowEditModal(false);
            setEditPlaylist(null);
            await refetchPlaylists();
        } catch {
            setEditMessage("Failed to delete playlist");
            setErrorToast("Failed to delete playlist");
        }
        setDeleteLoading(false);
    };

    // Helper function to reorder array items
    function reorder<T>(arr: T[], from: number, to: number): T[] {
        const updated = [...arr];
        const [removed] = updated.splice(from, 1);
        updated.splice(to, 0, removed);
        return updated;
    }

    return (
        <div className="app">
            <MobileNavbar onToggleSidebar={toggleSidebar} isOpen={sidebarOpen} />
            <Navbar isOpen={sidebarOpen} />
            <ErrorToast message={errorToast} onClose={() => setErrorToast("")} />
            <div className="container" style={{ maxWidth: 700, margin: "0 auto" }}>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: 24 }}>
                    <h2 style={{ fontWeight: 700, fontSize: 24 }}>Your Playlists</h2>
                    <button className="upload-button" onClick={() => setShowCreateModal(true)}>
                        + Create Playlist
                    </button>
                </div>
                {loadingPlaylists ? (
                    <div>Loading playlists...</div>
                ) : (
                    <div>
                        {playlists.length === 0 ? (
                            <div>No playlists found.</div>
                        ) : (
                            <ul style={{ listStyle: "none", padding: 0 }}>
                                {playlists.map(pl => (
                                    <li
                                        key={pl.id}
                                        style={{
                                            padding: "16px",
                                            marginBottom: "12px",
                                            background: "var(--color-card)",
                                            borderRadius: 8,
                                            boxShadow: "0 2px 8px var(--color-shadow)",
                                            cursor: "pointer",
                                            border: "1px solid var(--color-border)"
                                        }}
                                        onClick={() => openEditModal(pl.id)}
                                    >
                                        <strong>{pl.title}</strong> ({pl.videos?.length || 0} videos)
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                )}

                {/* Create Playlist Modal */}
                {showCreateModal && (
                    <div className="modal-overlay">
                        <div className="modal-content">
                            <h3>Create Playlist</h3>
                            <form onSubmit={handleCreatePlaylist} style={{ display: "flex", flexDirection: "column", gap: "1rem" }}>
                                <label>
                                    Playlist Title
                                    <input
                                        type="text"
                                        value={createTitle}
                                        onChange={e => setCreateTitle(e.target.value)}
                                        required
                                        className="title-input"
                                        placeholder="Enter playlist title"
                                    />
                                </label>
                                <div style={{ display: "flex", gap: "1rem", justifyContent: "flex-end" }}>
                                    <button
                                        className="upload-button"
                                        type="submit"
                                        disabled={createLoading || !createTitle}
                                    >
                                        Create
                                    </button>
                                    <button
                                        type="button"
                                        className="upload-button"
                                        style={{ background: "#888" }}
                                        onClick={() => setShowCreateModal(false)}
                                    >
                                        Cancel
                                    </button>
                                </div>
                                {createMessage && <div className="upload-message">{createMessage}</div>}
                            </form>
                        </div>
                    </div>
                )}

                {/* Edit Playlist Modal */}
                {showEditModal && editPlaylist && (
                    <div className="modal-overlay">
                        <div
                            className="modal-content"
                            style={{
                                maxWidth: 600,
                                background: "var(--color-card)",
                                color: "var(--color-text)",
                                boxShadow: "0 2px 16px var(--modal-overlay-bg, rgba(0,0,0,0.15))"
                            }}
                        >
                            <h3 style={{ color: "var(--color-text)" }}>Edit Playlist</h3>
                            <label style={{ color: "var(--color-text)" }}>
                                Playlist Title
                                <input
                                    type="text"
                                    value={editTitle}
                                    onChange={handleEditTitleChange}
                                    className="title-input"
                                    style={{
                                        background: "var(--color-input-bg)",
                                        color: "var(--color-text)",
                                        border: "1px solid var(--color-input-border)"
                                    }}
                                />
                            </label>
                            <div>
                                <strong style={{ color: "var(--color-text)" }}>Videos (drag to reorder):</strong>
                                <ul style={{ listStyle: "none", padding: 0 }}>
                                    {editVideos.length === 0 ? (
                                        <li style={{ marginBottom: "12px", color: "var(--color-text-secondary)" }}>
                                            No videos in this playlist.
                                        </li>
                                    ) : (
                                        editVideos.map((pv, idx) => (
                                            <li
                                                key={pv.videoId}
                                                draggable
                                                onDragStart={() => handleDragStart(idx)}
                                                onDragOver={e => { e.preventDefault(); handleDragOver(idx); }}
                                                onDrop={e => { e.preventDefault(); handleDrop(editVideos.findIndex(v => v.dragging), idx); }}
                                                style={{
                                                    padding: "10px",
                                                    marginBottom: "8px",
                                                    background: pv.dragging ? "var(--color-input-bg)" : "var(--color-bg-secondary)",
                                                    borderRadius: 6,
                                                    border: "1px solid var(--color-border)",
                                                    cursor: "move",
                                                    display: "flex",
                                                    alignItems: "center",
                                                    gap: "1rem",
                                                    color: "var(--color-text)"
                                                }}
                                            >
                                                {/* Video thumbnail */}
                                                <img
                                                    src={pv.video?.thumbnailUrl || ""}
                                                    alt={pv.video?.title || "Video thumbnail"}
                                                    style={{
                                                        width: 64,
                                                        height: 40,
                                                        objectFit: "cover",
                                                        borderRadius: 4,
                                                        marginRight: 12,
                                                        background: "var(--color-bg-secondary)"
                                                    }}
                                                />
                                                {/* Video title */}
                                                <span style={{ fontWeight: 500, flex: 1, minWidth: 0, whiteSpace: "nowrap", overflow: "hidden", textOverflow: "ellipsis" }}>
                                                    {pv.video?.title || "Loading..."}
                                                </span>
                                            </li>
                                        ))
                                    )}
                                </ul>
                            </div>
                            {/* Suggest videos if playlist is empty */}
                            {editVideos.length === 0 && (
                                <div style={{ marginTop: "16px" }}>
                                    <strong style={{ color: "var(--color-text)" }}>Suggested Videos to Add:</strong>
                                    {suggestLoading ? (
                                        <div style={{ color: "var(--color-text-secondary)" }}>Loading suggestions...</div>
                                    ) : (
                                        <ul style={{ listStyle: "none", padding: 0 }}>
                                            {suggestedVideos.map(video => (
                                                <li key={video.id} style={{
                                                    display: "flex",
                                                    alignItems: "center",
                                                    gap: "1rem",
                                                    marginBottom: "10px",
                                                    background: "var(--color-bg-secondary)",
                                                    borderRadius: 6,
                                                    border: "1px solid var(--color-border)",
                                                    padding: "8px",
                                                    color: "var(--color-text)"
                                                }}>
                                                    <img src={video.thumbnailUrl} alt="" style={{ width: 48, height: 32, objectFit: "cover", borderRadius: 4 }} />
                                                    <span>{video.title}</span>
                                                    <button
                                                        className="upload-button"
                                                        style={{
                                                            padding: "6px 12px",
                                                            fontSize: "12px",
                                                            marginLeft: "auto",
                                                            background: "var(--color-primary)",
                                                            color: "var(--color-bg)"
                                                        }}
                                                        onClick={() => handleAddVideoToPlaylist(video.id)}
                                                        disabled={editLoading}
                                                    >
                                                        Add
                                                    </button>
                                                </li>
                                            ))}
                                            {suggestedVideos.length === 0 && !suggestLoading && (
                                                <li style={{ color: "var(--color-text-secondary)" }}>No suggestions available.</li>
                                            )}
                                        </ul>
                                    )}
                                </div>
                            )}
                            <div style={{ display: "flex", gap: "1rem", justifyContent: "flex-end", marginTop: 16 }}>
                                <button
                                    className="upload-button"
                                    type="button"
                                    onClick={handleEditSave}
                                    disabled={editLoading}
                                    style={{
                                        background: "var(--color-primary)",
                                        color: "var(--color-bg)"
                                    }}
                                >
                                    Save Changes
                                </button>
                                <button
                                    type="button"
                                    className="upload-button"
                                    style={{
                                        background: "#888",
                                        color: "var(--color-bg)"
                                    }}
                                    onClick={() => setShowEditModal(false)}
                                >
                                    Cancel
                                </button>
                                <button
                                    type="button"
                                    className="upload-button"
                                    style={{
                                        background: "var(--color-error)",
                                        color: "var(--color-bg)"
                                    }}
                                    onClick={handleDeletePlaylist}
                                    disabled={deleteLoading}
                                >
                                    Delete Playlist
                                </button>
                            </div>
                            {/* Remove inline error message for edit actions, use toast instead */}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}