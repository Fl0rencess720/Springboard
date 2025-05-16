import OSS from 'ali-oss';
import { request } from './request'; // Assuming your request wrapper handles API calls
import type { Work } from '../types/api'; // Assuming Work type is defined

// Read OSS config from environment variables
const OSS_REGION = import.meta.env.VITE_OSS_REGION;
const OSS_BUCKET = import.meta.env.VITE_OSS_BUCKET;

if (!OSS_REGION || !OSS_BUCKET) {
  console.error("Error: VITE_OSS_REGION or VITE_OSS_BUCKET is not defined in your environment variables.");
  // Optionally throw an error
  // throw new Error("OSS Region or Bucket is not configured.");
}

interface OssCredentials {
  AccessKeyId: string;
  AccessKeySecret: string;
  SecurityToken: string;
  // Expiration is not provided by the backend in code.md,
  // STS tokens from AssumeRole typically last 1 hour by default.
  // We'll fetch new ones when needed or if operations fail.
  // Adding a fetched time for basic expiry check.
  fetchedAt?: number;
}

let cachedCredentials: OssCredentials | null = null;
let credentialPromise: Promise<OssCredentials> | null = null;

// Function to fetch STS credentials from the backend
const fetchOssCredentials = async (): Promise<OssCredentials> => {
  console.log('[OSS] Fetching new STS credentials...');
  try {
    // Use the request wrapper which handles base URL and potential errors
    const creds = await request<OssCredentials>({
      url: '/oss/sts', // Matches controller/oss.go -> GetCredentials
      method: 'GET',
    });
    cachedCredentials = { ...creds, fetchedAt: Date.now() };
    console.log('[OSS] STS credentials fetched successfully.');
    return cachedCredentials;
  } catch (error) {
    console.error('[OSS] Failed to fetch STS credentials:', error);
    cachedCredentials = null; // Clear cache on failure
    throw error; // Re-throw the error to be handled by the caller
  }
};

// Function to get valid credentials, fetching or returning cached ones
const getValidCredentials = async (): Promise<OssCredentials> => {
  const now = Date.now();
  // Check if cached credentials exist and are less than 55 minutes old (giving 5 min buffer)
  if (cachedCredentials && cachedCredentials.fetchedAt && (now - cachedCredentials.fetchedAt < 55 * 60 * 1000)) {
    console.log('[OSS] Using cached STS credentials.');
    return cachedCredentials;
  }

  // If fetching is already in progress, return the existing promise
  if (credentialPromise) {
    console.log('[OSS] Waiting for ongoing credential fetch...');
    return credentialPromise;
  }

  // Fetch new credentials
  credentialPromise = fetchOssCredentials();
  try {
    const creds = await credentialPromise;
    return creds;
  } finally {
    // Reset promise regardless of success or failure
    credentialPromise = null;
  }
};

// Function to get an initialized OSS client instance
const getOssClient = async (): Promise<OSS> => {
  if (!OSS_REGION || !OSS_BUCKET) {
    throw new Error("OSS Region or Bucket is not configured.");
  }
  try {
    const credentials = await getValidCredentials();
    const client = new OSS({
      region: OSS_REGION,
      bucket: OSS_BUCKET,
      accessKeyId: credentials.AccessKeyId,
      accessKeySecret: credentials.AccessKeySecret,
      stsToken: credentials.SecurityToken,
      refreshSTSToken: async () => {
        // This function is called by the SDK if it detects the token might be expired
        // or needs refreshing based on its internal logic or timing.
        console.log('[OSS] SDK requesting STS token refresh...');
        // Force fetch new credentials when SDK requests refresh
        cachedCredentials = null; // Invalidate cache
        return await getValidCredentials();
      },
      // refreshSTSTokenInterval: 300000, // Refresh every 5 minutes (optional, refreshSTSToken is often sufficient)
      secure: true, // Use HTTPS
    });
    console.log('[OSS] Client initialized/refreshed.');
    return client;
  } catch (error) {
    console.error('[OSS] Failed to initialize OSS client:', error);
    throw error;
  }
};

