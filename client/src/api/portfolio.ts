
import { request } from './request';
import type { SavePortfolioRequest, SavePortfolioResponse, GetMyPortfolioResponse, GetHistoryTemplatesResponse } from '../types/api';

/**
 * Saves a portfolio.
 * POST /api/portfolio/portfolio/save
 * @param payload - The portfolio data to save.
 */
export const savePortfolio = (payload: SavePortfolioRequest) => {
  return request<SavePortfolioResponse['data']>({ // Specify the expected data type
    url: '/portfolio/portfolio/save',
    method: 'POST',
    data: payload,
  });
};

/**
 * Fetches the current user's portfolios.
 * GET /api/portfolio/portfolio/me
 */
export const fetchMyPortfolios = () => {
  return request<GetMyPortfolioResponse['data']>({ // Specify the expected data type
    url: '/portfolio/portfolio/me',
    method: 'GET',
  });
};

/**
 * Fetches the current user's historically used templates.
 * GET /api/portfolio/portfolio/history
 */
export const fetchHistoryTemplates = () => {
  return request<GetHistoryTemplatesResponse['data']>({ // Specify the expected data type
    url: '/portfolio/portfolio/history',
    method: 'GET',
  });
};
