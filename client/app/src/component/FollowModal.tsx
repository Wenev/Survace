import { useState, useEffect } from "react";
import { X, Search, Loader2 } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useGrpc } from "../context/ClientContext";
import "../style/component/follow-modal.css";
import defaultAvatar from "../assets/default.jpg";

interface User {
  userId: number;
  username: string;
  avatarUrl?: string;
}

interface FollowModalProps {
  isOpen: boolean;
  onClose: () => void;
  userId: number;
  type: "followers" | "following";
  title?: string;
}

const FollowModal = ({ isOpen, onClose, userId, type, title }: FollowModalProps) => {
  const navigate = useNavigate();
  const { social, galactus } = useGrpc();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState("");

  useEffect(() => {
    if (!isOpen || !userId) return;

    const fetchUsers = async () => {
      setLoading(true);
      setError(null);
      try {
        // Fetch either followers or following based on type
        const response = type === "followers"
          ? await social.listFollowers({ userId })
          : await social.listFollowing({ userId });

        if (!response.response?.data) {
          setUsers([]);
          setLoading(false);
          return;
        }

        // Extract user IDs from follow relationships
        const userIds = response.response.data.map((follow: any) => {
          return type === "followers" ? follow.followerId : follow.followeeId;
        });

        // Fetch user details for each ID
        const usersData: User[] = [];
        for (const id of userIds) {
          try {
            const userRes = await galactus.findByUserId({ id });
            if (userRes.response?.user) {
              const userData = userRes.response.user;
              usersData.push({
                userId: userData.userId,
                username: userData.username,
                avatarUrl: userData.avatarUrl || defaultAvatar
              });
            }
          } catch (err) {
            console.error(`Error fetching user ${id}:`, err);
          }
        }

        setUsers(usersData);
      } catch (err) {
        console.error(`Error fetching ${type}:`, err);
        setError(`Failed to load ${type}. Please try again.`);
      } finally {
        setLoading(false);
      }
    };

    fetchUsers();
  }, [isOpen, userId, type, social, galactus]);

  // Filter users based on search query
  const filteredUsers = searchQuery
    ? users.filter(user =>
        user.username.toLowerCase().includes(searchQuery.toLowerCase())
      )
    : users;

  if (!isOpen) return null;

  const handleUserClick = (userId: number) => {
    onClose();
    navigate(`/profile/${userId}`);
  };

  return (
    <div className="follow-modal-overlay">
      <div className="follow-modal">
        <div className="follow-modal-header">
          <h3>{title || (type === "followers" ? "Followers" : "Following")}</h3>
          <button className="follow-modal-close" onClick={onClose}>
            <X size={24} />
          </button>
        </div>

        <div className="follow-modal-search">
          <Search size={18} className="follow-search-icon" />
          <input
            type="text"
            placeholder="Search users"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="follow-search-input"
          />
        </div>

        <div className="follow-modal-content">
          {loading ? (
            <div className="follow-modal-loading">
              <Loader2 size={32} className="loading-spinner" />
              <p>Loading users...</p>
            </div>
          ) : error ? (
            <div className="follow-modal-error">{error}</div>
          ) : filteredUsers.length === 0 ? (
            <div className="follow-modal-empty">
              {searchQuery
                ? `No users found matching "${searchQuery}"`
                : `No ${type} yet`}
            </div>
          ) : (
            <ul className="follow-user-list">
              {filteredUsers.map(user => (
                <li
                  key={user.userId}
                  className="follow-user-item"
                  onClick={() => handleUserClick(user.userId)}
                >
                  <img
                    src={user.avatarUrl || defaultAvatar}
                    alt={user.username}
                    className="follow-user-avatar"
                  />
                  <div className="follow-user-info">
                    <span className="follow-user-name">{user.username}</span>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>
      </div>
    </div>
  );
};

export default FollowModal;
