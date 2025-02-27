import { ERROR_CODES } from "./errorCode";


export const ERROR_MESSAGES: Record<number, string> = {
    [ERROR_CODES.EMAIL_ALREADY_EXISTS]: "Email already exists",
    [ERROR_CODES.INTERNAL_SERVER_ERROR]: "Internal Server Error",
    
}