// src/constants/api.ts

/**
 * API 业务错误码常量 (Based on backend code.md)
 */
export const ApiErrorCodeServerError = 0;
export const ApiErrorCodeAuthError = 1; // Generic auth error / No permission
export const ApiErrorCodeTokenExpired = 2;
export const ApiErrorCodeLoginError = 3;
export const ApiErrorCodeRefreshTokenError = 4; // Added refresh token error
// Add other business error codes as needed