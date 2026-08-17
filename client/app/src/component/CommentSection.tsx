"use client"

import type React from "react"
import { useState, useEffect } from "react"
import "../style/component/comment.css"
import defaultAvatar from "../assets/default.jpg";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import { Heart, MessageCircle, HeartOff } from "lucide-react";
import RichText from "./RichText";

interface Reply {
    id: number;
    username: string;
    avatar: string;
    text: string;
    time: string;
    likes: number;
    isLiked?: boolean;
}

interface Comment {
    id: number;
    username: string;
    avatar: string;
    text: string;
    time: string;
    likes: number;
    isLiked?: boolean;
    replies?: Reply[];
}

interface CommentsSectionProps {
    isOpen: boolean
    onClose: () => void
    videoId: string | null
}

export default function CommentsSection({ isOpen, onClose, videoId }: CommentsSectionProps) {
    const { brainrot, galactus } = useGrpc();
    const { user } = useAuth();
    if (!user) return null;
    const [comments, setComments] = useState<Comment[]>([])
    const [newComment, setNewComment] = useState("")
    const [loading, setLoading] = useState(false)
    const [replyingTo, setReplyingTo] = useState<{commentId: number, replyId?: number} | null>(null);
    const [openReplies, setOpenReplies] = useState<{ [commentId: number]: boolean }>({});
    const [repliesLoading, setRepliesLoading] = useState<{ [commentId: number]: boolean }>({});

    useEffect(() => {
        if (videoId && isOpen) {
            setLoading(true);
            (async () => {
                try {
                    const res = await brainrot.getCommentsByVideo({ videoId: Number(videoId) });
                    let commentsFetched = res.response?.comments?.map((c: any) => ({
                        id: c. id,
                        userId: c.userId,
                        username: c.username || `User${c.userId}`,
                        avatar: c.avatarUrl ? c.avatarUrl : "/assets/default.jpg",
                        text: c.content,
                        time: c.createdAt || "",
                        likes: c.likeCount || 0,
                        isLiked: false,
                        replies: [],
                    })) || [];
                    await Promise.all(commentsFetched.map(async (c, idx) => {
                        const likedRes = await brainrot.isCommentLiked({ userId: user.userId, commentId: c.id });
                        commentsFetched[idx].isLiked = likedRes?.response?.liked === true;
                    }));
                    commentsFetched = await Promise.all(commentsFetched.map(async (c) => {
                        try {
                            const userRes = await galactus.findByUserId({ id: c.userId });
                            const userData = userRes.response.user;
                            return {
                                ...c,
                                username: userData?.username || c.username,
                                avatar: userData?.avatarUrl || c.avatar,
                            };
                        } catch {
                            return c;
                        }
                    }));
                    setComments(commentsFetched);
                } finally {
                    setLoading(false);
                }
            })();
        }
    }, [videoId, isOpen, brainrot, user, galactus])

    const fetchReplies = async (commentId: number) => {
        setRepliesLoading(r => ({ ...r, [commentId]: true }));
        try {
            const res = await brainrot.getRepliesByComment({ commentId });
            let fetchedReplies = res.response?.replies?.map((r: any) => ({
                id: r.id,
                userId: r.userId,
                username: r.username || `User${r.userId}`,
                avatar: r.avatarUrl ? r.avatarUrl : "/assets/default.jpg",
                text: r.content,
                time: r.createdAt || "",
                likes: r.likeCount || 0,
                isLiked: false,
            })) || [];
            await Promise.all(fetchedReplies.map(async (r, idx) => {
                const likedRes = await brainrot.isReplyLiked({ userId: user.userId, replyId: r.id });
                fetchedReplies[idx].isLiked = likedRes?.response?.liked === true;
            }));
            fetchedReplies = await Promise.all(fetchedReplies.map(async (r) => {
                try {
                    const userRes = await galactus.findByUserId({ id: r.userId });
                    const userData = userRes.response.user;
                    return {
                        ...r,
                        username: userData?.username,
                        avatar: userData?.avatarUrl,
                    };
                } catch {
                    return r;
                }
            }));
            setComments(comments => comments.map(c =>
                c.id === commentId ? { ...c, replies: fetchedReplies } : c
            ));
        } finally {
            setRepliesLoading(r => ({ ...r, [commentId]: false }));
        }
    };

    const handleToggleReplies = (commentId: number) => {
        setOpenReplies(open => {
            const isOpen = open[commentId];
            if (!isOpen) fetchReplies(commentId);
            return { ...open, [commentId]: !isOpen };
        });
    };

    const updateCommentLikeCount = async (commentId: number) => {
        if (!videoId) return;
        const res = await brainrot.getCommentsByVideo({ videoId: Number(videoId) });
        const updated = res?.response?.comments?.find((c: any) => c.id === commentId);
        if (!updated) return;
        let isLiked = false;
        if (user) {
            const likedRes = await brainrot.isCommentLiked({ userId: user.userId, commentId });
            isLiked = !!likedRes?.response?.liked;
        }
        setComments(comments => comments.map(c =>
            c.id === commentId ? { ...c, likes: updated.likeCount || 0, isLiked } : c
        ));
    };

    const updateReplyLikeCount = async (commentId: number, replyId: number) => {
        const res = await brainrot.getRepliesByComment({ commentId });
        const updated = res?.response?.replies?.find((r: any) => r.id === replyId);
        if (!updated) return;
        let isLiked = false;
        if (user) {
            const likedRes = await brainrot.isReplyLiked({ userId: user.userId, replyId });
            isLiked = !!likedRes?.response?.liked;
        }
        setComments(comments => comments.map(c =>
            c.id === commentId
                ? {
                    ...c,
                    replies: c.replies?.map(r =>
                        r.id === replyId ? { ...r, likes: updated.likeCount || 0, isLiked } : r
                    )
                }
                : c
        ));
    };

    const handleLikeComment = async (commentId: number) => {
        if (!user) return;
        const comment = comments.find(c => c.id === commentId);
        if (!comment) return;
        const res = await brainrot.isCommentLiked({ userId: user.userId, commentId });
        console.log('isCommentLiked response:', res);
        const backendLiked = !!res.response?.liked;
        if (backendLiked) {
            const unlikeRes = await brainrot.unlikeComment({ userId: user.userId, commentId });
            console.log('unlikeComment response:', unlikeRes);
            setComments(comments => comments.map(c =>
                c.id === commentId ? { ...c, isLiked: false, likes: Math.max(0, c.likes - 1) } : c
            ));
        } else {
            const likeRes = await brainrot.likeComment({ userId: user.userId, commentId });
            console.log('likeComment response:', likeRes);
            setComments(comments => comments.map(c =>
                c.id === commentId ? { ...c, isLiked: true, likes: c.likes + 1 } : c
            ));
        }
    };

    const handleLikeReply = async (commentId: number, replyId: number) => {
        if (!user) return;
        const comment = comments.find(c => c.id === commentId);
        if (!comment || !comment.replies) return;
        const reply = comment.replies.find(r => r.id === replyId);
        if (!reply) return;
        const res = await brainrot.isReplyLiked({ userId: user.userId, replyId });
        console.log('isReplyLiked response:', res);
        const backendLiked = !!res.response?.liked;
        if (backendLiked) {
            const unlikeRes = await brainrot.unlikeReply({ userId: user.userId, replyId });
            console.log('unlikeReply response:', unlikeRes);
            setComments(comments => comments.map(c =>
                c.id === commentId
                    ? {
                        ...c,
                        replies: c.replies?.map(r =>
                            r.id === replyId ? { ...r, isLiked: false, likes: Math.max(0, r.likes - 1) } : r
                        )
                    }
                    : c
            ));
        } else {
            const likeRes = await brainrot.likeReply({ userId: user.userId, replyId });
            console.log('likeReply response:', likeRes);
            setComments(comments => comments.map(c =>
                c.id === commentId
                    ? {
                        ...c,
                        replies: c.replies?.map(r =>
                            r.id === replyId ? { ...r, isLiked: true, likes: r.likes + 1 } : r
                        )
                    }
                    : c
            ));
        }
    };
    const handleReply = (username: string, commentId: number, replyId?: number) => {
        setReplyingTo({ commentId, replyId });
        setNewComment(`@${username} `);
    };
    const handleAddComment = async () => {
        if (!newComment.trim() || !user || !videoId) return;
        if (replyingTo) {
            setComments(comments => comments.map(c => {
                if (c.id !== replyingTo.commentId) return c;
                const reply: Reply = {
                    id: Date.now(),
                    username: "You",
                    avatar: "/placeholder.svg",
                    text: newComment.trim(),
                    time: "now",
                    likes: 0,
                    isLiked: false,
                };
                return {
                    ...c,
                    replies: [reply, ...(c.replies || [])]
                };
            }));
            setReplyingTo(null);
        } else {
            setLoading(true);
            try {
                const res = await brainrot.addComment({ userId: user.userId, videoId: Number(videoId), text: newComment.trim() });
                console.log('addComment response:', res);
                const res2 = await brainrot.getCommentsByVideo({ videoId: Number(videoId) });
                console.log('getCommentsByVideo after addComment response:', res2);
                const fetched = res2.response?.comments?.map((c: any) => ({
                    id: c.id,
                    username: c.username,
                    avatar: c.avatarUrl || defaultAvatar,
                    text: c.content,
                    time: c.createdAt || "",
                    likes: c.likeCount || 0,
                    isLiked: c.isLiked || false,
                    replies: [],
                })) || [];
                setComments(fetched);
            } finally {
                setLoading(false);
            }
        }
        setNewComment("")
    }

    const handleKeyPress = (e: React.KeyboardEvent) => {
        if (e.key === "Enter") {
            handleAddComment()
        }
    }

    const renderRichText = (text: string) => (
        <RichText value={text} placeholder="" />
    );

    const handleRichTextChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setNewComment(e.target.value);
    };

    return (
        <aside className={`comments-section${isOpen ? " active" : ""}`} style={{ display: isOpen ? undefined : 'none' }}>
            <div className="comments-header">
                <h3>Comments ({comments.length})</h3>
                <button className="close-comments" onClick={onClose}>
                    ✕
                </button>
            </div>
            <div className="comments-list">
                {loading ? (
                    <div className="comments-loading">
                        {[1, 2, 3].map((i) => (
                            <div key={i} className="comment-skeleton">
                                <div className="skeleton-avatar">
                                    <div className="skeleton-shimmer"></div>
                                </div>
                                <div className="skeleton-comment-content">
                                    <div className="skeleton-username">
                                        <div className="skeleton-shimmer"></div>
                                    </div>
                                    <div className="skeleton-comment-text">
                                        <div className="skeleton-shimmer"></div>
                                    </div>
                                </div>
                            </div>
                        ))}
                    </div>
                ) : (
                    comments.map((comment) => (
                        <div key={comment.id} className="comment" style={{background: 'var(--color-bg-secondary)', borderRadius: 12, boxShadow: '0 2px 8px var(--color-shadow)', marginBottom: 18, padding: 14}}>
                            <div className="comment-avatar">
                                <img src={comment.avatar || defaultAvatar} alt={`${comment.username} avatar`} style={{width: 40, height: 40, borderRadius: '50%', objectFit: 'cover', border: '2px solid var(--color-border)'}} />
                            </div>
                            <div className="comment-content">
                                <div className="comment-header" style={{display: 'flex', alignItems: 'center', gap: 8, marginBottom: 2}}>
                                    <span className="username" style={{fontWeight: 700, fontSize: 15}}>{comment.username}</span>
                                    <span className="comment-time" style={{color: 'var(--color-text-secondary)', fontSize: 12}}>{comment.time}</span>
                                </div>
                                <div className="comment-text" style={{fontSize: 15, margin: 0, color: 'var(--color-text)'}}>
                                    {renderRichText(comment.text)}
                                </div>
                                <div className="comment-actions" style={{display: 'flex', alignItems: 'center', gap: 10, marginTop: 6}}>
                                    <button className="reply-btn" onClick={() => handleReply(comment.username, comment.id)} style={{fontSize: 13, color: 'var(--color-primary)', background: 'none', border: 'none', cursor: 'pointer', fontWeight: 600}}>
                                        <MessageCircle size={16} style={{marginRight: 4}} /> Reply
                                    </button>
                                    <button className="reply-btn" onClick={() => handleLikeComment(comment.id)} style={{fontSize: 13, color: comment.isLiked ? 'var(--color-primary)' : 'var(--color-text-secondary)', background: 'none', border: 'none', cursor: 'pointer', fontWeight: 600, display: 'flex', alignItems: 'center', gap: 4}}>
                                        {comment.isLiked ? (
                                            <Heart size={16} fill="#fff" color="#fe2c55" data-like={comment.id + '-filled'} />
                                        ) : (
                                            <Heart size={16} data-like={comment.id + '-outline'} />
                                        )}
                                        <span className="count-text">{comment.likes}</span>
                                    </button>
                                    <button className="reply-btn" onClick={() => handleToggleReplies(comment.id)} style={{fontSize: 13, color: 'var(--color-text-secondary)', background: 'none', border: 'none', cursor: 'pointer', fontWeight: 600, display: 'flex', alignItems: 'center', gap: 4}}>
                                        <span style={{display: 'inline-block', transform: openReplies[comment.id] ? 'rotate(180deg)' : 'rotate(0deg)', transition: 'transform 0.2s'}}>
                                            ▼
                                        </span>
                                        {loading ? <span className="skeleton-reply-count" style={{display: 'inline-block', width: 32, height: 12, background: 'var(--color-bg)', borderRadius: 6}}></span> : <>{openReplies[comment.id] ? 'Hide Replies' : `Show Replies${comment.replies && comment.replies.length > 0 ? ` (${comment.replies.length})` : ''}`}</>}
                                    </button>
                                </div>
                                <div className={openReplies[comment.id] ? "replies-dropdown open" : "replies-dropdown"}>
                                    {openReplies[comment.id] && (
                                        <div>
                                            {repliesLoading[comment.id] ? (
                                                <div className="comments-loading">
                                                    {[1, 2].map((i) => (
                                                        <div key={i} className="comment-skeleton">
                                                            <div className="skeleton-avatar">
                                                                <div className="skeleton-shimmer"></div>
                                                            </div>
                                                            <div className="skeleton-comment-content">
                                                                <div className="skeleton-username">
                                                                    <div className="skeleton-shimmer"></div>
                                                                </div>
                                                                <div className="skeleton-comment-text">
                                                                    <div className="skeleton-shimmer"></div>
                                                                </div>
                                                            </div>
                                                        </div>
                                                    ))}
                                                </div>
                                            ) : (
                                                (comment.replies && comment.replies.length > 0) ? comment.replies.map(reply => (
                                                    <div key={reply.id} className="reply">
                                                        <div className="comment-avatar">
                                                            <img src={reply.avatarUrl || defaultAvatar} alt={`${reply.username}`} style={{width: 32, height: 32, borderRadius: '50%', objectFit: 'cover', border: '2px solid var(--color-border)'}} />
                                                        </div>
                                                        <div className="comment-content">
                                                            <div className="comment-header" style={{display: 'flex', alignItems: 'center', gap: 8, marginBottom: 2}}>
                                                                <span className="username" style={{fontWeight: 600, fontSize: 14}}>{reply.username}</span>
                                                                <span className="comment-time" style={{color: 'var(--color-text-secondary)', fontSize: 12}}>{reply.time}</span>
                                                            </div>
                                                            <div className="comment-text" style={{fontSize: 14, margin: 0, color: 'var(--color-text)'}}>
                                                                {renderRichText(reply.text)}
                                                            </div>
                                                            <div className="comment-actions" style={{display: 'flex', alignItems: 'center', gap: 8, marginTop: 4}}>
                                                                <button className="reply-btn" onClick={() => handleReply(reply.username, comment.id, reply.id)} style={{fontSize: 12, color: 'var(--color-primary)', background: 'none', border: 'none', cursor: 'pointer', fontWeight: 600}}>
                                                                    <MessageCircle size={14} style={{marginRight: 3}} /> Reply
                                                                </button>
                                                                <button className="reply-btn" onClick={() => handleLikeReply(comment.id, reply.id)} style={{fontSize: 12, color: reply.isLiked ? 'var(--color-primary)' : 'var(--color-text-secondary)', background: 'none', border: 'none', cursor: 'pointer', fontWeight: 600, display: 'flex', alignItems: 'center', gap: 3}}>
                                                                    {repliesLoading[comment.id] ? <span className="skeleton-like-count" style={{display: 'inline-block', width: 18, height: 10, background: 'var(--color-bg)', borderRadius: 5}}></span> : (
                                                                        <>
                                                                            {reply.isLiked ? <Heart fill="var(--color-primary)" color="var(--color-primary)" size={14} /> : <HeartOff size={14} color="var(--color-text-secondary)" />}
                                                                            {reply.likes}
                                                                        </>
                                                                    )}
                                                                </button>
                                                            </div>
                                                        </div>
                                                    </div>
                                                )) : null
                                            )}
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    ))
                )}
            </div>
            <div className="comment-input-container">
                <RichText
                    value={newComment}
                    placeholder={replyingTo ? "Reply..." : "Add comment..."}
                    onChange={handleRichTextChange}
                />
                <button className="send-comment" onClick={handleAddComment}>
                    📤
                </button>
            </div>
        </aside>
    )
}