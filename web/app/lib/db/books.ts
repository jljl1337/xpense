import { customFetch } from "~/lib/db/fetch";

export async function createBook(
  name: string,
  description: string,
  csrfToken: string,
) {
  const response = await customFetch(
    "/api/books",
    "POST",
    {
      name,
      description,
    },
    csrfToken,
  );

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}
