import { customFetch } from "~/lib/db/fetch";

export type Category = {
  id: string;
  bookID: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
};

export async function createCategory(
  bookID: string,
  name: string,
  description: string,
  csrfToken: string,
) {
  const response = await customFetch(
    "/api/categories",
    "POST",
    {
      bookID,
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

export async function getCategories(bookID: string) {
  const response = await customFetch(
    `/api/categories?book-id=${bookID}`,
    "GET",
  );

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: Category[] = await response.json();
  return { data, error: null };
}

export async function getCategory(categoryID: string) {
  const response = await customFetch(`/api/categories/${categoryID}`, "GET");

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: Category = await response.json();
  return { data, error: null };
}

export async function updateCategory(
  categoryID: string,
  name: string,
  description: string,
  csrfToken: string,
) {
  const response = await customFetch(
    `/api/categories/${categoryID}`,
    "PUT",
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

export async function deleteCategory(categoryID: string, csrfToken: string) {
  const response = await customFetch(
    `/api/categories/${categoryID}`,
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
