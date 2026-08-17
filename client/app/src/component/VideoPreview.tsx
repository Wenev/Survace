import { useNavigate } from "react-router-dom";
import "../style/component/video-preview.css";
import { Heart, Eye } from "lucide-react";
import type {VideoFeed} from "../generated/dto/video.ts";

interface VideoPreviewProps {
  video: VideoFeed;
  className?: string;
  onVideoClick?: (video: VideoFeed) => void;
}

const VideoPreview = ({ video, className = "", onVideoClick }: VideoPreviewProps) => {
  const navigate = useNavigate();

  const formatViewCount = (count: number): string => {
    if (count >= 1000000) {
      return `${(count / 1000000).toFixed(1)}M`;
    } else if (count >= 1000) {
      return `${(count / 1000).toFixed(1)}K`;
    }
    return count.toString();
  };

  const handleClick = () => {
    if (onVideoClick) {
      onVideoClick(video);
    }
  };

  return (
    <div
      className={`video-preview ${className}`}
      onClick={handleClick}
    >
      <div className="video-preview-thumbnail">
        {video.thumbnailUrl ? (
          <img
            src={video.thumbnailUrl}
            alt={video.title || "Video thumbnail"}
            loading="lazy"
          />
        ) : video.url ? (
          <video
            src={video.url}
            preload="metadata"
            muted
            playsInline
          />
        ) : (
          <div className="video-preview-placeholder">
            <span>No preview available</span>
          </div>
        )}
        <div className="video-preview-stats">
          <div className="video-preview-stat">
            <Heart size={14} />
            <span>{formatViewCount(video.likeCount || 0)}</span>
          </div>
          <div className="video-preview-stat">
            <Eye size={14} />
            <span>{formatViewCount(video.viewCount || 0)}</span>
          </div>
        </div>
      </div>

      <div className="video-preview-info">
        <h3 className="video-preview-title">{video.title}</h3>
      </div>
    </div>
  );
};

export default VideoPreview;
