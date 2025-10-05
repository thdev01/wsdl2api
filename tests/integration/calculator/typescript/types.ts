// Auto-generated TypeScript types from OpenAPI specification

export interface AddRequest {
  parameters: string;
}

export interface AddResponse {
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

