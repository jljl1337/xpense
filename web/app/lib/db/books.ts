import { customFetch } from "~/lib/db/fetch";

type BookCount = {
  count: number;
};

export type Book = {
  id: string;
  userID: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
};

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

export async function getBooksCount() {
  const response = await customFetch("/api/books/count", "GET");

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: BookCount = await response.json();
  return { data: data.count, error: null };
}

export async function getBooks(page: number, pageSize: number) {
  const response = await customFetch(
    `/api/books?page=${page}&page-size=${pageSize}`,
    "GET",
  );

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: Book[] = await response.json();
  return { data, error: null };
}
