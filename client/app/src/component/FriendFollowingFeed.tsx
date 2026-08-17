import { useState, useEffect, useCallback, useRef } from "react";
import { VideoProp } from "./VideoProp";
import defaultAvatar from "../assets/default.jpg";
import "../style/component/friend-following-feed.css";
import "../style/component/video-prop.css";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import VideoPreview from "./VideoPreview";
import type { VideoFeed } from "../generated/dto/video";
import { ArrowLeft, ChevronUp, ChevronDown } from "lucide-react";
import { useNavigate } from "react-router-dom";

type Mode = "friends" | "following";

interface FriendOrFollowing {
    id: number;
    username: string;
    avatarUrl?: string;
    isLive?: boolean;
}

interface Video {
    id: number;
    userId: number;
    url: string;
    title: string;
    createdAt: string;
    likeCount?: number;
    commentCount?: number;
    isLiked?: boolean;
    thumbnailUrl?: string;
}

interface FriendFollowingFeedProps {
    mode: Mode;
}

interface FriendData {
    friendId: number;
}
interface FollowingData {
    followeeId: number;
}

export default function FriendFollowingFeed({ mode }: FriendFollowingFeedProps) {
    const { user } = useAuth();
    const { social, galactus, brainrot, stream } = useGrpc();
    const navigate = useNavigate();

    const [list, setList] = useState<FriendOrFollowing[]>([]);
    const [selectedId, setSelectedId] = useState<number | null>(null);
    const [videos, setVideos] = useState<VideoFeed[]>([]);
    const [cachedVideos, setCachedVideos] = useState<VideoFeed[]>([]);
    const [loading, setLoading] = useState(false);
    const [loadingVideos, setLoadingVideos] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [hasMore, setHasMore] = useState(true);
    const [page, setPage] = useState(0);
    const [cycleIndex, setCycleIndex] = useState(0);
    const [currentIndex, setCurrentIndex] = useState(0);
    const loadingRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (!user) return;
        setLoading(true);
        const fetchList = async () => {
            try {
                let ids: number[] = [];
                if (mode === "friends") {
                    const res = await social.listFriends({ userId: user.userId });
                    ids = (res.response?.data as FriendData[] || []).map((f) => f.friendId);
                } else {
                    const res = await social.listFollowing({ userId: user.userId });
                    ids = (res.response?.data as FollowingData[] || []).map((f) => f.followeeId);
                }
                const users: FriendOrFollowing[] = await Promise.all(
                    ids.map(async (id) => {
                        try {
                            const ures = await galactus.findByUserId({ id });
                            const u = ures.response?.user;
                            return {
                                id: u.userId,
                                username: u.username,
                                avatarUrl: u.avatarUrl || defaultAvatar,
                                isLive: false,
                            };
                        } catch {
                            return {
                                id,
                                username: "Unknown",
                                avatarUrl: defaultAvatar,
                                isLive: false,
                            };
                        }
                    })
                );
                let liveMap: Record<number, boolean> = {};
                try {
                    const liveRes = await stream.listLiveStreams({});
                    const liveIds = (liveRes.response?.streams || []).map((s: { userId: string }) => Number(s.userId));
                    liveMap = Object.fromEntries(liveIds.map((id) => [id, true]));
                } catch (e) { }
                setList(users.map(u => ({ ...u, isLive: !!liveMap[u.id] })));
                setSelectedId(users.length > 0 ? users[0].id : null);
            } catch (e: any) {
                setError(e?.message || "Failed to load list");
            }
            setLoading(false);
        };
        fetchList();
    }, [user, mode, galactus, social, stream]);

    const fetchFeed = useCallback(async (reset = false) => {
        if (!selectedId) return;
        setLoadingVideos(true);
        try {
            let res;
            if (mode === "friends") {
                res = await brainrot.friendVideo({ userId: user.userId, friendId: selectedId, limit: 10, offset: reset ? 0 : page * 10 });
            } else {
                res = await brainrot.followingVideo({ userId: user.userId, followeeId: selectedId, limit: 10, offset: reset ? 0 : page * 10 });
            }
            const newVideos = res?.response?.videos || [];
            if (reset) {
                setVideos(newVideos);
                setCachedVideos(newVideos);
                setCycleIndex(0);
                setHasMore(true);
            } else {
                if (newVideos.length === 0) {
                    if (cachedVideos.length > 0) {
                        let nextBatch = [];
                        if (cycleIndex + 10 <= cachedVideos.length) {
                            nextBatch = cachedVideos.slice(cycleIndex, cycleIndex + 10);
                        } else {
                            nextBatch = [
                                ...cachedVideos.slice(cycleIndex),
                                ...cachedVideos.slice(0, (cycleIndex + 10) % cachedVideos.length)
                            ];
                        }
                        setVideos(prev => [...prev, ...nextBatch]);
                        setCycleIndex((prev) => (prev + 10) % cachedVideos.length);
                        setHasMore(true);
                    } else {
                        setHasMore(false);
                    }
                    return;
                } else {
                    setVideos(prev => [...prev, ...newVideos]);
                    setCachedVideos(prev => [...prev, ...newVideos]);
                    setHasMore(true);
                }
            }
        } catch (e: any) {
            setError(e?.message || "Failed to load videos");
        } finally {
            setLoadingVideos(false);
        }
    }, [selectedId, user, mode, brainrot, page, cachedVideos, cycleIndex]);

    useEffect(() => {
        setPage(0);
        fetchFeed(true);
    }, [selectedId, mode]);

    useEffect(() => {
        if (!hasMore || loadingVideos) return;
        const observer = new window.IntersectionObserver(entries => {
            if (entries[0].isIntersecting) {
                setPage(prev => prev + 1);
            }
        }, { threshold: 1 });
        if (loadingRef.current) observer.observe(loadingRef.current);
        return () => observer.disconnect();
    }, [hasMore, loadingVideos]);

    useEffect(() => {
        if (page === 0) return;
        fetchFeed();
    }, [page]);

    useEffect(() => {
        if (videos.length === 0) return;
        const videoId = videos[currentIndex]?.id;
        if (videoId !== undefined) {
            const el = document.querySelector(`[data-video-id='${videoId}']`);
            if (el) {
                el.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }
    }, [currentIndex, videos]);

    const handleUp = () => {
        setCurrentIndex((prev) => {
            const newIdx = Math.max(prev - 1, 0);
            scrollToVideo(newIdx);
            return newIdx;
        });
    };
    const handleDown = () => {
        setCurrentIndex((prev) => {
            const newIdx = Math.min(prev + 1, videos.length - 1);
            scrollToVideo(newIdx);
            return newIdx;
        });
    };

    const scrollToVideo = (idx: number) => {
        const videoId = videos[idx]?.id;
        if (videoId !== undefined) {
            const el = document.querySelector(`[data-video-id='${videoId}']`);
            if (el) {
                el.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }
    };

    useEffect(() => {
        const feedEl = document.querySelector('.vertical-feed');
        if (!feedEl) return;
        const handleScroll = () => {
            const videoEls = Array.from(feedEl.querySelectorAll('[data-video-id]'));
            let closestIdx = 0;
            let minDist = Infinity;
            videoEls.forEach((el, idx) => {
                const rect = el.getBoundingClientRect();
                const dist = Math.abs(rect.top - window.innerHeight / 2);
                if (dist < minDist) {
                    minDist = dist;
                    closestIdx = idx;
                }
            });
            setCurrentIndex(closestIdx);
        };
        feedEl.addEventListener('scroll', handleScroll, { passive: true });
        return () => feedEl.removeEventListener('scroll', handleScroll);
    }, [videos.length]);

    useEffect(() => {
        setCurrentIndex(0);
    }, [videos]);

    const handleLikeUpdate = (videoId: number, delta: number, liked: boolean) => {
        setVideos(videos =>
            videos.map(v =>
                v.id === videoId
                    ? { ...v, likeCount: (v.likeCount ?? 0) + delta, isLiked: liked }
                    : v
            )
        );
    };
    const wrapperClass = "friend-following-feed";

    return (
        <div className={wrapperClass} style={{ flexDirection: "column", position: "relative" }}>
            <div
                className="side-scroll-list"
                style={{
                    position: "sticky",
                    top: 0,
                    zIndex: 2,
                    background: "var(--color-bg, #fff)",
                    boxShadow: "0 2px 8px rgba(0,0,0,0.03)",
                    paddingBottom: 8,
                    marginBottom: 8,
                }}
            >
                {loading ? (
                    <div style={{ padding: 16 }}>Loading...</div>
                ) : (
                    <div className="side-scroll-inner" style={{ display: "flex", flexDirection: "row", gap: 16, overflowX: "auto" }}>
                        {list.map((item) => (
                            <div
                                key={item.id}
                                className={`side-scroll-item${selectedId === item.id ? " selected" : ""}`}
                                onClick={() => {
                                    navigate(`/profile/${item.id}`);
                                }}
                                style={{
                                    cursor: "pointer",
                                    display: "flex",
                                    flexDirection: "column",
                                    alignItems: "center",
                                    minWidth: 80,
                                }}
                            >
                                <div style={{ position: "relative" }}>
                                    <img
                                        src={item.avatarUrl || defaultAvatar}
                                        alt={item.username}
                                        style={{
                                            width: 54,
                                            height: 54,
                                            borderRadius: "50%",
                                            border: selectedId === item.id ? "2px solid #fe2c55" : "2px solid #eee",
                                            objectFit: "cover",
                                            boxShadow: item.isLive ? "0 0 0 2px #fe2c55" : undefined,
                                        }}
                                    />
                                    {item.isLive && (
                                        <span
                                            style={{
                                                position: "absolute",
                                                bottom: -6,
                                                left: "50%",
                                                transform: "translateX(-50%)",
                                                background: "#fe2c55",
                                                color: "#fff",
                                                borderRadius: 8,
                                                fontSize: 11,
                                                padding: "1px 6px",
                                                fontWeight: 600,
                                                letterSpacing: 0.5,
                                                boxShadow: "0 1px 4px rgba(0,0,0,0.08)",
                                            }}
                                        >
                                            LIVE
                                        </span>
                                    )}
                                </div>
                                <div style={{ fontSize: 13, marginTop: 4, maxWidth: 70, overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap" }}>
                                    {item.username}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </div>
            <div className="vertical-feed">
                {loadingVideos ? (
                    <div style={{ padding: 24 }}>Loading videos...</div>
                ) : videos.length === 0 ? (
                    <div style={{ padding: 24, color: "#aaa" }}>No videos found.</div>
                ) : (
                    <div className="video-list">
                        {videos.map((video, idx) => (
                            <VideoProp
                                key={video.id}
                                video={{
                                    ...video,
                                    isPlaying: idx === currentIndex,
                                    thumbnailUrl: (video as any).thumbnailUrl || undefined,
                                    isLiked: video.isLiked ?? false
                                }}
                                onTogglePlay={() => {}}
                                onLikeUpdate={handleLikeUpdate}
                                data-video-id={video.id}
                            />
                        ))}
                        <div ref={loadingRef} className="loading-trigger" />
                    </div>
                )}
                {error && (
                    <div style={{ color: "red", marginTop: 12, textAlign: "center" }}>{error}</div>
                )}
                <div className="traverse-btn-group">
                    <button
                        className="traverse-btn up"
                        onClick={handleUp}
                        disabled={currentIndex === 0}
                    >
                        <ChevronUp size={28} />
                    </button>
                    <button
                        className="traverse-btn down"
                        onClick={handleDown}
                        disabled={currentIndex === videos.length - 1}
                    >
                        <ChevronDown size={28} />
                    </button>
                </div>
            </div>
        </div>
    );
}