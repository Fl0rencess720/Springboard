import { request } from './request';
import type { 
    GetAllTemplatesResponse,
    GetHotTemplatesResponse,
    Template as TemplateDetails, // Rename imported Template to avoid conflict 
    ApiResult 
} from '../types/api';

/**
 * Fetches all available templates.
 * GET /api/portfolio/template/all
 */
export const fetchAllTemplates = () => {
    return request<GetAllTemplatesResponse['data']>({ // Specify the expected data type
    url: '/portfolio/template/all',
    method: 'GET',
  });
};

/**
 * Fetches hot templates.
 * GET /api/portfolio/template/hot
 */
export const fetchHotTemplates = () => {
  // Assuming the endpoint returns data matching GetHotTemplatesResponseItem[]
  // This might already be defined elsewhere, ensure consistency
  return request<GetHotTemplatesResponse['data']>({
    url: '/portfolio/template/hot', // Replace with actual endpoint
    method: 'GET',
  });
};

/**
 * Fetches detailed information for a specific template.
 * GET /api/template?uid={uid}
 * @param uid - The UID of the template to fetch.
 */
export const fetchTemplateDetails = (uid: string) => {
  // Assuming the endpoint returns the full Template structure
  return request<TemplateDetails>({ // Use the renamed TemplateDetails type
    url: `/portfolio/template?uid=${uid}`, // Replace with actual endpoint
    method: 'GET',
  });
};
