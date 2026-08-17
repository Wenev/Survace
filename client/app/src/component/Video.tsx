import "../style/component/video-prop.css.css"
import { useState } from "react"

const hashtags = ["#anime", "#animeedit", "#animemood", "#aesthetic", "#sad", "#aesthetic", "#fyp", "#foryou", "#fypシ"]

export const VideoPlayer = () => {
    const [isPlaying, setIsPlaying] = useState(false)

    const togglePlay = () => {
        setIsPlaying(!isPlaying)
    }

    const handleHashtagClick = (hashtag: string) => {
        console.log("Clicked hashtag:", hashtag)
    }

    return (
        <div className="video-container">
            <div className="video-player">
                <img
                    src="https://hebbkx1anhila5yf.public.blob.vercel-storage.com/image-nriWDDAmvtsiYnvDg8JzyFZBv4tXhu.png"
                    alt="Video thumbnail"
                    className="video-thumbnail"
                />
                <div className="video-overlay">
                    <h2 className="video-text">but i need to</h2>
                </div>
                <div className="video-controls">
                    <button className="play-button" onClick={togglePlay}>
                        {isPlaying ? "⏸️" : "▶️"}
                    </button>
                </div>
            </div>

            <div className="video-info">
                <div className="hashtags">
                    {hashtags.map((hashtag, index) => (
                        <span key={index} className="hashtag" onClick={() => handleHashtagClick(hashtag)}>
                            {hashtag}
                        </span>
                    ))}
                </div>
            </div>
        </div>
    )
}
