import React, { useState, useEffect } from "react";
import { useGrpc } from "../context/ClientContext";
import { useNavigate } from "react-router-dom";

interface RichTextProps {
    value?: string;
    placeholder?: string;
    onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

const RichText: React.FC<RichTextProps> = ({ value = "", placeholder = "Type @username or #hashtag...", onChange }) => {
    const [input, setInput] = useState(value);
    const [error, setError] = useState<string | null>(null);
    const [mentionUserId, setMentionUserId] = useState<string | null>(null);
    const { galactus } = useGrpc();
    const navigate = useNavigate();

    useEffect(() => {
        setError(null);
        setMentionUserId(null);
        if (input.startsWith("@")) {
            const username = input.slice(1).split(" ")[0];
            if (username.length > 0) {
                galactus.findByUsername({ username }).then(res => {
                    const userData = res.response.user;
                    if (userData && userData.userId) {
                        setMentionUserId(userData.userId.toString());
                    } else {
                        setError("User not found.");
                    }
                }).catch(() => setError("User not found."));
            }
        }
    }, [input, galactus]);

    useEffect(() => {
        setInput(value);
    }, [value]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setInput(e.target.value);
        if (onChange) onChange(e);
    };

    const handleMentionClick = () => {
        if (mentionUserId) {
            navigate(`/profile/${mentionUserId}`);
        }
    };

    return (
        <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
            <input
                type="text"
                value={input}
                onChange={handleChange}
                placeholder={placeholder}
                style={{
                    padding: "8px 12px",
                    borderRadius: 6,
                    border: "1px solid #ccc",
                    fontSize: 16,
                }}
            />
            {input.startsWith("@") && mentionUserId && (
                <div
                    style={{
                        color: "#fe2c55",
                        fontWeight: 600,
                        cursor: "pointer",
                        textDecoration: "underline"
                    }}
                    onClick={handleMentionClick}
                >
                    Mention: {input}
                </div>
            )}
            {input.startsWith("#") && (
                <div style={{ color: "#fe2c55", fontWeight: 600 }}>
                    Hashtag: {input}
                </div>
            )}
            {error && (
                <div style={{ color: "red", fontSize: 14 }}>
                    {error}
                </div>
            )}
        </div>
    );
};

export default RichText;