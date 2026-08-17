import "../style/component/loading-skeleton.css"
export const VideoSkeleton = () => {
    return (
        <div className="video-item skeleton">
            <div className="video-content">
                <div className="video-player">
                    <div className="skeleton-thumbnail">
                        <div className="skeleton-shimmer"></div>
                    </div>
                    <div className="video-overlay">
                        <div className="video-info-overlay">
                            <div className="user-info">
                                <div className="skeleton-avatar">
                                    <div className="skeleton-shimmer"></div>
                                </div>
                                <div className="skeleton-username">
                                    <div className="skeleton-shimmer"></div>
                                </div>
                            </div>
                            <div className="skeleton-text">
                                <div className="skeleton-shimmer"></div>
                            </div>
                            <div className="skeleton-hashtags">
                                <div className="skeleton-hashtag">
                                    <div className="skeleton-shimmer"></div>
                                </div>
                                <div className="skeleton-hashtag">
                                    <div className="skeleton-shimmer"></div>
                                </div>
                                <div className="skeleton-hashtag">
                                    <div className="skeleton-shimmer"></div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="video-actions">
                    {[1, 2, 3, 4].map((i) => (
                        <div key={i} className="action-item">
                            <div className="skeleton-action-button">
                                <div className="skeleton-shimmer"></div>
                            </div>
                            <div className="skeleton-action-count">
                                <div className="skeleton-shimmer"></div>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
        </div>
    )
}
