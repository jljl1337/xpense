import { getUrl } from "~/lib/db/url";

export async function signUp(username: string, password: string) {
  const response = await fetch(getUrl("/api/auth/sign-up"), {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}

export async function login(username: string, password: string) {
  const response = await fetch(getUrl("/api/auth/login"), {
    method: "POST",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ username, password }),
  });

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}

export async function logout() {
  const response = await fetch(getUrl("/api/auth/logout"), {
    method: "POST",
  });

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}

export async function logoutAll() {
  const response = await fetch(getUrl("/api/auth/logout-all"), {
    method: "POST",
  });

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}
