import { customFetch } from "~/lib/db/fetch";

type User = {
  id: string;
  username: string;
  createdAt: string;
};

export async function checkUsernameExists(username: string) {
  const response = await customFetch(
    `/api/users/exists?username=${encodeURIComponent(username)}`,
    "GET",
  );

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: { exists: boolean } = await response.json();
  return { data: data.exists, error: null };
}

export async function getMe() {
  const response = await customFetch("/api/users/me", "GET");

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: User = await response.json();
  return { data, error: null };
}

export async function deleteMe(csrfToken: string) {
  const response = await customFetch(
    "/api/users/me",
    "DELETE",
    null,
    csrfToken,
  );

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}
