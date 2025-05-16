import { request } from './request';
// Import specific types
import type { LoginResponseData, RefreshTokenResponseData } from '../types/api';

/**
 * Logs in the user. Attempts to use WeChat's code for login first.
 * If wx.login fails or the backend call with the code fails,
 * it falls back to using a mock OpenID.
 * Stores both access_token and refresh_token upon successful login.
 * @returns Promise<LoginResponseData> The login response data.
 */
export const login = (): Promise<LoginResponseData> => {
  return new Promise((resolve, reject) => {
    // 1. Attempt to get the real WeChat login code
    uni.login({
      provider: 'weixin',
      success: (loginRes) => {
        console.log('uni.login success, code:', loginRes.code);
        // 2. Send the code to the backend
        request<LoginResponseData>({
          url: `/login?code=${loginRes.code}`, // Send code as query parameter
          method: 'POST',
        })
          .then(data => handleLoginResponse(data, resolve, reject))
          .catch(error => {
            console.error('Backend login with code failed:', error);
            // 3. Fallback to mock OpenID if backend login fails
            loginWithMockIdInternal(resolve, reject);
          });
      },
      fail: (err) => {
        console.error('uni.login failed:', err);
        // 3. Fallback to mock OpenID if uni.login fails
        loginWithMockIdInternal(resolve, reject);
      },
    });
  });
};

// Internal function to handle login with mock OpenID
const loginWithMockIdInternal = (resolve: (value: LoginResponseData | PromiseLike<LoginResponseData>) => void, reject: (reason?: any) => void) => {
  const mockOpenId = `mock_openid_${Date.now()}_${Math.random().toString(36).substring(2)}`;
  console.log('Using mock OpenID for login:', mockOpenId);

  request<LoginResponseData>({
    url: `/login?code=${mockOpenId}`, // Assuming the same endpoint can handle mock openid
    method: 'POST',
  })
    .then(data => handleLoginResponse(data, resolve, reject))
    .catch(error => {
      console.error('Login with mock OpenID failed:', error);
      reject(error); // Reject the main promise if mock login also fails
    });
};

// Helper function to process login response and store tokens
const handleLoginResponse = (data: LoginResponseData, resolve: (value: LoginResponseData | PromiseLike<LoginResponseData>) => void, reject: (reason?: any) => void) => {
  console.log('call func \nLogin response data:', data);
  if (data && data.access_token && data.refresh_token) {
    try {
      uni.setStorageSync('access_token', data.access_token);
      uni.setStorageSync('refresh_token', data.refresh_token); // Store refresh token
      console.log('Tokens stored successfully.');
      resolve(data);
    } catch (e) {
      console.error('Failed to store tokens:', e);
      reject(new Error('Failed to store tokens'));
    }
  } else {
    console.error('Login response missing tokens:', data);
    // Don't reject here immediately if falling back from real code login,
    // let the fallback mechanism handle it. If this is the mock login failing,
    // then the final catch in loginWithMockIdInternal will reject.
    // However, if the backend *successfully* responds but without tokens, it's an error.
    reject(new Error('Invalid login response: missing tokens'));
  }
};

/**
 * Attempts to refresh the access token using the stored refresh token.
 * Updates the stored access token on success.
 * @returns Promise<string> The new access token.
 */
export const refreshAccessToken = (): Promise<string> => {
  return new Promise((resolve, reject) => {
    let refreshToken: string | null = null;
    let currentAccessToken: string | null = null; // Needed for the Authorization header

    try {
      refreshToken = uni.getStorageSync('refresh_token');
      currentAccessToken = uni.getStorageSync('access_token'); // Get current token for header
    } catch (e) {
      console.error('Failed to retrieve tokens from storage:', e);
      return reject(new Error('Failed to retrieve tokens'));
    }

    if (!refreshToken || !currentAccessToken) {
      console.warn('Missing refresh token or access token for refresh attempt.');
      return reject(new Error('Missing tokens for refresh'));
    }

    // Note: The backend refresh endpoint requires the *current* (potentially expired)
    // access token in the Authorization header and the refresh_token as a query parameter.
    request<RefreshTokenResponseData>({
      url: '/refresh',
      method: 'POST', // Assuming POST based on typical practices, adjust if needed
      data: { // Send refresh_token as part of the request body based on Go code
        refresh_token: refreshToken,
      },
      // The request function automatically adds the Authorization header with the current access_token
      // We might need a way to bypass the standard token check *within* the request function
      // for this specific call if the token is known expired, or ensure the backend
      // specifically allows expired access tokens for the refresh endpoint.
      // For now, assume the request function handles adding the current token header.
    })
      .then(data => {
        if (data && data.access_token) {
          try {
            uni.setStorageSync('access_token', data.access_token); // Update access token
            console.log('Access token refreshed successfully.');
            resolve(data.access_token);
          } catch (e) {
            console.error('Failed to store refreshed access token:', e);
            reject(new Error('Failed to store refreshed token'));
          }
        } else {
          console.error('Refresh token response missing access_token:', data);
          reject(new Error('Invalid refresh response: missing access_token'));
        }
      })
      .catch(error => {
        console.error('Refresh token request failed:', error);
        // If refresh fails (e.g., refresh token itself expired or invalid), clear tokens
        try {
          uni.removeStorageSync('access_token');
          uni.removeStorageSync('refresh_token');
          console.log('Cleared tokens due to refresh failure.');
          // TODO: Redirect to login page
          // uni.reLaunch({ url: '/pages/login/login' }); // Example redirect
        } catch (e) {
          console.error('Failed to clear tokens after refresh failure:', e);
        }
        reject(error); // Reject the promise so the original request knows refresh failed
      });
  });
};

// Example usage (optional, can be called from a page like index.vue)
/*
login() // Use the new login function
  .then(loginData => {
    console.log('Login successful:', loginData);
    // Navigate to another page or update UI
  })
  .catch(error => {
    console.error('Login failed:', error);
    // Show error message to the user
  });
*/