import { customFetch } from "~/lib/db/url";

type CsrfToken = {
  csrfToken: string;
};

export async function signUp(username: string, password: string) {
  const response = await customFetch("/api/auth/sign-up", "POST", {
    username,
    password,
  });

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}

export async function login(username: string, password: string) {
  const response = await customFetch("/api/auth/login", "POST", {
    username,
    password,
  });

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}

export async function getCsrfToken() {
  const response = await customFetch("/api/auth/csrf-token", "GET");

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: CsrfToken = await response.json();
  return { data: data.csrfToken, error: null };
}

export async function logout(csrfToken: string) {
  const response = await customFetch(
    "/api/auth/logout",
    "POST",
    null,
    csrfToken,
  );

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}

export async function logoutAll(csrfToken: string) {
  const response = await customFetch(
    "/api/auth/logout-all",
    "POST",
    null,
    csrfToken,
  );

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}
