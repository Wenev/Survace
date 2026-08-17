import { X } from "lucide-react";
import React, { useEffect } from "react";

interface ErrorToastProps {
    message: string;
    onClose: () => void;
    duration?: number;
}

const ErrorToast: React.FC<ErrorToastProps> = ({ message, onClose, duration = 3500 }) => {
    useEffect(() => {
        if (!message) return;
        const timer = setTimeout(onClose, duration);
        return () => clearTimeout(timer);
    }, [message, onClose, duration]);

    if (!message) return null;

    return (
        <div
            style={{
                position: "fixed",
                top: 24,
                right: 24,
                zIndex: 9999,
                background: "#fff",
                color: "#d32f2f",
                border: "1px solid #f8d7da",
                borderRadius: 8,
                boxShadow: "0 2px 8px rgba(0,0,0,0.12)",
                padding: "16px 24px 16px 16px",
                display: "flex",
                alignItems: "center",
                minWidth: 240,
                maxWidth: 360,
                fontSize: 16,
                gap: 12,
                animation: "fadeInDown 0.3s"
            }}
        >
            <X
                size={20}
                style={{ cursor: "pointer", marginRight: 8 }}
                onClick={onClose}
                aria-label="Close"
            />
            <span style={{ flex: 1 }}>{message}</span>
        </div>
    );
};

export default ErrorToast;
