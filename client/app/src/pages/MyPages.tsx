import { useEffect, useState } from "react";
import Navbar from "../component/Navbar";
import MobileNavbar from "../component/MobileNavbar";
import ErrorToast from "../component/ErrorToast";
import { useGrpc } from "../context/ClientContext";
import { useAuth } from "../context/AuthContext";
import defaultAvatar from "../assets/default.jpg";
import "../style/profile.css";

interface PageInfo {
  id: number;
  title: string;
  url: string;
  createdAt: string;
  updatedAt?: string;
  thumbnailUrl?: string;
}

export default function MyPages() {
  const { user } = useAuth();
  const { galactus, brainrot } = useGrpc();
  const [pages, setPages] = useState<PageInfo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!user) return;
    setLoading(true);
    setError(null);
    (async () => {
      try {
        const res = await brainrot.getUserVideo({ userId: user.userId });
        const data = res.response?.pages || [];
        setPages(
          data.map((p: any) => ({
            id: p.id,
            title: p.title,
            url: p.url,
            createdAt: typeof p.createdAt === "string" ? p.createdAt : (p.createdAt && "seconds" in p.createdAt ? new Date(Number(p.createdAt.seconds) * 1000).toISOString() : ""),
            updatedAt: p.updatedAt ? (typeof p.updatedAt === "string" ? p.updatedAt : (p.updatedAt && "seconds" in p.updatedAt ? new Date(Number(p.updatedAt.seconds) * 1000).toISOString() : "")) : undefined,
            thumbnailUrl: p.thumbnailUrl || defaultAvatar,
          }))
        );
      } catch (e: any) {
        setError(e?.message || "Failed to load your pages");
      }
      setLoading(false);
    })();
  }, [user, galactus]);

  return (
    <div className="app">
      <MobileNavbar />
      <Navbar />
      <ErrorToast message={error || ""} onClose={() => setError(null)} />
      <main className="main-content profile-main-content">
        <div className="profile-container">
          <div className="section-title">My Pages</div>
          {loading ? (
            <div style={{ textAlign: "center", padding: "2rem" }}>Loading...</div>
          ) : pages.length === 0 ? (
            <div className="no-videos"><p>You have no pages yet.</p></div>
          ) : (
            <div className="playlists-grid">
              {pages.map(page => (
                <div key={page.id} className="playlist-grid-item" style={{ cursor: "pointer" }} onClick={() => window.open(page.url, "_blank") }>
                  <div className="playlist-thumbnail">
                    <img src={page.thumbnailUrl || defaultAvatar} alt={page.title} />
                  </div>
                  <div className="playlist-info">
                    <h3 className="playlist-name">{page.title}</h3>
                    <div className="playlist-meta">
                      <span className="meta-item">Created: {new Date(page.createdAt).toLocaleDateString()}</span>
                      {page.updatedAt && <span className="meta-item">Updated: {new Date(page.updatedAt).toLocaleDateString()}</span>}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}

