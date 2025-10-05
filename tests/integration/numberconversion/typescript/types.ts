// Auto-generated TypeScript types from OpenAPI specification

export interface NumbertowordsRequest {
  parameters: string;
}

export interface NumbertowordsResponse {
  parameters: string;
}

export interface NumbertodollarsRequest {
  parameters: string;
}

export interface NumbertodollarsResponse {
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

