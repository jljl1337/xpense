import { getUrl } from "~/lib/db/url";

type User = {
  id: string;
  username: string;
  createdAt: number;
};

export async function me() {
  const response = await fetch(getUrl("/api/users/me"), {
    method: "GET",
    credentials: "include",
  });

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: User = await response.json();
  return { data, error: null };
}
