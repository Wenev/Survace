import { createBrowserRouter } from "react-router-dom"
import {LoginPage} from "../pages/LoginPage.tsx";
import {RegisterPage} from "../pages/RegisterPage.tsx";
import { UploadPage } from "../pages/UploadPage.tsx";
import HomePage from "../pages/HomePage.tsx";
import ChatPage from "../pages/ChatPage";
import GoLivePage from "../pages/GoLivePage.tsx";
import LivePage from "../pages/LivePage";
import ProfilePage from "../pages/ProfilePage.tsx";
import FriendPage from "../pages/FriendPage.tsx";
import FollowingPage from "../pages/FollowingPage.tsx";
import SettingPage from "../pages/SettingPage.tsx";
import ForgetPassword from "../pages/ForgetPassword";
import StudioPage from "../pages/StudioPage.tsx";
import CreatePlaylist from "../pages/CreatePlaylist.tsx";
import WatchLive from "../pages/WatchLive.tsx";

export const routes = createBrowserRouter([
    {
      path: "/",
      element: <HomePage />,
    },
    {
        path: "/signup",
        element: <RegisterPage />,
        children: [
            // {
            //     index: true,
            //     element: <HomePage />,
            // },
            // {
            //     path: "signup",
            //     element: <SignupPage />,
            // },
        ],
    },
    {
        path: "/login",
        element: <LoginPage />,
    },
    {
        path: "/signup",
        element: <LoginPage />,
    },
    {
        path: "/upload",
        element: <UploadPage />,
    },
    {
        path: "/chat",
        element: <ChatPage />,
    },
    {
        path: "/live/streaming",
        element: <GoLivePage />,
    },
    {
        path: "/live",
        element: <LivePage />,
    },
    {
      path: "/live/:callId",
      element: <WatchLive />
    },
    {
        path: "/profile",
        element: <ProfilePage />,
    },
    {
        path: "/profile/:userId",
        element: <ProfilePage />,
    },
    {
        path: "/friends",
        element: <FriendPage />
    },
    {
        path: "/following",
        element: <FollowingPage />
    },
    {
        path: "/settings",
        element: <SettingPage />
    },
    {
        path: "/forgetpassword",
        element: <ForgetPassword />,
    },
    {
        path: "/studio",
        element: <StudioPage />
    },
    {
        path: "/playlist/create",
        element: <CreatePlaylist />
    },
    {
        path: "/playlist/:playlistId",
        // element: <EditPlaylist />
    }
])