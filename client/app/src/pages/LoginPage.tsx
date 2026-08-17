import "../style/login.css"
import {type CredentialResponse, GoogleLogin, GoogleOAuthProvider} from "@react-oauth/google";
import {useNavigate} from "react-router-dom";
import type {GooglePayload} from "../types/type.ts";
import {jwtDecode} from "jwt-decode";
import type {GoogleLoginRequest, LoginRequest, LoginResponse} from "../generated/galactus_controller.ts";
import {useAuth} from "../context/AuthContext.tsx";
import {useState} from "react";
import ErrorToast from "../component/ErrorToast";

export const LoginPage = () => {
    const { googleLogin, login } = useAuth()
    const navigate = useNavigate();
    const [usernameOrEmail, setUsernameOrEmail] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState<string | null>(null)
    const googleLoginHandler = async (cred: CredentialResponse) => {
        try {
            const token = cred.credential
            const payload: GooglePayload = jwtDecode(token)
            const input: GoogleLoginRequest = {
                email: payload.email
            }
            const res = await googleLogin(input)
            navigate("/upload")
            console.log(res)
        }
        catch(error: any) {
            setError(error?.message || "Google login failed");
            console.log(error)
        }
    }

    const loginHandler = async () => {
        try {
            const input: LoginRequest = {
                emailOrUsername: usernameOrEmail,
                password: password
            }
            const res = await login(input)
            navigate("/upload")
            console.log(res)
        }
        catch (error: any) {
            setError(error?.message || "Login failed");
            console.log(error)
        }
    }
    return (
            <div className="tiktok-container">
                <main className="main-content">
                    <ErrorToast message={error || ""} onClose={() => setError(null)} />
                    <div className="signup-container">
                        <h1 className="title">Login for TikTok</h1>
                        <div className="signup-options">
                            <form className="login-form" onSubmit={e => {
                                e.preventDefault()
                                loginHandler()
                            }}>
                                <input
                                    type="text"
                                    className="input-form"
                                    placeholder="Username or Email"
                                    value={usernameOrEmail}
                                    onChange={e => setUsernameOrEmail(e.target.value)}
                                    required
                                />
                                <input
                                    type="password"
                                    className="input-form"
                                    placeholder="Password"
                                    value={password}
                                    onChange={e => setPassword(e.target.value)}
                                    required
                                />
                                <button type="submit" className="login-btn">Login</button>
                            </form>
                            <GoogleOAuthProvider clientId="377333372407-4i4f9it1sm0ci0djskkboe6n15ds6qce.apps.googleusercontent.com">
                                <GoogleLogin onSuccess={(res: CredentialResponse) => {
                                    googleLoginHandler(res)
                                }} onError={() => {
                                    setError("Google login failed");
                                    console.log("error")
                                }}/>
                            </GoogleOAuthProvider>
                        </div>
                        <div className="login-link">
                            Don't have an account? <a onClick={() => {
                                navigate("/signup")
                            }}>Signup</a>
                        </div>
                        <div className="login-link" style={{marginTop: 8}}>
                            <span
                                style={{
                                    textDecoration: "underline",
                                    cursor: "pointer",
                                    fontSize: 14,
                                }}
                                onClick={() => navigate("/forgetpassword")}
                            >
                                Forgot password?
                            </span>
                        </div>
                    </div>
                </main>
            </div>
    )
}