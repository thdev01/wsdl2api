// Auto-generated TypeScript types from OpenAPI specification

export interface WindchillincelsiusRequest {
  parameters: string;
}

export interface WindchillincelsiusResponse {
  parameters: string;
}

export interface WindchillinfahrenheitRequest {
  parameters: string;
}

export interface WindchillinfahrenheitResponse {
  parameters: string;
}

export interface CelsiustofahrenheitRequest {
  parameters: string;
}

export interface CelsiustofahrenheitResponse {
  parameters: string;
}

export interface FahrenheittocelsiusRequest {
  parameters: string;
}

export interface FahrenheittocelsiusResponse {
  parameters: string;
}

// Error types
export interface SOAPFault {
  faultcode: string;
  faultstring: string;
  detail?: string;
}

export interface APIError {
  message: string;
  status: number;
  fault?: SOAPFault;
}

