import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import {GalactusControllerClient} from "../generated/controller/galactus_controller.client.ts";
import {createContext, type FC, useContext} from "react";
import {BrainrotServiceClient} from "../generated/controller/video_controller.client.ts";
import {SocialServiceClient} from "../generated/controller/social_controller.client.ts";
import {StreamServiceClient} from "../generated/controller/stream_controller.client.ts";
import { NotificationServiceClient } from "../generated/controller/notification_controller.client.ts";
import { QueryClient, QueryClientProvider, useQuery, useMutation } from "@tanstack/react-query";


const transport = new GrpcWebFetchTransport({
    baseUrl: "http://localhost:8080",
})

const galactusClient = new GalactusControllerClient(transport)
const brainrotClient = new BrainrotServiceClient(transport)
const chatClient = new SocialServiceClient(transport)
const streamClient = new StreamServiceClient(transport)
const notificationClient = new NotificationServiceClient(transport)

const queryClient = new QueryClient()

type GrpcClientContext = {
    galactus: GalactusControllerClient,
    brainrot: BrainrotServiceClient,
    social: SocialServiceClient,
    stream: StreamServiceClient,
    notification: NotificationServiceClient,
    queryClient: QueryClient,
    useVideoQuery: typeof useVideoQuery,
    useVideoMutation: typeof useVideoMutation,
    useGalactusQuery: typeof useGalactusQuery,
    useGalactusMutation: typeof useGalactusMutation,
    useSocialQuery: typeof useSocialQuery,
    useSocialMutation: typeof useSocialMutation,
    useStreamQuery: typeof useStreamQuery,
    useStreamMutation: typeof useStreamMutation,
    useNotificationQuery: typeof useNotificationQuery,
    useNotificationMutation: typeof useNotificationMutation,
}

// Generic query hook for any gRPC client
function useGrpcQuery<
    TClient,
    TKey extends keyof TClient,
    TReq extends object,
    TRes
>(
    client: TClient,
    key: TKey,
    req: TReq,
    options?: any
) {
    return useQuery({
        queryKey: [key, req],
        queryFn: async () => {
            // @ts-ignore
            const result = await client[key](req)
            return result
        },
        ...options
    })
}

// Generic mutation hook for any gRPC client
function useGrpcMutation<
    TClient,
    TKey extends keyof TClient,
    TReq extends object,
    TRes
>(
    client: TClient,
    key: TKey,
    options?: any
) {
    return useMutation({
        mutationFn: async (req: TReq) => {
            // @ts-ignore
            const result = await client[key](req)
            return result
        },
        ...options
    })
}

// Specific hooks for each controller
function useVideoQuery<TKey extends keyof BrainrotServiceClient, TReq extends object, TRes>(key: TKey, req: TReq, options?: any) {
    const { brainrot } = useGrpc()
    return useGrpcQuery(brainrot, key, req, options)
}
function useVideoMutation<TKey extends keyof BrainrotServiceClient, TReq extends object, TRes>(key: TKey, options?: any) {
    const { brainrot } = useGrpc()
    return useGrpcMutation(brainrot, key, options)
}

function useGalactusQuery<TKey extends keyof GalactusControllerClient, TReq extends object, TRes>(key: TKey, req: TReq, options?: any) {
    const { galactus } = useGrpc()
    return useGrpcQuery(galactus, key, req, options)
}
function useGalactusMutation<TKey extends keyof GalactusControllerClient, TReq extends object, TRes>(key: TKey, options?: any) {
    const { galactus } = useGrpc()
    return useGrpcMutation(galactus, key, options)
}

function useSocialQuery<TKey extends keyof SocialServiceClient, TReq extends object, TRes>(key: TKey, req: TReq, options?: any) {
    const { social } = useGrpc()
    return useGrpcQuery(social, key, req, options)
}
function useSocialMutation<TKey extends keyof SocialServiceClient, TReq extends object, TRes>(key: TKey, options?: any) {
    const { social } = useGrpc()
    return useGrpcMutation(social, key, options)
}

function useStreamQuery<TKey extends keyof StreamServiceClient, TReq extends object, TRes>(key: TKey, req: TReq, options?: any) {
    const { stream } = useGrpc()
    return useGrpcQuery(stream, key, req, options)
}
function useStreamMutation<TKey extends keyof StreamServiceClient, TReq extends object, TRes>(key: TKey, options?: any) {
    const { stream } = useGrpc()
    return useGrpcMutation(stream, key, options)
}

function useNotificationQuery<TKey extends keyof NotificationServiceClient, TReq extends object, TRes>(key: TKey, req: TReq, options?: any) {
    const { notification } = useGrpc()
    return useGrpcQuery(notification, key, req, options)
}
function useNotificationMutation<TKey extends keyof NotificationServiceClient, TReq extends object, TRes>(key: TKey, options?: any) {
    const { notification } = useGrpc()
    return useGrpcMutation(notification, key, options)
}

const ClientContext = createContext<GrpcClientContext | undefined>(undefined)

export const ClientProvider: FC<{children: React.ReactNode}> = ({children}) => {
    return (
        <QueryClientProvider client={queryClient}>
            <ClientContext.Provider value={{
                galactus: galactusClient,
                brainrot: brainrotClient,
                social: chatClient,
                stream: streamClient,
                notification: notificationClient,
                queryClient,
                useVideoQuery,
                useVideoMutation,
                useGalactusQuery,
                useGalactusMutation,
                useSocialQuery,
                useSocialMutation,
                useStreamQuery,
                useStreamMutation,
                useNotificationQuery,
                useNotificationMutation
            }}>
                {children}
            </ClientContext.Provider>
        </QueryClientProvider>
    )
}

export const useGrpc = (): GrpcClientContext => {
    const context = useContext(ClientContext)
    if (context === undefined) {
        throw new Error("useGrpc must be used within a ClientProvider")
    }
    return context
}