// src/api/request.ts
import type { ApiResponse } from '../types/api.d.ts';
// Import the actual constants from the new location
import { ApiErrorCodeTokenExpired, ApiErrorCodeAuthError } from '../constants/api';
// Import the refresh function, ensuring no circular dependency issues
// If 'auth' imports 'request', we might need to move refreshAccessToken elsewhere or restructure.
// Assuming it's safe for now.
import { refreshAccessToken } from './auth';

// Read Base URL from environment variable
const BASE_URL = import.meta.env.VITE_API_BASE_URL;

if (!BASE_URL) {
  console.error("Error: VITE_API_BASE_URL is not defined in your environment variables (.env.* file).");
  // Optionally throw an error or set a default, though erroring is safer
  // throw new Error("API Base URL is not configured.");
}

interface RequestOptions extends UniApp.RequestOptions {
  // Flag to prevent retrying refresh token request itself
  _isRetry?: boolean;
  _isRefreshing?: boolean; // Internal flag to track refresh state
  /**
   * Query parameters to be appended to the URL
   */
  params?: Record<string, any>;
}

// --- Token Management ---
let isRefreshing = false; // Flag to prevent multiple concurrent refresh attempts
let pendingRequests: Array<(token: string) => void> = []; // Queue for requests waiting for token refresh

const getToken = (): string | null => {
  try {
    const token = uni.getStorageSync('access_token');
    return token ? `Bearer ${token}` : null;
  } catch (e) {
    console.error('Error getting token from storage:', e);
    return null;
  }
};

const clearTokensAndRedirect = () => {
  try {
    uni.removeStorageSync('access_token');
    uni.removeStorageSync('refresh_token');
    console.log('Tokens cleared, redirecting to login.');
    // TODO: Implement actual redirection logic
    // Example: uni.reLaunch({ url: '' });
  } catch (e) {
    console.error('Failed to clear tokens or redirect:', e);
  }
};

// Function to process the pending queue after token refresh
const processPendingRequests = (newToken: string | null, error?: any) => {
  pendingRequests.forEach(callback => {
    if (newToken) {
      callback(newToken);
    } else {
      // If refresh failed, we might want to reject the pending requests instead
      // For simplicity now, let them retry and potentially fail again or trigger login
      // A better approach might involve passing the error to the callbacks
       console.error("Refresh failed, cannot resolve pending request.", error);
       // Consider rejecting promises associated with these callbacks if possible
    }
  });
  pendingRequests = []; // Clear the queue
};

/**
 * Encapsulated request function with automatic token refresh.
 */
export const request = <T = any>(options: RequestOptions): Promise<T> => {
  return new Promise((resolve, reject) => {
    // If a refresh is already in progress, queue this request
    if (isRefreshing && !options._isRefreshing) { // Don't queue the refresh request itself
        console.log(`[Request Queue] ${options.method} ${options.url} waiting for token refresh.`);
        pendingRequests.push((newToken) => {
            console.log(`[Request Retry] ${options.method} ${options.url} retrying after refresh.`);
            options.header = { ...options.header, Authorization: `Bearer ${newToken}` };
            // Re-run the request logic for the queued request
            makeRequest<T>(options).then(resolve).catch(reject);
        });
        return; // Don't proceed with the original request yet
    }

    // Proceed with making the actual request
    makeRequest<T>(options).then(resolve).catch(async (error) => {
        // Check if the error indicates an expired token and it's not a retry/refresh attempt
        const apiResponse = error?.response?.data as ApiResponse<any>; // Extract response if available
        const statusCode = error?.response?.statusCode; // HTTP status if available

        // Check for TokenExpired business code OR 401 Unauthorized HTTP status
        // AND ensure this isn't already a retry after a refresh attempt
        if (
            !options._isRetry &&
            ( (apiResponse?.code === ApiErrorCodeTokenExpired) || (statusCode === 401 && apiResponse?.code === ApiErrorCodeAuthError) ) // Backend uses 401 + AuthError for expired token
           )
        {
            console.warn(`[Token Refresh] Detected expired token for ${options.method} ${options.url}. Attempting refresh.`);

            // Prevent multiple refresh attempts
            if (!isRefreshing) {
                isRefreshing = true;
                options._isRefreshing = true; // Mark this specific call chain part of refresh

                try {
                    const newToken = await refreshAccessToken(); // Attempt refresh
                    console.log('[Token Refresh] Success.');
                    isRefreshing = false;
                    processPendingRequests(newToken); // Process queued requests with the new token

                    // Retry the original request with the new token
                    console.log(`[Request Retry] ${options.method} ${options.url} retrying with new token.`);
                    options.header = { ...options.header, Authorization: `Bearer ${newToken}` };
                    options._isRetry = true; // Mark as retry to prevent infinite loop
                    makeRequest<T>(options).then(resolve).catch(reject); // Retry the original request

                } catch (refreshError) {
                    console.error('[Token Refresh] Failed:', refreshError);
                    isRefreshing = false;
                    processPendingRequests(null, refreshError); // Notify queue about failure
                    clearTokensAndRedirect(); // Clear tokens and redirect on refresh failure
                    reject(refreshError); // Reject the original request's promise
                }
            } else {
                 // If refresh is already happening, queue this request (should have been caught above, but as fallback)
                 console.log(`[Request Queue] ${options.method} ${options.url} joining wait for ongoing refresh.`);
                 pendingRequests.push((newToken) => {
                     console.log(`[Request Retry] ${options.method} ${options.url} retrying after refresh.`);
                     options.header = { ...options.header, Authorization: `Bearer ${newToken}` };
                     options._isRetry = true;
                     makeRequest<T>(options).then(resolve).catch(reject);
                 });
            }
        } else {
            // Not a token expiry error, or it's already a retry, reject directly
            reject(error);
        }
    });
  });
};


