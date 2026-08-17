import { useState, useEffect, useRef, useLayoutEffect } from "react";
import { useParams, useNavigate } from "react-router-dom";
import "../style/profile.css";
import Navbar from "../component/Navbar";
import MobileNavbar from "../component/MobileNavbar";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import VideoPreview from "../component/VideoPreview";
import { VideoProp } from "../component/VideoProp";
import FollowModal from "../component/FollowModal";
import ErrorToast from "../component/ErrorToast";
import type { VideoFeed, Playlist as PlaylistDto } from "../generated/dto/video";
import type { Video } from "../types/type";
import defaultAvatar from "../assets/default.jpg";
import { Grid, UserCheck2, Users2, Settings, PenSquare, ListPlus, Video as VideoIcon, ListVideo, ThumbsUp, ArrowLeft } from "lucide-react";

interface UserProfile {
  userId: number;
  username: string;
  bio?: string;
  avatarUrl?: string;
  followerCount: number;
  followingCount: number;
  isFollowing: boolean;
}

type Playlist = PlaylistDto;

type ProfileTab = 'videos' | 'playlists' | 'likes';

export default function ProfilePage() {
  const { userId } = useParams<{ userId: string }>();
  const navigate = useNavigate();
  const { user: currentUser } = useAuth();
  const { galactus, brainrot, social, stream } = useGrpc();
  
  const [user, setUser] = useState<UserProfile | null>(null);
  const [videos, setVideos] = useState<VideoFeed[]>([]);
  const [playlists, setPlaylists] = useState<Playlist[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isLive, setIsLive] = useState(false);
  const [activeVideoId, setActiveVideoId] = useState<number | null>(null);
  const [activeVideo, setActiveVideo] = useState<VideoFeed | null>(null);
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [isMobile, setIsMobile] = useState(false);
  const [showFollowersModal, setShowFollowersModal] = useState(false);
  const [showFollowingModal, setShowFollowingModal] = useState(false);
  const [activeTab, setActiveTab] = useState<ProfileTab>('videos');
  const tabContentRef = useRef<HTMLDivElement | null>(null);
  const [tabPositions, setTabPositions] = useState<{ [key in ProfileTab]: number }>({
    videos: 0,
    playlists: 0,
    likes: 0,
  });
  const tabRefs = useRef<(HTMLButtonElement | null)[]>([]);

  const [showPlaylistModal, setShowPlaylistModal] = useState(false);
  const [playlistModalData, setPlaylistModalData] = useState<Playlist | null>(null);
  const [playlistVideos, setPlaylistVideos] = useState<VideoFeed[]>([]);
  const [playlistVideosLoading, setPlaylistVideosLoading] = useState(false);
  const [playlistError, setPlaylistError] = useState<string | null>(null);
  const [playlistActiveVideo, setPlaylistActiveVideo] = useState<VideoFeed | null>(null);
  const [errorToast, setErrorToast] = useState<string>("");

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
  
  useEffect(() => {
    if (!userId) {
      if (currentUser) {
        navigate(`/profile/${currentUser.userId}`);
      } else {
        navigate('/');
      }
      return;
    }
  }, [userId, currentUser, navigate]);

  useEffect(() => {
    if (!userId) {
      return;
    }
    
    const fetchProfileData = async () => {
      try {
        setIsLoading(true);

        const userIdNum = parseInt(userId, 10);
        if (isNaN(userIdNum)) {
          setError("Invalid user ID");
          setIsLoading(false);
          return;
        }

        const userRes = await galactus.findByUserId({ id: userIdNum });
        const userData = userRes.response?.user;
        
        if (!userData) {
          setError("User not found");
          setIsLoading(false);
          return;
        }
        
        let isFollowingUser = false;
        
        if (currentUser && currentUser.userId !== userData.userId) {
          try {
            const followingRes = await social.listFollowing({ userId: currentUser.userId });
            const followingList = followingRes.response?.data || [];
            isFollowingUser = followingList.some((f: any) => f.followeeId === userData.userId);
          } catch (error) {
            console.error("Error checking following status:", error);
          }
        }
        
        let followerCount = 0;
        try {
          const followersRes = await social.listFollowers({ userId: userData.userId });
          followerCount = (followersRes.response?.data || []).length;
        } catch (error) {
          console.error("Error fetching followers:", error);
        }
        
        let followingCount = 0;
        try {
          const followingRes = await social.listFollowing({ userId: userData.userId });
          followingCount = (followingRes.response?.data || []).length;
        } catch (error) {
          console.error("Error fetching following:", error);
        }
        
        try {
          const liveStreamsRes = await stream.listLiveStreams({});
          const liveStreams = liveStreamsRes.response?.streams || [];
          setIsLive(liveStreams.some(stream => stream.userId === userData.userId.toString()));
        } catch (error) {
          console.error("Error checking live status:", error);
        }
        
        setUser({
          userId: userData.userId,
          username: userData.username,
          bio: userData.bio || "No bio yet",
          avatarUrl: userData.avatarUrl || defaultAvatar,
          followerCount,
          followingCount,
          isFollowing: isFollowingUser
        });
        
        try {
          const videosRes = await brainrot.getUserVideo({ userId: userData.userId });
          setVideos(videosRes.response?.videos || []);
        } catch (error) {
          console.error("Error fetching videos:", error);
        }
        
        try {
          const playlistsRes = await brainrot.getUserPlaylists({ userId: userData.userId });
          setPlaylists(playlistsRes.response?.playlists || []);
        } catch (error) {
          console.error("Error fetching playlists:", error);
        }

        setIsLoading(false);
      } catch (error) {
        console.error("Error fetching profile data:", error);
        setError("Failed to load profile data");
        setIsLoading(false);
      }
    };
    
    fetchProfileData();
  }, [userId, currentUser, galactus, social, brainrot, stream]);

  const handleTogglePlay = (video: VideoFeed) => {
    setActiveVideo(prevVideo => prevVideo?.id === video.id ? null : video);
    setActiveVideoId(prevVideo => prevVideo === video.id ? null : video.id);
  };

  const handleBackToProfile = () => {
    setActiveVideo(null);
    setActiveVideoId(null);
  };
  
  const handleFollow = async () => {
    if (!currentUser || !user) return;
    
    try {
      await social.follow({ userId: currentUser.userId, targetId: user.userId });
      setUser(prev => prev ? { ...prev, isFollowing: true, followerCount: prev.followerCount + 1 } : null);
    } catch (error) {
      console.error("Error following user:", error);
    }
  };
  
  const handleUnfollow = async () => {
    if (!currentUser || !user) return;
    
    try {
      await social.unfollow({ userId: currentUser.userId, targetId: user.userId });
      setUser(prev => prev ? { ...prev, isFollowing: false, followerCount: prev.followerCount - 1 } : null);
    } catch (error) {
      console.error("Error unfollowing user:", error);
    }
  };
  
  const handleWatchLive = () => {
    if (user) {
      navigate(`/live?user=${user.userId}`);
    }
  };
  
  const handleLikeUpdate = (videoId: number, delta: number, liked: boolean) => {
    setVideos(prevVideos => 
      prevVideos.map(video => 
        video.id === videoId 
          ? { ...video, likeCount: (video.likeCount || 0) + delta, isLiked: liked } 
          : video
      )
    );
  };
  
  const handleEditProfile = () => {
    navigate('/settings');
  };

  const handleCreatePlaylist = () => {
    navigate('/playlist/create');
  };

  const handleGoToStudio = () => {
    navigate('/studio');
  };

  const handleTabChange = (tab: ProfileTab) => {
    setActiveTab(tab);
  };

  const handleOpenPlaylistModal = async (playlist: Playlist) => {
    setShowPlaylistModal(true);
    setPlaylistModalData(playlist);
    setPlaylistVideos([]);
    setPlaylistActiveVideo(null);
    setPlaylistVideosLoading(true);
    setPlaylistError(null);

    try {
      const videos: VideoFeed[] = [];
      for (const pv of playlist.videos || []) {
        try {
          const res = await brainrot.getVideoByID({ id: pv.videoId });
          if (res.response?.video) {
            videos.push(res.response.video);
          }
        } catch (err) {
        }
      }
      setPlaylistVideos(videos);
    } catch (err) {
      setPlaylistError("Failed to load playlist videos");
      setErrorToast("Failed to load playlist videos");
    }
    setPlaylistVideosLoading(false);
  };

  const handleClosePlaylistModal = () => {
    setShowPlaylistModal(false);
    setPlaylistModalData(null);
    setPlaylistVideos([]);
    setPlaylistActiveVideo(null);
    setPlaylistError(null);
  };

  const handlePlaylistVideoClick = (video: VideoFeed) => {
    setActiveVideo(video);
    setActiveVideoId(video.id);
    setShowPlaylistModal(false);
  };

  const handleBackToPlaylistModal = () => {
    setActiveVideo(null);
    setActiveVideoId(null);
    setShowPlaylistModal(true);
  };

  useLayoutEffect(() => {
    if (tabRefs.current.length === 0) return;

    const updateTabPositions = () => {
      const newPositions: { [key in ProfileTab]: number } = {
        videos: 0,
        playlists: 0,
        likes: 0,
      };

      tabRefs.current.forEach((tab, index) => {
        if (tab) {
          const { offsetLeft, clientWidth } = tab;
          newPositions[Object.keys(newPositions)[index] as ProfileTab] = offsetLeft + clientWidth / 2;
        }
      });

      setTabPositions(newPositions);
    };

    updateTabPositions();
    window.addEventListener("resize", updateTabPositions);
    return () => window.removeEventListener("resize", updateTabPositions);
  }, [tabRefs]);

  return (
    <div className="app">
      <MobileNavbar onToggleSidebar={() => setSidebarOpen(!sidebarOpen)} isOpen={isMobile && sidebarOpen} />
      <Navbar isOpen={sidebarOpen} />
      
      <ErrorToast message={errorToast} onClose={() => setErrorToast("")} />
      <main className="main-content profile-main-content">
        {activeVideo ? (
          <div className="video-viewer">
            <div className="video-viewer-header">
              <button className="back-button" onClick={showPlaylistModal ? handleBackToPlaylistModal : handleBackToProfile}>
                <ArrowLeft size={24} /> {showPlaylistModal ? "Back to Playlist" : "Back to Profile"}
              </button>
            </div>
            <div className="video-viewer-content">
              <VideoProp
                video={{
                  ...activeVideo,
                  isPlaying: true,
                  thumbnailUrl: (activeVideo as any).thumbnailUrl || undefined,
                  isLiked: (activeVideo as any).isLiked || false
                }}
                onTogglePlay={() => {}}
                onLikeUpdate={handleLikeUpdate}
              />
            </div>
          </div>
        ) : (
          <div className="profile-container">
            {isLoading ? (
              <div className="profile-loading">
                <div className="profile-avatar-skeleton"></div>
                <div className="profile-info-skeleton">
                  <div className="profile-username-skeleton"></div>
                  <div className="profile-bio-skeleton"></div>
                  <div className="profile-stats-skeleton"></div>
                </div>
              </div>
            ) : error ? (
              <div className="profile-error">{error}</div>
            ) : user ? (
              <>
                <div className="profile-header">
                  <div className="profile-avatar-container">
                    <img src={user.avatarUrl} alt={user.username} className="profile-avatar" />
                    {isLive && (
                      <div className="live-badge" onClick={handleWatchLive}>
                        LIVE
                      </div>
                    )}
                  </div>

                  <div className="profile-info">
                    <div className="profile-username">
                      <h1>{user.username}</h1>
                      {currentUser && currentUser.userId !== user.userId ? (
                        user.isFollowing ? (
                          <button className="unfollow-button" onClick={handleUnfollow}>Unfollow</button>
                        ) : (
                          <button className="follow-button" onClick={handleFollow}>Follow</button>
                        )
                      ) : currentUser && currentUser.userId === user.userId && (
                        <div className="user-actions">
                          <button className="edit-profile-button" onClick={handleEditProfile}>
                            <Settings size={16} /> Edit Profile
                          </button>
                          <button className="create-playlist-button" onClick={handleCreatePlaylist}>
                            <ListPlus size={16} /> Create Playlist
                          </button>
                          <button className="studio-page-button" onClick={handleGoToStudio}>
                            <VideoIcon size={16} /> Studio Page
                          </button>
                        </div>
                      )}
                    </div>

                    <div className="profile-stats">
                      <div className="stat-item">
                        <span className="stat-value">{videos.length}</span>
                        <span className="stat-label"><Grid size={16} /> Videos</span>
                      </div>
                      <div className="stat-item clickable" onClick={() => setShowFollowersModal(true)}>
                        <span className="stat-value">{user.followerCount}</span>
                        <span className="stat-label"><Users2 size={16} /> Followers</span>
                      </div>
                      <div className="stat-item clickable" onClick={() => setShowFollowingModal(true)}>
                        <span className="stat-value">{user.followingCount}</span>
                        <span className="stat-label"><UserCheck2 size={16} /> Following</span>
                      </div>
                    </div>

                    <div className="profile-bio">
                      {user.bio}
                    </div>
                  </div>
                </div>

                <div className="profile-content">
                  <div className="profile-tabs">
                    <button
                      className={`profile-tab ${activeTab === 'videos' ? 'active' : ''}`}
                      onClick={() => handleTabChange('videos')}
                      ref={el => tabRefs.current[0] = el}
                    >
                      <VideoIcon size={16} /> Videos
                    </button>
                    <button
                      className={`profile-tab ${activeTab === 'playlists' ? 'active' : ''}`}
                      onClick={() => handleTabChange('playlists')}
                      ref={el => tabRefs.current[1] = el}
                    >
                      <ListVideo size={16} /> Playlists
                    </button>
                    <button
                      className={`profile-tab ${activeTab === 'likes' ? 'active' : ''}`}
                      onClick={() => handleTabChange('likes')}
                      ref={el => tabRefs.current[2] = el}
                    >
                      <ThumbsUp size={16} /> Likes
                    </button>

                    <div className="tab-slider" style={{ left: tabPositions[activeTab] }} />
                  </div>

                  <div className="tab-content" ref={tabContentRef}>
                    {activeTab === 'videos' && (
                      <div className="profile-videos">
                        {videos.length === 0 ? (
                          <div className="no-videos">
                            <p>No videos yet</p>
                          </div>
                        ) : (
                          <div className="videos-grid">
                            {videos.map(video => (
                                <div key={video.id} className="video-grid-item">
                                <VideoPreview
                                  video={video}
                                  onVideoClick={handleTogglePlay}
                                />
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    )}

                    {activeTab === 'playlists' && (
                      <div className="profile-playlists">
                        {playlists.length === 0 ? (
                          <div className="no-playlists">
                            <p>No playlists yet</p>
                          </div>
                        ) : (
                          <div className="playlists-grid">
                            {playlists.map(playlist => (
                              <div
                                key={playlist.id}
                                className="playlist-grid-item"
                                style={{ cursor: "pointer" }}
                                onClick={() => handleOpenPlaylistModal(playlist)}
                              >
                                <div className="playlist-thumbnail">
                                  <img
                                    src={
                                      playlist.videos && playlist.videos.length > 0
                                        ? playlist.videos[0].video?.thumbnailUrl || defaultAvatar
                                        : defaultAvatar
                                    }
                                    alt={playlist.title}
                                  />
                                </div>
                                <div className="playlist-info">
                                  <h3 className="playlist-name">{playlist.title}</h3>
                                  <div className="playlist-meta">
                                    <span className="meta-item">{playlist.videos?.length || 0} videos</span>
                                  </div>
                                </div>
                              </div>
                            ))}
                          </div>
                        )}
                        {showPlaylistModal && playlistModalData && (
                          <div className="modal-overlay" style={{ zIndex: 2000 }}>
                            <div className="modal-content" style={{ maxWidth: 600, minHeight: 300, position: "relative" }}>
                              <button
                                style={{
                                  position: "fixed",
                                  top: 24,
                                  left: 24,
                                  background: "var(--color-card)",
                                  border: "1px solid var(--color-border)",
                                  borderRadius: 8,
                                  fontSize: 15,
                                  cursor: "pointer",
                                  color: "var(--color-primary)",
                                  fontWeight: 500,
                                  display: "flex",
                                  alignItems: "center",
                                  gap: 4,
                                  zIndex: 2100,
                                  boxShadow: "0 2px 8px rgba(0,0,0,0.10)",
                                  padding: "8px 16px"
                                }}
                                onClick={handleClosePlaylistModal}
                              >
                                <ArrowLeft size={18} /> Back
                              </button>
                              <button
                                style={{
                                  position: "fixed",
                                  top: 24,
                                  right: 24,
                                  background: "var(--color-card)",
                                  border: "1px solid var(--color-border)",
                                  borderRadius: 8,
                                  fontSize: 20,
                                  cursor: "pointer",
                                  color: "var(--color-text)",
                                  zIndex: 2100,
                                  boxShadow: "0 2px 8px rgba(0,0,0,0.10)",
                                  padding: "8px 16px"
                                }}
                                onClick={handleClosePlaylistModal}
                                aria-label="Close"
                              >
                                ×
                              </button>
                              <h3 style={{ marginTop: 0, marginLeft: 32 }}>{playlistModalData.title}</h3>
                              {playlistVideosLoading ? (
                                <div>Loading videos...</div>
                              ) : playlistError ? (
                                <div style={{ color: "var(--color-error)" }}>{playlistError}</div>
                              ) : playlistVideos.length === 0 ? (
                                <div>No videos in this playlist.</div>
                              ) : (
                                <ul style={{ listStyle: "none", padding: 0, margin: 0 }}>
                                  {playlistVideos.map(video => (
                                    <li
                                      key={video.id}
                                      style={{
                                        display: "flex",
                                        alignItems: "center",
                                        gap: "1rem",
                                        padding: "10px 0",
                                        borderBottom: "1px solid var(--color-border)",
                                        cursor: "pointer"
                                      }}
                                      onClick={() => handlePlaylistVideoClick(video)}
                                    >
                                      <img
                                        src={video.thumbnailUrl || defaultAvatar}
                                        alt={video.title}
                                        style={{
                                          width: 80,
                                          height: 48,
                                          objectFit: "cover",
                                          borderRadius: 4,
                                          background: "var(--color-bg-secondary)"
                                        }}
                                      />
                                      <div style={{ flex: 1, minWidth: 0 }}>
                                        <div style={{
                                          fontWeight: 500,
                                          whiteSpace: "nowrap",
                                          overflow: "hidden",
                                          textOverflow: "ellipsis"
                                        }}>
                                          {video.title}
                                        </div>
                                        <div style={{ fontSize: 13, color: "var(--color-text-secondary)", marginTop: 2 }}>
                                          {video.viewCount} views · {video.commentCount} comments · {video.likeCount} likes
                                        </div>
                                      </div>
                                    </li>
                                  ))}
                                </ul>
                              )}
                            </div>
                          </div>
                        )}
                      </div>
                    )}

                    {activeTab === 'likes' && (
                      <div className="profile-likes">
                        {videos.length === 0 ? (
                          <div className="no-likes">
                            <p>No liked videos yet</p>
                          </div>
                        ) : (
                          <div className="likes-grid">
                            {videos.filter(video => video.isLiked).map(video => (
                              <div key={video.id} className="video-grid-item">
                                <VideoPreview
                                  video={video}
                                  onVideoClick={handleTogglePlay}
                                />
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                </div>

                {showFollowersModal && (
                  <FollowModal
                    userId={user.userId}
                    onClose={() => setShowFollowersModal(false)}
                    title="Followers"
                    fetchFollowers={() => social.listFollowers({ userId: user.userId })}
                  />
                )}

                {showFollowingModal && (
                  <FollowModal
                    userId={user.userId}
                    onClose={() => setShowFollowingModal(false)}
                    title="Following"
                    fetchFollowers={() => social.listFollowing({ userId: user.userId })}
                  />
                )}
              </>
            ) : (
              <div className="profile-error">User not found</div>
            )}
          </div>
        )}
      </main>
    </div>
  );
}