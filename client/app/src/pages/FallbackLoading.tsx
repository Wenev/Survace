import "../style/component/feed.css";
import "../style/theme.css";

export default function FallbackLoading() {
    return (
        <div
            style={{
                minHeight: "100vh",
                width: "100vw",
                background: "var(--color-bg)",
                color: "var(--color-text)",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                transition: "background 0.2s, color 0.2s"
            }}
            className={window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? "dark-mode" : ""}
        >
            <div style={{ width: 400, maxWidth: "90vw" }}>
                <div className="video-skeleton" style={{
                    background: "var(--color-card)",
                    borderRadius: 12,
                    padding: 24,
                    boxShadow: "0 2px 16px rgba(0,0,0,0.15)",
                    marginBottom: 24,
                }}>
                    <div className="skeleton-thumbnail" style={{
                        width: "100%",
                        height: 220,
                        borderRadius: 8,
                        background: "var(--color-bg-secondary)",
                        marginBottom: 16,
                        animation: "skeleton-loading 1.2s infinite linear alternate"
                    }} />
                    <div className="skeleton-title" style={{
                        width: "70%",
                        height: 24,
                        borderRadius: 6,
                        background: "var(--color-bg-secondary)",
                        marginBottom: 12,
                        animation: "skeleton-loading 1.2s infinite linear alternate"
                    }} />
                    <div className="skeleton-meta" style={{
                        width: "40%",
                        height: 16,
                        borderRadius: 6,
                        background: "var(--color-bg-secondary)",
                        animation: "skeleton-loading 1.2s infinite linear alternate"
                    }} />
                </div>
            </div>
            <style>
                {`
                @keyframes skeleton-loading {
                    0% { opacity: 0.6; }
                    100% { opacity: 1; }
                }
                `}
            </style>
        </div>
    );
}
