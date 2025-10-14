import { customFetch } from "~/lib/db/fetch";

export async function getVersion() {
  const response = await customFetch("/api/version", "GET");
  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }
  const data: { version: string } = await response.json();
  return { data: data.version, error: null };
}
