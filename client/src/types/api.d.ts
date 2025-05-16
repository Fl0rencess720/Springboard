// src/types/api.d.ts

/**
 * 通用后端响应结构
 */
export interface ApiResponse<T = any> {
  code: number; // Business status code, 200 means success
  msg: string; // Response message
  data: T; // Response data payload
}

/**
 * Represents the status of feedback.
 * Pending = 0, Approved = 1, Rejected = 2
 */
export type FeedbackStatus = 0 | 1 | 2;

/**
 * Feedback data structure.
 * Matches Go: data.Feedback
 */
export interface Feedback {
  uid: string;
  content: string;
  timestamp: string; // Assuming time.Time serializes to string
  status: FeedbackStatus;
}

/**
 * Work data structure (e.g., an image within a project).
 * Matches Go: data.Work
 */
export interface Work {
  // ID?: number; // Backend primary key, usually not needed in frontend models unless for specific ops
  oss_key: string;
  project_uid: string;
  size: string; // Format: "axb", e.g., "1920x1080"
  margin_top: string;
  margin_left: string;
  page: number;
  CreatedAt?: string; // Optional: Added from Go model
  UpdatedAt?: string; // Optional: Added from Go model
}

/**
 * Text data structure.
 * Matches Go: data.Text
 */
export interface Text {
    // ID?: number; // Backend primary key
    content: string;
    font_size: string;
    size: string; // Format: "axb"
    margin_top: string;
    margin_left: string;
    page: number;
    portfolio_uid: string;
}

/**
 * Page data structure (part of a Template).
 * Matches Go: data.Page
 */
export interface Page {
    // ID?: number; // Backend primary key
    uid: string;
    oss_key: string; // Background/layout image for the page template
    preview_oss_key: string; // Preview image for the page template
    bleed: string[]; // Array of 4 strings: svg x, y, width, height
    template_uid: string;
    margin_top: string;
    margin_left: string;
    size: string; // Format: "axb"
    is_content_page: boolean;
}

/**
 * Project data structure (a group of works within a portfolio).
 * Matches Go: data.Project
 */
export interface Project {
  // ID?: number; // Backend primary key
  uid: string;
  name: string;
  order: number;
  portfolio_uid: string;
  works: Work[];
  CreatedAt?: string; // Optional: Added from Go model
  UpdatedAt?: string; // Optional: Added from Go model
}

/**
 * Template data structure.
 * Matches Go: data.Template
 */
export interface Template {
  // ID?: number; // Backend primary key
  uid: string;
  name: string;
  font_oss_key: string; // Changed from oss_key
  pages: Page[]; // Added from Go model
  CreatedAt?: string; // Optional: Added from Go model
}

/**
 * Portfolio data structure.
 * Matches Go: data.Portfolio
 */
export interface Portfolio {
  // ID?: number; // Backend primary key
  uid: string;
  openid: string;
  title: string;
  projects: Project[];
  texts: Text[]; // Added from Go model
  template_uid: string;
  template: Template; // Updated Template type
  CreatedAt?: string; // Optional: Added from Go model
  UpdatedAt?: string; // Optional: Added from Go model
}

// --- API Request/Response Specific Types ---

/**
 * Generic type for API responses wrapping the actual data type.
 */
export type ApiResult<T> = ApiResponse<T>;

/**
 * Request body for saving a portfolio.
 * Matches Go: controller.SavePortfolioRequest
 */
export interface SavePortfolioRequest {
  uid: string; // Portfolio UID. Empty string for first save.
  title: string;
  template_uid: string;
  // Use Omit to exclude fields potentially added by backend/frontend state
  // Also omit fields from nested Work objects
  projects: Omit<Project, 'ID' | 'CreatedAt' | 'UpdatedAt' | 'portfolio_uid' | 'works'> & {
      works: Omit<Work, 'ID' | 'CreatedAt' | 'UpdatedAt' | 'project_uid'>[];
  }[];
  // Note: Backend model shows Portfolio has Texts, but SavePortfolioRequest doesn't.
  // Assuming texts are saved separately or implicitly handled.
}

/**
 * Response data structure for saving a portfolio.
 */
export interface SavePortfolioResponseData {
    projects: Project[]; // Includes backend-generated fields like ID, UID, CreatedAt, etc.
    uid: string; // The portfolio UID (new or existing)
}

export type SavePortfolioResponse = ApiResult<SavePortfolioResponseData>;

/**
 * Response structure for getting all templates.
 * Note: Assumes API returns a simplified structure with 'oss_key' for preview.
 * The main Template interface uses 'font_oss_key' and 'pages'.
 */
export interface GetAllTemplatesResponseItem {
  uid: string;
  name: string;
  oss_key: string;
  // Add CreatedAt if it exists in the actual response based on be-interface.md
  CreatedAt?: string;
}

// Define the full response type using ApiResult
export type GetAllTemplatesResponse = ApiResult<GetAllTemplatesResponseItem[]>;

/**
 * Response structure for getting hot templates.
  */
export interface GetHotTemplatesResponseItem {
  uid: string;
  name: string;
  oss_key: string;
  CreatedAt: string;
}
export type GetHotTemplatesResponse = ApiResult<GetHotTemplatesResponseItem[]>;

/**
 * Response structure for getting my portfolios.
 */
export interface GetMyPortfolioResponseItem extends Portfolio {
    // Portfolio interface already includes all fields from the Go model including nested ones and CreatedAt/UpdatedAt
}
export type GetMyPortfolioResponse = ApiResult<GetMyPortfolioResponseItem[]>;

/**
 * Response structure for getting historical usage templates.
  */
export interface GetHistoryTemplatesResponseItem {
    uid: string;
    name: string;
    oss_key: string;
    CreatedAt: string;
}
export type GetHistoryTemplatesResponse = ApiResult<GetHistoryTemplatesResponseItem[]>;

/**
 * Response structure for getting feedbacks.
 */
export type GetAllFeedbacksResponse = ApiResult<Feedback[]>;
export type GetFeedbacksByStatusResponse = ApiResult<Feedback[]>;

/**
 * Request body for adding feedback.
 */
export interface AddFeedbackRequest {
    content: string;
}
export type AddFeedbackResponse = ApiResult<null>; // Assuming data is null on success

/**
 * Request body for updating feedback status.
 */
export interface UpdateFeedbackStatusRequest {
    uid: string;
    status: FeedbackStatus;
}
export type UpdateFeedbackResponse = ApiResult<null>; // Assuming data is null on success

/**
 * Response structure for getting OSS STS credentials.
 */
export interface GetCredentialsResponseData {
    AccessKeyId: string;
    AccessKeySecret: string;
    SecurityToken: string;
}
export type GetCredentialsResponse = ApiResult<GetCredentialsResponseData>;

/**
 * 登录接口响应数据结构 (/login)
 */
export interface LoginResponseData {
  access_token: string;
  refresh_token: string; // Added refresh token
}
export type LoginResponse = ApiResult<LoginResponseData>;

/**
 * 刷新 Token 接口响应数据结构 (/refresh)
 */
export interface RefreshTokenResponseData {
  access_token: string;
}
export type RefreshTokenResponse = ApiResult<RefreshTokenResponseData>;
