import type {
    FindByUserIdRequest,
    GoogleLoginRequest, GoogleRegisterRequest,
    LoginRequest,
    LoginResponse,
    RegisterRequest,
    RegisterResponse, User
} from "../generated/dto/galactus.ts";
// import type {User} from "../types/type.ts";
import {createContext, type FC, type ReactNode, useContext, useEffect, useState} from "react";
import {useGrpc} from "./ClientContext.tsx";
import {useNavigate} from "react-router-dom";

type UserAuthContext = {
    register: (req: RegisterRequest) => Promise<RegisterResponse>
    login: (req: LoginRequest) => Promise<LoginResponse>
    googleLogin: (req: GoogleLoginRequest) => Promise<LoginResponse>
    googleRegister: (req: GoogleRegisterRequest) => Promise<RegisterResponse>
    logout: () => void
    user: User | null
}

const AuthContext = createContext<UserAuthContext | undefined>(undefined)

export const AuthProvider: FC<{children: ReactNode}> = ({children}: { children: React.ReactNode }) => {
    const { galactus } = useGrpc()
    const [user, setUser] = useState<User | null>(null)

    useEffect(() => {
        const stored = localStorage.getItem("user");
        if (stored) {
            try {
                const parsed = JSON.parse(stored);
                if (parsed && typeof parsed === "object" && parsed.userId) {
                    galactus.findByUserId({ id: parsed.userId }).then(res => {
                        if (res.response.user) setUser(res.response.user);
                    });
                } else if (typeof parsed === "number") {
                    galactus.findByUserId({ id: parsed }).then(res => {
                        if (res.response.user) setUser(res.response.user);
                    });
                }
            } catch(error) {
                console.log(error);
            }
        }
    }, [galactus])

    const register = async (req: RegisterRequest) => {
        try {
            const res = await galactus.register(req)
            return res.response
        }
        catch(error) {
            throw error
        }
    }
    const login = async (req: LoginRequest) => {
        try {
            const res = await galactus.login(req);
            const userRes: FindByUserIdRequest= {
                id: res.response.userId
            }
            const userdata = await galactus.findByUserId(userRes)
            if(userdata.response.user) {
                const safeUser = JSON.parse(JSON.stringify(userdata.response.user, (key, value) =>
                    typeof value === 'bigint' ? Number(value) : value
                ));
                setUser(safeUser)
                localStorage.setItem("user", JSON.stringify(safeUser))
                console.log(userdata.response.user)
                console.log(user)
            }
            return res.response;
        } catch (error) {
            console.error("Login error:", error);
            throw error;
        }
    };

    const googleLogin = async (req: GoogleLoginRequest) => {
        try {
            const res = await galactus.googleLogin(req);
            const userRes: FindByUserIdRequest= {
                id: res.response.userId
            }
            console.log(userRes)
            const userdata = await galactus.findByUserId(userRes)
            if(userdata.response.user) {
                const safeUser = JSON.parse(JSON.stringify(userdata.response.user, (key, value) =>
                    typeof value === 'bigint' ? Number(value) : value
                ));
                setUser(safeUser)
                localStorage.setItem("user", JSON.stringify(safeUser))
                console.log(userdata)
                console.log(safeUser)
            }
            return res.response;
        } catch (error) {
            console.error("Google login error:", error);
            throw error;
        }
    };

    const googleRegister = async (req: GoogleRegisterRequest) => {
        try {
            const res = await galactus.googleRegister(req);
            return res.response;
        } catch (error) {
            console.error("Google register error:", error);
            throw error;
        }
    };

    const logout = () => {
        localStorage.removeItem("user");
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{ register, login, googleLogin, googleRegister, logout, user }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = (): UserAuthContext => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
};
