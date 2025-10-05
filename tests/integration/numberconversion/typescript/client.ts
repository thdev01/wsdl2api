// Auto-generated API client from OpenAPI specification

import type * as Types from './types';

export interface ClientConfig {
  baseURL?: string;
  headers?: Record<string, string>;
  timeout?: number;
}

export class APIClient {
  private baseURL: string;
  private headers: Record<string, string>;
  private timeout: number;

  constructor(config: ClientConfig = {}) {
    this.baseURL = config.baseURL || 'https://www.dataaccess.com/webservicesserver/numberconversion.wso';
    this.headers = config.headers || {};
    this.timeout = config.timeout || 30000;
  }

  private async request<T>(
    path: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = this.baseURL + path;
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...this.headers,
          ...options.headers,
        },
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        const error: Types.APIError = {
          message: response.statusText,
          status: response.status,
        };

        try {
          const fault = await response.json();
          error.fault = fault;
        } catch {
          // No JSON body
        }

        throw error;
      }

      return await response.json();
    } catch (err) {
      clearTimeout(timeoutId);

      if (err instanceof Error && err.name === 'AbortError') {
        throw new Error('Request timeout');
      }

      throw err;
    }
  }

  /**
   * NumberToWords
   * Returns the word corresponding to the positive number passed as parameter. Limited to quadrillions.
   */
  async numbertowords(request: Types.NumbertowordsRequest): Promise<Types.NumbertowordsResponse> {
    return this.request<Types.NumbertowordsResponse>('/api/NumberToWords', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  /**
   * NumberToDollars
   * Returns the non-zero dollar amount of the passed number.
   */
  async numbertodollars(request: Types.NumbertodollarsRequest): Promise<Types.NumbertodollarsResponse> {
    return this.request<Types.NumbertodollarsResponse>('/api/NumberToDollars', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

}

// Export a default client instance
export const apiClient = new APIClient();
