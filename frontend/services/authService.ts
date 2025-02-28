import apiClient from "@/utils/apiClient";
import { ApiResponse } from "@/utils/apiResponse";


interface SignUpRequest{
    email: string;
    password: string;
}

interface SignUpData{
    token: string;
}


interface VerifyOTPRequest{
    token: string;
    otp: string;
}

interface LoginRequest {
    email: string;
    password: string;
}

interface LoginResponse {
    access_token: string;
    refresh_token: string;
}

export const authService = {
    // Login
    logIn: async(params: LoginRequest) => {
        try {
            const res = await apiClient.post<ApiResponse<LoginResponse>>(`/auth/login`, params);
            return res.data;
        } catch (error) {
            throw error
        }
    },

    // Register
    signUp: async (params: SignUpRequest) => {
        try {
            const res = await apiClient.post<ApiResponse<SignUpData>>(`/auth/signup`, params);
            return res.data;
        } catch (error) {
            console.log(error);
            throw error;
        }
    },

    // Verity OTP
    verifyOTP: async (params: VerifyOTPRequest) => {
        try {
            const res = await apiClient.post<ApiResponse<null>>(`/auth/verify-otp`, params);
            return res.data;
        } catch (error) {
            throw error
        }
    },

    // GetTTLOTP
    getTTLOtp: async (token: string) => {
        try {
            const res = await apiClient.get<ApiResponse<{ttl: string}>>(`/auth/verify-otp?token=${token}`)
            return res.data;
        } catch (error) {
            throw error
        }
    }
}