/**
 * Internal function to make the actual uni.request call.
 */
const makeRequest = <T = any>(options: RequestOptions): Promise<T> => {
    return new Promise((resolve, reject) => {
        // Ensure BASE_URL is defined before proceeding
        if (!BASE_URL) {
            return reject(new Error("API Base URL is not configured. Check environment variables."));
        }
        const url = BASE_URL + options.url;
        const method = options.method || 'GET';
        const header = {
            'Content-Type': 'application/json',
            ...options.header,
        };

        // Add Authorization header if token exists and not explicitly removed
        if (!header.hasOwnProperty('Authorization') || header['Authorization'] !== null) {
            const token = getToken();
            if (token) {
                header['Authorization'] = token;
            }
        }
         // Remove null Authorization header if it was set explicitly to null
        if (header['Authorization'] === null) {
            delete header['Authorization'];
        }


        console.log(`[Request] ${method} ${url}`, { data: options.data, params: options.params, header: header }); // Added params logging

        uni.request({
            ...options,
            url,
            method,
            header,
            success: (res) => {
                console.log(`[Response] ${method} ${url}`, { statusCode: res.statusCode, data: res.data });

                // **Moved the declaration and check outside the main success/error logic**
                const apiResponse = res.data as ApiResponse<T>; // Declare and cast ONCE
                const statusCode = res.statusCode;

                // Check for token expiration or specific auth error indicating expiration
                // This check is primarily for logging/debugging here, the main refresh logic is in the outer catch
                if (
                    ( (apiResponse?.code === ApiErrorCodeTokenExpired) || (statusCode === 401 && apiResponse?.code === ApiErrorCodeAuthError) )
                ) {
                  console.log('[Debug] Token might be expired based on this response.');
                }

                // Main success/error handling
                if (statusCode !== 200) {
                    uni.showToast({ title: `网络错误 (${statusCode})`, icon: 'none' });
                    return reject({ message: `HTTP Error ${statusCode}`, response: res });
                }

                // Use the already declared apiResponse
                if (apiResponse.code === 200) {
                    resolve(apiResponse.data);
                } else {
                    console.error(`Business Error ${apiResponse.code}: ${apiResponse.msg}`);
                    uni.showToast({ title: apiResponse.msg || '请求失败', icon: 'none' });
                    reject({ message: apiResponse.msg || `Business Error ${apiResponse.code}`, response: res });
                }
            },
            fail: (err) => {
                console.error(`[Request Fail] ${method} ${url}`, err);
                uni.showToast({ title: '请求失败，请检查网络', icon: 'none' });
                reject(err); // Reject with the original error object
            },
        });
    });
};

// Example usage (can be moved to specific api files like src/api/portfolio.ts)
/*
import type { GetAllTemplatesResponse } from '../types/api';

export const fetchAllTemplates = () => {
  return request<GetAllTemplatesResponse['data']>({ // Specify the expected data type
    url: '/portfolio/template/all',
    method: 'GET',
  });
};

fetchAllTemplates()
  .then(data => {
    console.log('Fetched templates:', data);
    // data is now typed as GetAllTemplatesResponseItem[]
  })
  .catch(error => {
    console.error('Failed to fetch templates:', error);
  });
*/
