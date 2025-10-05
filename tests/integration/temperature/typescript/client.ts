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
    this.baseURL = config.baseURL || 'http://webservices.daehosting.com/services/TemperatureConversions.wso';
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
   * CelsiusToFahrenheit
   * Converts a Celsius Temperature to a Fahrenheit value
   */
  async celsiustofahrenheit(request: Types.CelsiustofahrenheitRequest): Promise<Types.CelsiustofahrenheitResponse> {
    return this.request<Types.CelsiustofahrenheitResponse>('/api/CelsiusToFahrenheit', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  /**
   * FahrenheitToCelsius
   * Converts a Fahrenheit Temperature to a Celsius value
   */
  async fahrenheittocelsius(request: Types.FahrenheittocelsiusRequest): Promise<Types.FahrenheittocelsiusResponse> {
    return this.request<Types.FahrenheittocelsiusResponse>('/api/FahrenheitToCelsius', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  /**
   * WindChillInCelsius
   * Windchill temperature calculated with the formula of Steadman
   */
  async windchillincelsius(request: Types.WindchillincelsiusRequest): Promise<Types.WindchillincelsiusResponse> {
    return this.request<Types.WindchillincelsiusResponse>('/api/WindChillInCelsius', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

  /**
   * WindChillInFahrenheit
   * Windchill temperature calculated with the formula of Steadman
   */
  async windchillinfahrenheit(request: Types.WindchillinfahrenheitRequest): Promise<Types.WindchillinfahrenheitResponse> {
    return this.request<Types.WindchillinfahrenheitResponse>('/api/WindChillInFahrenheit', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  }

}

// Export a default client instance
export const apiClient = new APIClient();
