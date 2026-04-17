export interface SuccessResponse<T> {
  data: T;
}

export interface ErrorResponse {
  error: {
    code: string;
    message: string;
  };
}

export interface ProbeResult {
  ok: boolean;
}
