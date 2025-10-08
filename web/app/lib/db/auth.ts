import { customFetch } from "~/lib/db/fetch";

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

export async function getPreSession() {
  const response = await customFetch("/api/auth/pre-session", "POST");

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: CsrfToken = await response.json();
  return { data: data.csrfToken, error: null };
}

export async function signIn(
  username: string,
  password: string,
  csrfToken: string,
) {
  const response = await customFetch(
    "/api/auth/sign-in",
    "POST",
    {
      username,
      password,
    },
    csrfToken,
  );

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

export async function signOut(csrfToken: string) {
  const response = await customFetch(
    "/api/auth/sign-out",
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

export async function signOutAll(csrfToken: string) {
  const response = await customFetch(
    "/api/auth/sign-out-all",
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
