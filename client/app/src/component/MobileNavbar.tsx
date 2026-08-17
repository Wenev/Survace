"use client"
import "../style/component/comment.css"

interface MobileHeaderProps {
    onToggleSidebar: () => void
    isOpen: boolean
    searchQuery: string
    setSearchQuery: (q: string) => void
}

export default function MobileHeader({ onToggleSidebar, isOpen, searchQuery, setSearchQuery }: MobileHeaderProps) {
    return (
        <header className={`mobile-header${isOpen ? '' : ' hidden'}`} style={{ display: isOpen ? undefined : 'none' }}>
            <button className={`menu-toggle${isOpen ? " active" : ""}`} onClick={onToggleSidebar}>
                <span></span>
                <span></span>
                <span></span>
            </button>
            <div className="logo">
                <span className="tiktok-icon">🎵</span>
                TikTok
            </div>
            <form className="search-container" style={{ flex: 1, marginLeft: 16 }}>
                <input
                    type="text"
                    placeholder="Search"
                    className="search-input"
                    value={searchQuery}
                    onChange={e => setSearchQuery(e.target.value)}
                    style={{ width: '100%' }}
                />
                <span className="search-icon">🔍</span>
            </form>
        </header>
    )
}
