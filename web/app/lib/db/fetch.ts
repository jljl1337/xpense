/**
 * A custom fetch function that automatically includes credentials and handles
 * JSON body and CSRF token.
 * @param path The API endpoint path
 * @param method The HTTP method (GET, POST, etc.)
 * @param body Optional request body (will be `JSON.stringify`-ed)
 * @param csrfToken Optional CSRF token to include in the headers
 * @returns The fetch response
 */
export function customFetch(
  path: string,
  method: string,
  body?: any,
  csrfToken?: string,
) {
  let options: RequestInit = {
    method,
    credentials: "include",
  };
  if (body) {
    options.body = JSON.stringify(body);
    options.headers = {
      "Content-Type": "application/json",
    };
  }
  if (csrfToken) {
    options.headers = {
      ...options.headers,
      "X-CSRF-Token": csrfToken,
    };
  }
  return fetch(getUrl(path), options);
}

function getUrl(path: string) {
  // If in development, prefix the log with http://localhost:8080
  if (import.meta.env.DEV) {
    return `http://localhost:8080${path}`;
  }
  return path;
}
