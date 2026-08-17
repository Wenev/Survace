import "../style/login.css"
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useGrpc } from "../context/ClientContext";
import type {
  ForgetPasswordRequest,
  ResetPasswordRequest,
} from "../generated/galactus_controller";
import ErrorToast from "../component/ErrorToast";

export default function ForgetPassword() {
  const { galactus } = useGrpc();
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [otpCode, setOtpCode] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [step, setStep] = useState<"request" | "reset">("request");
  const [error, setError] = useState<string | null>(null);
  const [info, setInfo] = useState<string | null>(null);

  const handleSendOtp = async () => {
    try {
      if (!email) throw new Error("Email must be filled");
      const input: ForgetPasswordRequest = { email };
      const res = await galactus.forgetPassword(input);
      if (res?.response?.statusCode === 0) {
        setInfo("OTP sent to your email.");
        setStep("reset");
      } else {
        setError(res?.response?.message || "Failed to send OTP");
      }
    } catch (e: any) {
      setError(e?.message || "Failed to send OTP");
    }
  };

  const handleResetPassword = async () => {
    try {
      if (!email || !otpCode || !newPassword) throw new Error("All fields are required");
      const input: ResetPasswordRequest = {
        email,
        otpCode,
        newPassword,
      };
      const res = await galactus.resetPassword(input);
      console.log(email)
      if (res?.response?.statusCode === 0) {
        setInfo("Password reset successful. Please login.");
        setTimeout(() => navigate("/login"), 1200);
      } else {
        setError(res?.response?.message || "Failed to reset password");
      }
    } catch (e: any) {
      setError(e?.message || "Failed to reset password");
    }
  };

  return (
    <div className="tiktok-container">
      <main className="main-content">
        <ErrorToast message={error || ""} onClose={() => setError(null)} />
        <div className="signup-container">
          <h1 className="title">Forgot Password</h1>
          <p className="subtitle">
            Reset your password using your email and OTP.
          </p>
          <div className="signup-options">
            {step === "request" && (
              <form
                className="login-form"
                onSubmit={e => {
                  e.preventDefault();
                  handleSendOtp();
                }}
              >
                <input
                  type="email"
                  className="input-form"
                  placeholder="Email"
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                  required
                />
                <button type="submit" className="login-btn">
                  Send OTP
                </button>
              </form>
            )}
            {step === "reset" && (
              <form
                className="login-form"
                onSubmit={e => {
                  e.preventDefault();
                  handleResetPassword();
                }}
              >
                <input
                  type="email"
                  className="input-form"
                  placeholder="Email"
                  value={email}
                  disabled
                  required
                />
                <input
                  type="text"
                  className="input-form otp-input"
                  placeholder="OTP Code"
                  value={otpCode}
                  onChange={e => setOtpCode(e.target.value)}
                  required
                />
                <input
                  type="password"
                  className="input-form"
                  placeholder="New Password"
                  value={newPassword}
                  onChange={e => setNewPassword(e.target.value)}
                  required
                />
                <button type="submit" className="login-btn">
                  Reset Password
                </button>
              </form>
            )}
            {info && (
              <div style={{ color: "#16a34a", marginTop: 8, fontSize: 14 }}>
                {info}
              </div>
            )}
          </div>
          <div className="login-link" style={{ marginTop: 16 }}>
            <span
              style={{
                textDecoration: "underline",
                cursor: "pointer",
                fontSize: 14,
              }}
              onClick={() => navigate("/login")}
            >
              Back to Login
            </span>
          </div>
        </div>
      </main>
    </div>
  );
}
