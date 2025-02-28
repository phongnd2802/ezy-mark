import { ERROR_CODES } from "./errorCode";


export const ERROR_MESSAGES: Record<number, string> = {
    [ERROR_CODES.EMAIL_ALREADY_EXISTS]: "This email is already in use.",
    [ERROR_CODES.EXPIRED_SESSION]: "This OTP has expired. Tap 'Resend' to receive a new one.",
    [ERROR_CODES.OTP_DOES_NOT_MATCH]: "The OTP you entered is incorrect. Please try again.",
    [ERROR_CODES.INTERNAL_SERVER_ERROR]: "Something went wrong on our end. Please try again later.",
}   