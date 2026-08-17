import "../style/login.css"
import {useNavigate} from "react-router-dom";
import {type CredentialResponse, GoogleLogin, GoogleOAuthProvider} from "@react-oauth/google";
import {useGrpc} from "../context/ClientContext.tsx";
import {jwtDecode} from "jwt-decode";
import type {
    GoogleRegisterRequest,
    LoginRequest,
    LoginResponse,
    RegisterRequest, SendOTPRequest
} from "../generated/galactus_controller.ts";
import type {GooglePayload} from "../types/type.ts";
import {Timestamp} from "../generated/google/protobuf/timestamp.ts";
import {useAuth} from "../context/AuthContext.tsx";
import {useState} from "react";



export const RegisterPage = () => {
    const navigate = useNavigate();
    const { googleRegister, register } = useAuth();
    const { galactus } = useGrpc()
    const [dateOfBirth, setDateOfBirth] = useState("")
    const [email, setEmail] = useState("")
    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [oTPCode, setOTPCode] = useState("")
    const googleRegisterHandler = async (cred: CredentialResponse )=> {
        try {
            const token = cred.credential
            const payload: GooglePayload = jwtDecode(token)
            const input: GoogleRegisterRequest = {
                username: payload.name,
                email: payload.email
            }
            const res = await googleRegister(input)
            console.log(res)
        }
        catch(error) {
            console.log(error)
        }
    }

    const registerHandler = async () => {
        try {
            let dateOfBirthTimestamp: Timestamp | undefined = undefined;
            if (dateOfBirth) {
                const date = new Date(dateOfBirth);
                dateOfBirthTimestamp = {
                    seconds: Math.floor(date.getTime() / 1000),
                    nanos: (date.getTime() % 1000) * 1e6,
                };
            }
            const input: RegisterRequest = {
                dateOfBirth: dateOfBirthTimestamp,
                email,
                oTPCode,
                password,
                username
            }
            const res = await register(input)
            navigate("/login")
            console.log(res)
        }
        catch (error) {
            console.log(error)
        }
    }

    const handleSendRegisterOtp = async () => {
        try {
            if (!email) {
                throw "Email must be filled"
            }
            const input: SendOTPRequest =  {
                email: email
            }
            const res = await galactus.sendOTP(input)
            console.log(res.response)
        }
        catch (error) {
            console.log(error)
        }
    }

    return (
        <div className="tiktok-container">
            <main className="main-content">
                <div className="signup-container">
                    <h1 className="title">Sign up for TikTok</h1>
                    <p className="subtitle">
                        Register Your Profile!
                    </p>

                    <div className="signup-options">
                        <form className="login-form" onSubmit={e => {
                            e.preventDefault()
                            registerHandler()
                        }}>
                            <input
                                type="email"
                                className="input-form"
                                placeholder="Email"
                                value={email}
                                onChange={e => setEmail(e.target.value)}
                                required
                            />
                            <input
                                type="text"
                                className="input-form"
                                placeholder="Username"
                                value={username}
                                onChange={e => setUsername(e.target.value)}
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
                            <input
                                type="date"
                                className="input-form"
                                placeholder="Date of Birth"
                                value={dateOfBirth}
                                onChange={e => setDateOfBirth(e.target.value)}
                                required
                            />
                            <div className="otp-row">
                                <input
                                    type="text"
                                    className="input-form otp-input"
                                    placeholder="OTP Code"
                                    value={oTPCode}
                                    onChange={e => setOTPCode(e.target.value)}
                                    required
                                />
                                <button
                                    type="button"
                                    className="login-btn"
                                    onClick={handleSendRegisterOtp}
                                >
                                    Send OTP
                                </button>
                            </div>
                            <button type="submit" className="login-btn">Register</button>
                        </form>
                        <GoogleOAuthProvider clientId="377333372407-4i4f9it1sm0ci0djskkboe6n15ds6qce.apps.googleusercontent.com">
                            <GoogleLogin onSuccess={(res) => {
                                googleRegisterHandler(res)
                            }} onError={() => {
                                console.log("error")
                            }}/>
                        </GoogleOAuthProvider>
                    </div>

                    <div className="login-link">
                        Already have an account? <a onClick={() => {
                        navigate("/login")
                    }}>Login</a>
                    </div>
                </div>
            </main>
        </div>
    )
}