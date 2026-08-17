import React, { createContext, useContext, useState, useEffect, useCallback } from "react";
import { useAuth } from "./AuthContext";
import { useGrpc } from "./ClientContext";

export interface Notification {
  id: number;
  type: string;
  message: string;
  createdAt: string;
  isRead: boolean;
}

interface NotifContextType {
  notifications: Notification[];
  unreadCount: number;
  refreshNotifications: () => Promise<void>;
  markAsRead: (notifId: number) => Promise<void>;
  loading: boolean;
  error: string | null;
}

const NotifContext = createContext<NotifContextType | undefined>(undefined);

export const useNotif = () => {
  const ctx = useContext(NotifContext);
  if (!ctx) throw new Error("useNotif must be used within NotifProvider");
  return ctx;
};

export const NotifProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user } = useAuth();
  const { notification } = useGrpc();
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchNotifications = useCallback(async () => {
    if (!user) return;
    setLoading(true);
    setError(null);
    try {
      const res = await notification.getNotificationsByUserID({ userId: user.userId });
      const data = res.response?.notifications || [];
      setNotifications(
        data.map((n: any) => {
          let message = n.message;
          let mentions: string[] = [];
          if (n.type === "new-follower") {
            if (n.mentions && Array.isArray(n.mentions)) {
              mentions = n.mentions;
            } else if (n.payload && n.payload.followerUsername) {
              mentions = [n.payload.followerUsername];
            }
            if (mentions.length > 0) {
              message = `${mentions[0]} started following you`;
            } else {
              message = "You have a new follower";
            }
          }
          return {
            id: n.id,
            type: n.type,
            message,
            mentions,
            createdAt: typeof n.createdAt === "string" ? n.createdAt : (n.createdAt && "seconds" in n.createdAt ? new Date(Number(n.createdAt.seconds) * 1000).toISOString() : ""),
            isRead: !!n.isRead,
          };
        })
      );

      if (typeof window !== "undefined" && window.Notification && Notification.permission === "granted") {
        data.forEach((n: any) => {
          if (n.isRead === false) {
            let notifMsg = n.message;
            if (n.type === "new-follower" && n.mentions && n.mentions.length > 0) {
              notifMsg = `${n.mentions[0]} started following you`;
            }
            new Notification("Notification", {
              body: notifMsg,
              icon: "/vite.svg"
            });
          }
        });
      }
    } catch (e: any) {
      setError(e?.message || "Failed to fetch notifications");
    }
    setLoading(false);
  }, [user, notification]);

  useEffect(() => {
    if (typeof window !== "undefined" && window.Notification && Notification.permission !== "granted") {
      Notification.requestPermission();
    }
  }, []);

  useEffect(() => {
    fetchNotifications();
  }, [fetchNotifications]);

  const markAsRead = useCallback(async (notifId: number) => {
    if (!user) return;
    try {
      await notification.markAsRead({ notificationId: notifId });
      setNotifications((prev) => prev.map(n => n.id === notifId ? { ...n, isRead: true } : n));
    } catch (e) {
    }
  }, [user, notification]);

  const unreadCount = notifications.filter(n => !n.isRead).length;

  return (
    <NotifContext.Provider value={{ notifications, unreadCount, refreshNotifications: fetchNotifications, markAsRead, loading, error }}>
      {children}
    </NotifContext.Provider>
  );
};