/**
 * Uploads a file (Blob or File) to OSS.
 * @param file The file or blob to upload.
 * @param objectKey The full path and name for the object in OSS (e.g., 'user/avatars/image.jpg').
 * @returns The result from OSS client.put.
 */
export const uploadFile = async (file: Blob | File, objectKey: string): Promise<OSS.PutObjectResult> => {
  console.log(`[OSS Upload] Starting upload for key: ${objectKey}`);
  try {
    const client = await getOssClient();
    const result = await client.put(objectKey, file);
    console.log(`[OSS Upload] Success for key: ${objectKey}`, result);
    return result;
  } catch (error) {
    console.error(`[OSS Upload] Failed for key: ${objectKey}`, error);
    // Attempt to refresh credentials and retry once if it's likely an auth error
    if ((error as any)?.name === 'RequestTimeTooSkewed' || (error as any)?.status === 403) {
        console.warn(`[OSS Upload] Retrying upload for ${objectKey} after potential auth error.`);
        cachedCredentials = null; // Force refresh
        try {
            const client = await getOssClient();
            const result = await client.put(objectKey, file);
            console.log(`[OSS Upload] Retry Success for key: ${objectKey}`, result);
            return result;
        } catch (retryError) {
             console.error(`[OSS Upload] Retry Failed for key: ${objectKey}`, retryError);
             throw retryError; // Throw the retry error
        }
    }
    throw error; // Throw original error if not an auth-like error or retry failed
  }
};

/**
 * Generates a signed URL for previewing an object in the browser.
 * Default expiry is 1 hour (3600 seconds).
 * @param objectKey The full path and name of the object in OSS.
 * @param options Optional parameters for signatureUrl (e.g., process for image processing).
 * @returns The signed URL string.
 */
export const getPreviewUrl = async (objectKey: string, options?: OSS.SignatureUrlOptions): Promise<string> => {
  console.log(`[OSS Preview URL] Generating for key: ${objectKey}`);
   if (!objectKey) {
    console.warn('[OSS Preview URL] Received empty objectKey, returning empty string.');
    return '';
  }
  try {
    const client = await getOssClient();
    // Ensure Content-Disposition is not set to attachment for preview
    const defaultOptions: OSS.SignatureUrlOptions = { expires: 3600 }; // 1 hour default expiry
    const finalOptions = { ...defaultOptions, ...options };
    // Remove content-disposition if accidentally passed for preview
    if (finalOptions.response?.['content-disposition']) {
        console.warn('[OSS Preview URL] Removing content-disposition for preview URL.');
        delete finalOptions.response['content-disposition'];
        if (Object.keys(finalOptions.response).length === 0) {
            delete finalOptions.response;
        }
    }

    const url = client.signatureUrl(objectKey, finalOptions);
    console.log(`[OSS Preview URL] Generated for key: ${objectKey}`);
    return url;
  } catch (error) {
    console.error(`[OSS Preview URL] Failed for key: ${objectKey}`, error);
     // Attempt to refresh credentials and retry once if it's likely an auth error
    if ((error as any)?.name === 'RequestTimeTooSkewed' || (error as any)?.status === 403) {
        console.warn(`[OSS Preview URL] Retrying generation for ${objectKey} after potential auth error.`);
        cachedCredentials = null; // Force refresh
        try {
            const client = await getOssClient();
            const url = client.signatureUrl(objectKey, options);
            console.log(`[OSS Preview URL] Retry Success for key: ${objectKey}`);
            return url;
        } catch (retryError) {
             console.error(`[OSS Preview URL] Retry Failed for key: ${objectKey}`, retryError);
             throw retryError;
        }
    }
    throw error;
  }
};

/**
 * Generates a signed URL for downloading an object.
 * Sets Content-Disposition to attachment.
 * Default expiry is 1 hour (3600 seconds).
 * @param objectKey The full path and name of the object in OSS.
 * @param filename Optional desired filename for the download. Defaults to the objectKey basename.
 * @param options Optional parameters for signatureUrl.
 * @returns The signed URL string.
 */
