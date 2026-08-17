export type GooglePayload = {
    sub: string;
    email: string;
    email_verified: boolean;
    name: string;
    picture: string;
    given_name: string;
    family_name: string;
}

export type User = {
    id: number;
    username: string;
    email: string;
    dateOfBirth: Date | null;
    avatarUrl: string;
}

export interface Video {
    id: number;
    userId: number;
    title: string;
    url: string;
    objectName: string;
    enableComment: boolean;
    visibility: string;
    createdAt: { seconds: number; nanos?: number };
    isPlaying: boolean;
    likeCount?: number;
    commentCount?: number;
}
