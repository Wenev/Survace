"use client"

import { useState, useEffect, useCallback, useRef } from "react"
import { VideoProp } from "./VideoProp"
import { VideoSkeleton } from "./LoadingSkeleton"
import "../style/component/feed.css"
import { useAuth } from "../context/AuthContext"
import { useGrpc } from "../context/ClientContext"
import type {Video, VideoFeed} from "../generated/dto/video.ts";
import { ChevronUp, ChevronDown } from "lucide-react";

interface VideoFeedProps {
    activeVideoId: string | null
    onVideoChange: (videoId: string) => void
    searchQuery?: string
    setError?: (msg: string) => void
}

function useDebounce<T>(value: T, delay: number): T {
    const [debouncedValue, setDebouncedValue] = useState(value);
    useEffect(() => {
        const handler = setTimeout(() => setDebouncedValue(value), delay);
        return () => clearTimeout(handler);
    }, [value, delay]);
    return debouncedValue;
}

export default function VideoFeed({ activeVideoId, onVideoChange, searchQuery = "", setError }: VideoFeedProps) {
    const { user } = useAuth()
    const { brainrot, useVideoQuery } = useGrpc()
    const [videos, setVideos] = useState<VideoFeed[]>([])
    const [cachedVideos, setCachedVideos] = useState<VideoFeed[]>([])
    const [loading, setLoading] = useState(false)
    const loadingRef = useRef<HTMLDivElement>(null)
    const [page, setPage] = useState(1)
    const [hasMore, setHasMore] = useState(true)
    const [cycleIndex, setCycleIndex] = useState(0)
    const [currentIndex, setCurrentIndex] = useState(0);
    const debouncedSearch = useDebounce(searchQuery, 400)

        // const {data: VideoQuery, isLoading: VideoLoading, error: VideoError } = useVideoQuery(
        //     "getRandomFeed",
        //     {userId: user.userId, limit: 10, offset: 0}
        // )

    const adVideo: VideoFeed = {
        id: -1,
        userId: 0,
        title: "Ad",
        url: "http://localhost:9000/video/4.mp4",
        objectName: "ad",
        enableComment: false,
        visibility: "public",
        createdAt: { seconds: Date.now() / 1000 },
        likeCount: 0,
        commentCount: 0,
        isDraft: false,
        thumbnailUrl: "http://localhost:9000/WE25-1_KPI_Report_Odd2025.jpg",
        isLiked: false
    }

    function insertAds(videos: VideoFeed[]): VideoFeed[] {
        const result: VideoFeed[] = [];
        for (let i = 0; i < videos.length; i += 2) {
            const block = videos.slice(i, i + 2);
            const insertIdx = Math.floor(Math.random() * (block.length + 1));
            block.splice(insertIdx, 0, { ...adVideo, id: -1 * (i + 1) });
            result.push(...block);
        }
        return result;
    }

    const fetchFeed = useCallback(async (reset = false) => {
        setLoading(true)
        try {
            let res;
            if (debouncedSearch && debouncedSearch.trim() !== "") {
                res = await brainrot.searchVideo({ query: debouncedSearch, threshold: 0.05, limit: 10, offset: reset ? 0 : page * 10 })
            } else if (user) {
                res = await brainrot.getRandomFeed({ userId: user.userId, limit: 10, offset: reset ? 0 : page * 10 })
                // console.log(VideoQuery?.response)
            } else {
                res = await brainrot.getRandomFeedLoggedOut({ limit: 10, offset: reset ? 0 : page * 10 })
            }

            const newVideos = (res?.response?.videos || []).map(video => ({
                ...video,
                thumbnailUrl: video.thumbnailUrl
            }));

            const videosWithAds = insertAds(newVideos);

            if (reset) {
                setVideos(videosWithAds)
                setCachedVideos(videosWithAds)
                setCycleIndex(0)
            } else {
                if (newVideos.length === 0) {
                    if (cachedVideos.length > 0) {
                        let nextBatch = []
                        if (cycleIndex + 10 <= cachedVideos.length) {
                            nextBatch = cachedVideos.slice(cycleIndex, cycleIndex + 10)
                        } else {
                            nextBatch = [
                                ...cachedVideos.slice(cycleIndex),
                                ...cachedVideos.slice(0, (cycleIndex + 10) % cachedVideos.length)
                            ]
                        }
                        const nextBatchWithAds = insertAds(nextBatch);
                        setVideos(prev => [...prev, ...nextBatchWithAds])
                        setCycleIndex((prev) => (prev + 10) % cachedVideos.length)
                    }
                    setHasMore(false);
                    return
                } else {
                    setVideos(prev => [...prev, ...videosWithAds])
                    setCachedVideos(prev => [...prev, ...videosWithAds])
                }
            }
            setHasMore(true)
        } catch (e: any) {
            if (setError) setError(e?.message || "Failed to load feed");
        } finally {
            setLoading(false)
        }
    }, [user, brainrot, page, cachedVideos, cycleIndex, debouncedSearch, setError])

    useEffect(() => {
        setPage(0)
        fetchFeed(true)
    }, [user, debouncedSearch])

    useEffect(() => {
        if (!hasMore || loading) return
        const observer = new window.IntersectionObserver(entries => {
            if (entries[0].isIntersecting) {
                setPage(prev => prev + 1)
            }
        }, { threshold: 1 })
        if (loadingRef.current) observer.observe(loadingRef.current)
        return () => observer.disconnect()
    }, [hasMore, loading])

    useEffect(() => {
        if (page === 0) return
        fetchFeed()
    }, [page])

    const handleLikeUpdate = (videoId: number, delta: number, liked: boolean) => {
        setVideos(videos =>
            videos.map(v =>
                v.id === videoId
                    ? { ...v, likeCount: (v.likeCount ?? 0) + delta, isLiked: liked }
                    : v
            )
        );
    };

    useEffect(() => {
        if (videos.length === 0) return;
        const videoId = videos[currentIndex]?.id;
        if (videoId !== undefined) {
            onVideoChange(String(videoId));
            // Scroll into view and sync traversal buttons
            const el = document.querySelector(`[data-video-id='${videoId}']`);
            if (el) {
                el.scrollIntoView({ behavior: 'smooth', block: 'center' });
            }
        }
    }, [currentIndex, videos, onVideoChange]);

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
        const feedEl = document.querySelector('.video-list');
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

    return (
        <div className="video-feed">
            <div className="video-list" style={{ position: 'relative' }}>
                {loading ? (
                    <>
                        <VideoSkeleton />
                        <VideoSkeleton />
                        <VideoSkeleton />
                        <VideoSkeleton />
                    </>
                ) : (
                    videos.map((video, idx) => {
                        console.log("Feed thumbnailUrl:", video.thumbnailUrl);
                        return (
                            <VideoProp
                                key={video.id}
                                video={{
                                    ...video,
                                    isPlaying: activeVideoId === String(video.id), // Use string comparison for activeVideoId
                                    thumbnailUrl: (video as any).thumbnailUrl || undefined,
                                    isLiked: video.isLiked ?? false
                                }}
                                onTogglePlay={() => onVideoChange(String(video.id))}
                                onLikeUpdate={handleLikeUpdate}
                                data-video-id={video.id}
                            />
                        );
                    })
                )}
                <div ref={loadingRef} className="loading-trigger" />
            </div>
            {/* Up/Down buttons fixed on right */}
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
    )
}