export const getDownloadUrl = async (objectKey: string, filename?: string, options?: OSS.SignatureUrlOptions): Promise<string> => {
  console.log(`[OSS Download URL] Generating for key: ${objectKey}`);
   if (!objectKey) {
    console.warn('[OSS Download URL] Received empty objectKey, returning empty string.');
    return '';
  }
  try {
    const client = await getOssClient();
    const downloadFilename = filename || objectKey.split('/').pop() || 'download';
    const responseHeaders = {
      'content-disposition': `attachment; filename=${encodeURIComponent(downloadFilename)}`,
      // Add other headers if needed, e.g., 'content-type'
    };
    const defaultOptions: OSS.SignatureUrlOptions = {
        expires: 3600, // 1 hour default expiry
        response: responseHeaders
    };
    // Merge options, ensuring response headers are preserved/merged correctly
    const finalOptions = { ...defaultOptions, ...options };
    finalOptions.response = { ...responseHeaders, ...options?.response }; // Ensure our disposition is set

    const url = client.signatureUrl(objectKey, finalOptions);
    console.log(`[OSS Download URL] Generated for key: ${objectKey}`);
    return url;
  } catch (error) {
    console.error(`[OSS Download URL] Failed for key: ${objectKey}`, error);
     // Attempt to refresh credentials and retry once if it's likely an auth error
    if ((error as any)?.name === 'RequestTimeTooSkewed' || (error as any)?.status === 403) {
        console.warn(`[OSS Download URL] Retrying generation for ${objectKey} after potential auth error.`);
        cachedCredentials = null; // Force refresh
        try {
            const client = await getOssClient();
            // Reconstruct options for retry
            const downloadFilename = filename || objectKey.split('/').pop() || 'download';
            const responseHeaders = { 'content-disposition': `attachment; filename=${encodeURIComponent(downloadFilename)}` };
            const finalOptions = { expires: 3600, response: responseHeaders, ...options };
             finalOptions.response = { ...responseHeaders, ...options?.response };

            const url = client.signatureUrl(objectKey, finalOptions);
            console.log(`[OSS Download URL] Retry Success for key: ${objectKey}`);
            return url;
        } catch (retryError) {
             console.error(`[OSS Download URL] Retry Failed for key: ${objectKey}`, retryError);
             throw retryError;
        }
    }
    throw error;
  }
};

// --- Convenience function for uploading Work ---
// You'll need to define the structure of 'Work' and how files relate to it.
// This is a placeholder structure. Adjust according to your actual data model.

/**
 * Uploads a file associated with a Work item.
 * Determines the object key based on portfolio UID, project UID, and work UID/filename.
 *
 * @param file The file to upload (e.g., from uni.chooseImage).
 * @param portfolioUid The UID of the portfolio.
 * @param projectUid The UID of the project this work belongs to.
 * @param workUid The UID of the work item (or a unique identifier for the file).
 * @param fileSuffix The file extension (e.g., '.jpg', '.png').
 * @returns The generated OSS object key.
 */
export const uploadWorkFile = async (
    file: File | Blob, // Or the specific type returned by uni.chooseImage/uploadFile
    portfolioUid: string,
    projectUid: string,
    workUid: string, // Or maybe use a timestamp/random string if work UID isn't known yet
    fileSuffix: string = '.jpg' // Default or determine from file type
): Promise<string> => {
    if (!portfolioUid || !projectUid || !workUid) {
        throw new Error("Missing UIDs for generating work file object key.");
    }
    // Example object key structure: portfolios/{portfolioUid}/projects/{projectUid}/works/{workUid}{suffix}
    // Adjust this structure as needed!
    const objectKey = `portfolios/${portfolioUid}/projects/${projectUid}/works/${workUid}${fileSuffix}`;

    await uploadFile(file, objectKey);
    return objectKey; // Return the key so it can be saved in the Work data
};

// --- Helper to get file suffix ---
export const getFileSuffix = (filePathOrName: string): string => {
    const match = filePathOrName.match(/\.([^.]+)$/);
    return match ? match[0] : '';
}
