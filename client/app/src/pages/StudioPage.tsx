import { useEffect, useState, useRef } from "react";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";
import Navbar from "../component/Navbar";
import MobileNavbar from "../component/MobileNavbar";
import "../style/upload.css";
import type {VideoFeed} from "../generated/dto/video.ts";

type SortKey = "likeCount" | "viewCount" | "commentCount";
type SortOrder = "asc" | "desc";
type StatusFilter = "all" | "published" | "draft";

export default function StudioPage() {
    const { brainrot } = useGrpc();
    const { user } = useAuth();
    const navigate = useNavigate();
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [isMobile, setIsMobile] = useState(false);
    const [videos, setVideos] = useState<VideoFeed[]>([]);
    const [sortKey, setSortKey] = useState<SortKey>("likeCount");
    const [sortOrder, setSortOrder] = useState<SortOrder>("desc");
    const [loading, setLoading] = useState(false);
    const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");
    
    // Add search state variables
    const [searchQuery, setSearchQuery] = useState<string>("");
    const [debouncedSearchQuery, setDebouncedSearchQuery] = useState<string>("");

    const [editModalOpen, setEditModalOpen] = useState(false);
    const [editVideo, setEditVideo] = useState<any | null>(null);
    const [editTitle, setEditTitle] = useState("");
    const [editEnableComment, setEditEnableComment] = useState(true);
    const [editVisibility, setEditVisibility] = useState("public");
    const [editThumbnail, setEditThumbnail] = useState<File | null>(null);
    const [editThumbnailUrl, setEditThumbnailUrl] = useState("");
    const [editIsDraft, setEditIsDraft] = useState(false);
    const [editPostedAt, setEditPostedAt] = useState<string>("");
    const [editMessage, setEditMessage] = useState("");
    const editThumbnailInputRef = useRef<HTMLInputElement>(null);

    // Add delete video modal state
    const [deleteModalOpen, setDeleteModalOpen] = useState(false);
    const [deleteLoading, setDeleteLoading] = useState(false);

    useEffect(() => {
        if (!user) {
            navigate("/login");
            return;
        }
    }, [user, navigate]);

    useEffect(() => {
        const checkMobile = () => {
            setIsMobile(window.innerWidth <= 768);
            if (window.innerWidth > 768) setSidebarOpen(false);
        };
        checkMobile();
        window.addEventListener("resize", checkMobile);
        return () => window.removeEventListener("resize", checkMobile);
    }, []);

    const fetchVideos = async () => {
        setLoading(true);
        try {
            const res = await brainrot.getUserVideoAndDraft({ userId: user.userId });
            const feeds = res.response.videos;
            setVideos(feeds);
        } catch {
            setVideos([]);
        }
        setLoading(false);
    };

    useEffect(() => {
        if (user) fetchVideos();
        // eslint-disable-next-line
    }, [user]);

    // Add debounce effect for search
    useEffect(() => {
        const timer = setTimeout(() => {
            setDebouncedSearchQuery(searchQuery);
        }, 300); // 300ms debounce delay

        return () => {
            clearTimeout(timer);
        };
    }, [searchQuery]);

    const filteredVideos = videos.filter(video => {
        // First filter by status
        if (statusFilter !== "all" && 
            ((statusFilter === "draft" && !video.isDraft) || 
             (statusFilter === "published" && video.isDraft))) {
            return false;
        }
        
        // Then filter by search query
        if (debouncedSearchQuery) {
            return video.title.toLowerCase().includes(debouncedSearchQuery.toLowerCase());
        }
        
        return true;
    });

    const sortedVideos = [...filteredVideos].sort((a, b) => {
        const aVal = a[sortKey] ?? 0;
        const bVal = b[sortKey] ?? 0;
        return sortOrder === "asc" ? aVal - bVal : bVal - aVal;
    });

    const toggleSidebar = () => setSidebarOpen(!sidebarOpen);

    // Open edit modal
    const handleEditClick = (video: any) => {
        setEditVideo(video);
        setEditTitle(video.title);
        setEditEnableComment(video.enableComment);
        setEditVisibility(video.visibility);
        setEditThumbnailUrl(video.thumbnailUrl);
        setEditIsDraft(video.isDraft);
        let postedAtStr = "";
        if (video.postedAt?.seconds !== undefined && video.postedAt?.seconds !== null) {
            let seconds = video.postedAt.seconds;
            if (typeof seconds === "bigint") {
                seconds = Number(seconds);
            } else if (typeof seconds === "string") {
                seconds = Number(seconds);
            }
            postedAtStr = new Date(seconds * 1000).toISOString().slice(0, 16);
        }
        setEditPostedAt(postedAtStr);
        setEditMessage("");
        setEditModalOpen(true);
    };

    // Delete video handler
    const handleDeleteVideo = async () => {
        if (!editVideo) return;
        setDeleteLoading(true);
        try {
            await brainrot.deleteVideoByID({ videoId: editVideo.id });
            setDeleteLoading(false);
            setDeleteModalOpen(false);
            setEditModalOpen(false);
            fetchVideos();
        } catch {
            setDeleteLoading(false);
        }
    };

    const handleEditThumbnailChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const selectedFile = e.target.files?.[0];
        if (selectedFile) {
            setEditThumbnail(selectedFile);
            setEditThumbnailUrl(URL.createObjectURL(selectedFile));
        }
    };

    const handleEditSubmit = async () => {
        if (!editVideo) return;
        setEditMessage("");
        let thumbnailBuff: Uint8Array | undefined = undefined;
        if (editThumbnail) {
            const buff = await editThumbnail.arrayBuffer();
            thumbnailBuff = new Uint8Array(buff);
        }
        // Prepare postedAt from string (like UploadPage)
        let postedAtProto: any = undefined;
        if (editPostedAt) {
            const dateObj = new Date(editPostedAt);
            postedAtProto = {
                seconds: Math.floor(dateObj.getTime() / 1000),
                nanos: (dateObj.getTime() % 1000) * 1e6,
            };
        }
        try {
            const res = await brainrot.updateVideo({
                id: editVideo.id,
                userId: user.userId,
                title: editTitle,
                url: editVideo.url,
                objectName: editVideo.objectName,
                enableComment: editEnableComment,
                visibility: editVisibility,
                thumbnail: thumbnailBuff, // send thumbnail as bytes
                postedAt: postedAtProto,
                isDraft: editIsDraft,
            });
            if (res?.response?.code === 0) {
                setEditMessage("Video updated!");
                setEditModalOpen(false);
                fetchVideos();
            } else {
                setEditMessage(res?.response?.message || "Failed to update video.");
            }
        } catch {
            setEditMessage("Failed to update video.");
        }
    };

    return (
        <div className="app">
            <MobileNavbar onToggleSidebar={toggleSidebar} isOpen={sidebarOpen} />
            <Navbar isOpen={sidebarOpen} />
            <div className="container" style={{ padding: "32px 16px", maxWidth: 1000 }}>
                <div className="upload-card" style={{ maxWidth: 900, padding: "32px 24px" }}>
                    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center", marginBottom: "1.5rem" }}>
                        <h2 style={{ fontWeight: 700, fontSize: 28 }}>Studio</h2>
                        <button 
                            className="upload-button" 
                            onClick={() => navigate('/upload')}
                            style={{
                                padding: "10px 16px",
                                fontSize: 15,
                                display: "flex",
                                alignItems: "center",
                                gap: 8
                            }}
                        >
                            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                                <path d="M12 5V19M5 12H19" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                            </svg>
                            Upload New Video
                        </button>
                    </div>
                    
                    {/* Add search input */}
                    <div style={{ marginBottom: "1.5rem" }}>
                        <div style={{ position: "relative" }}>
                            <input
                                type="text"
                                placeholder="Search videos by title..."
                                value={searchQuery}
                                onChange={(e) => setSearchQuery(e.target.value)}
                                style={{
                                    width: "100%",
                                    padding: "10px 16px 10px 40px",
                                    borderRadius: "8px",
                                    border: "1px solid var(--color-border)",
                                    fontSize: "15px",
                                    background: "var(--color-input-bg, #f7f7f7)",
                                    color: "var(--color-text)"
                                }}
                            />
                            <svg
                                width="18"
                                height="18"
                                viewBox="0 0 24 24"
                                fill="none"
                                xmlns="http://www.w3.org/2000/svg"
                                style={{
                                    position: "absolute",
                                    left: "14px",
                                    top: "50%",
                                    transform: "translateY(-50%)",
                                    color: "var(--color-text-secondary)"
                                }}
                            >
                                <path
                                    d="M21 21L15 15M17 10C17 13.866 13.866 17 10 17C6.13401 17 3 13.866 3 10C3 6.13401 6.13401 3 10 3C13.866 3 17 6.13401 17 10Z"
                                    stroke="currentColor"
                                    strokeWidth="2"
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                />
                            </svg>
                            {searchQuery && (
                                <button
                                    onClick={() => setSearchQuery("")}
                                    style={{
                                        position: "absolute",
                                        right: "14px",
                                        top: "50%",
                                        transform: "translateY(-50%)",
                                        background: "none",
                                        border: "none",
                                        cursor: "pointer",
                                        color: "var(--color-text-secondary)"
                                    }}
                                    aria-label="Clear search"
                                >
                                    ×
                                </button>
                            )}
                        </div>
                    </div>
                    
                    <div style={{
                        display: "flex",
                        gap: "2rem",
                        alignItems: "center",
                        marginBottom: "2rem",
                        flexWrap: "wrap"
                    }}>
                        <div style={{ display: "flex", gap: "1rem", alignItems: "center" }}>
                            <label style={{ fontWeight: 500 }}>Sort by:</label>
                            <select value={sortKey} onChange={e => setSortKey(e.target.value as SortKey)} style={{ padding: "8px", borderRadius: 6 }}>
                                <option value="likeCount">Likes</option>
                                <option value="viewCount">Views</option>
                                <option value="commentCount">Comments</option>
                            </select>
                            <select value={sortOrder} onChange={e => setSortOrder(e.target.value as SortOrder)} style={{ padding: "8px", borderRadius: 6 }}>
                                <option value="desc">Descending</option>
                                <option value="asc">Ascending</option>
                            </select>
                        </div>
                        <div style={{ display: "flex", gap: "1rem", alignItems: "center" }}>
                            <label style={{ fontWeight: 500 }}>Filter:</label>
                            <select value={statusFilter} onChange={e => setStatusFilter(e.target.value as StatusFilter)} style={{ padding: "8px", borderRadius: 6 }}>
                                <option value="all">All</option>
                                <option value="published">Published</option>
                                <option value="draft">Draft</option>
                            </select>
                        </div>
                    </div>
                    
                    {loading ? (
                        <div style={{ padding: "2rem", textAlign: "center" }}>Loading...</div>
                    ) : (
                        <>
                            {filteredVideos.length > 0 ? (
                                <table className="studio-table" style={{ width: "100%", borderCollapse: "collapse", marginBottom: "2rem" }}>
                                    <thead>
                                        <tr style={{ background: "var(--color-bg-secondary)" }}>
                                            <th style={{ padding: "12px" }}>Thumbnail</th>
                                            <th style={{ padding: "12px" }}>Title</th>
                                            <th style={{ padding: "12px" }}>Status</th>
                                            <th style={{ padding: "12px" }}>Likes</th>
                                            <th style={{ padding: "12px" }}>Views</th>
                                            <th style={{ padding: "12px" }}>Comments</th>
                                            <th style={{ padding: "12px" }}>Posted At</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {sortedVideos.map(video => {
                                            let postedAtDisplay = "-";
                                            if (video.postedAt?.seconds !== undefined && video.postedAt?.seconds !== null) {
                                                let seconds = video.postedAt.seconds;
                                                if (typeof seconds === "bigint") {
                                                    seconds = Number(seconds);
                                                } else if (typeof seconds === "string") {
                                                    seconds = Number(seconds);
                                                }
                                                postedAtDisplay = new Date(seconds * 1000).toLocaleString();
                                            }
                                            return (
                                                <tr key={video.id} style={{ cursor: "pointer", borderBottom: "1px solid var(--color-border)" }} onClick={() => handleEditClick(video)}>
                                                    <td style={{ padding: "10px" }}>
                                                        <img 
                                                            src={video.thumbnailUrl} 
                                                            alt="thumb" 
                                                            style={{ width: 80, borderRadius: 8 }} 
                                                            loading="lazy"
                                                        />
                                                    </td>
                                                    <td style={{ padding: "10px", fontWeight: 500 }}>{video.title}</td>
                                                    <td style={{ padding: "10px" }}>{video.isDraft ? "Draft" : "Published"}</td>
                                                    <td style={{ padding: "10px" }}>{video.likeCount ?? 0}</td>
                                                    <td style={{ padding: "10px" }}>{video.viewCount ?? 0}</td>
                                                    <td style={{ padding: "10px" }}>{video.commentCount ?? 0}</td>
                                                    <td style={{ padding: "10px" }}>{postedAtDisplay}</td>
                                                </tr>
                                            );
                                        })}
                                    </tbody>
                                </table>
                            ) : (
                                <div style={{ marginTop: "2rem", textAlign: "center", color: "var(--color-text-secondary)" }}>
                                    {searchQuery ? `No videos found matching "${searchQuery}"` : "No videos found."}
                                </div>
                            )}
                        </>
                    )}
                </div>
            </div>
            {/* Move modal outside of .container to avoid CSS stacking issues */}
            {editModalOpen && (
                <div className="modal-overlay" style={{
                    position: "fixed",
                    top: 0, left: 0, right: 0, bottom: 0,
                    background: "var(--modal-overlay-bg, rgba(0,0,0,0.5))",
                    zIndex: 1000,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center"
                }}>
                    <div className="modal-content" style={{
                        maxWidth: 400,
                        background: "var(--color-card, #fff)",
                        borderRadius: 12,
                        padding: 24,
                        boxShadow: "0 2px 16px rgba(0,0,0,0.15)",
                        color: "var(--color-text, #222)",
                        display: "flex",
                        flexDirection: "column",
                        gap: 16,
                        transition: "background 0.2s, color 0.2s"
                    }}>
                        <h3 style={{ color: "var(--color-text, #222)" }}>Edit Video</h3>
                        <label>
                            Title:
                            <input
                                type="text"
                                value={editTitle}
                                onChange={e => setEditTitle(e.target.value)}
                                className="title-input"
                                style={{ background: "var(--color-input-bg, #f7f7f7)", color: "var(--color-text, #222)", border: "1px solid var(--color-input-border, #ddd)" }}
                            />
                        </label>
                        <label>
                            Enable Comments:
                            <input
                                type="checkbox"
                                checked={editEnableComment}
                                onChange={e => setEditEnableComment(e.target.checked)}
                                style={{ marginLeft: 8 }}
                            />
                        </label>
                        <label>
                            Visibility:
                            <select
                                value={editVisibility}
                                onChange={e => setEditVisibility(e.target.value)}
                                className="title-input"
                                style={{ background: "var(--color-input-bg, #f7f7f7)", color: "var(--color-text, #222)", border: "1px solid var(--color-input-border, #ddd)" }}
                            >
                                <option value="public">Public</option>
                                <option value="private">Private</option>
                                <option value="unlisted">Unlisted</option>
                            </select>
                        </label>
                        <label>
                            Thumbnail:
                            <input
                                type="file"
                                accept="image/*"
                                ref={editThumbnailInputRef}
                                onChange={handleEditThumbnailChange}
                                style={{ marginTop: 8 }}
                            />
                            {editThumbnailUrl && (
                                <img src={editThumbnailUrl} alt="thumb" style={{ width: 80, borderRadius: 8, marginTop: 8 }} />
                            )}
                        </label>
                        <label>
                            Status:
                            <select
                                value={editIsDraft ? "draft" : "published"}
                                onChange={e => setEditIsDraft(e.target.value === "draft")}
                                className="title-input"
                                style={{ background: "var(--color-input-bg, #f7f7f7)", color: "var(--color-text, #222)", border: "1px solid var(--color-input-border, #ddd)" }}
                            >
                                <option value="published">Published</option>
                                <option value="draft">Draft</option>
                            </select>
                        </label>
                        <label>
                            Scheduled Post Date:
                            <input
                                type="datetime-local"
                                value={editPostedAt}
                                onChange={e => setEditPostedAt(e.target.value)}
                                className="title-input"
                                style={{
                                    background: editIsDraft ? "var(--color-input-bg, #f7f7f7)" : "#eee",
                                    color: editIsDraft ? "var(--color-text, #222)" : "#aaa",
                                    border: "1px solid var(--color-input-border, #ddd)"
                                }}
                                disabled={!editIsDraft}
                            />
                        </label>
                        {editMessage && (
                            <div style={{ color: "var(--color-error, red)", marginTop: 8 }}>{editMessage}</div>
                        )}
                        <div style={{ display: "flex", gap: "1rem", marginTop: "1rem" }}>
                            <button className="upload-button" onClick={handleEditSubmit}>Save</button>
                            <button className="upload-button" style={{ background: "#888" }} onClick={() => setEditModalOpen(false)}>Cancel</button>
                            <button
                                className="upload-button"
                                style={{ background: "var(--color-error)", color: "var(--color-bg)" }}
                                type="button"
                                onClick={() => setDeleteModalOpen(true)}
                                disabled={deleteLoading}
                            >
                                Delete Video
                            </button>
                        </div>
                    </div>
                </div>
            )}
            {/* Delete Video Modal */}
            {deleteModalOpen && (
                <div className="modal-overlay" style={{
                    position: "fixed",
                    top: 0, left: 0, right: 0, bottom: 0,
                    background: "var(--modal-overlay-bg, rgba(0,0,0,0.5))",
                    zIndex: 1100,
                    display: "flex",
                    alignItems: "center",
                    justifyContent: "center"
                }}>
                    <div className="modal-content" style={{
                        maxWidth: 350,
                        background: "var(--color-card, #fff)",
                        borderRadius: 12,
                        padding: 24,
                        boxShadow: "0 2px 16px rgba(0,0,0,0.15)",
                        color: "var(--color-text, #222)",
                        display: "flex",
                        flexDirection: "column",
                        gap: 16,
                        transition: "background 0.2s, color 0.2s"
                    }}>
                        <h3 style={{ color: "var(--color-error, #fe2c55)" }}>Delete Video</h3>
                        <div>Are you sure you want to delete this video?</div>
                        <div style={{ display: "flex", gap: "1rem", marginTop: "1rem" }}>
                            <button
                                className="upload-button"
                                style={{ background: "var(--color-error)", color: "var(--color-bg)" }}
                                onClick={handleDeleteVideo}
                                disabled={deleteLoading}
                            >
                                {deleteLoading ? "Deleting..." : "Delete"}
                            </button>
                            <button
                                className="upload-button"
                                style={{ background: "#888" }}
                                onClick={() => setDeleteModalOpen(false)}
                            >
                                Cancel
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}