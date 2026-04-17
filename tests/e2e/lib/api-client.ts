import { APIRequestContext, request as playwrightRequest } from "@playwright/test";

export class ApiClient {
  constructor(
    private request: APIRequestContext,
    private token: string = ""
  ) {}

  private headers() {
    const h: Record<string, string> = {};
    if (this.token) {
      h["Authorization"] = `Bearer ${this.token}`;
    }
    return h;
  }

  async get(path: string, params?: Record<string, string>) {
    return this.request.get(path, {
      headers: this.headers(),
      params,
    });
  }

  async post(path: string, data?: unknown) {
    return this.request.post(path, {
      headers: this.headers(),
      data,
    });
  }

  async put(path: string, data?: unknown) {
    return this.request.put(path, {
      headers: this.headers(),
      data,
    });
  }

  async delete(path: string) {
    return this.request.delete(path, {
      headers: this.headers(),
    });
  }
}

export async function createApiClient(
  baseURL: string,
  token: string = ""
): Promise<ApiClient> {
  const ctx = await playwrightRequest.newContext({ baseURL });
  return new ApiClient(ctx, token);
}
