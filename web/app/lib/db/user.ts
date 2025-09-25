import { customFetch } from "~/lib/db/url";

type User = {
  id: string;
  username: string;
  createdAt: number;
};

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
