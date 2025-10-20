import { redirect } from "react-router";
import type { Route } from "./+types/delete";

import DeletePage from "~/components/pages/delete-page";
import { getCsrfToken } from "~/lib/db/auth";
import { deleteCategory, getCategory } from "~/lib/db/categories";
import { redirectIfNeeded } from "~/lib/db/common";

export async function clientLoader({ params }: Route.LoaderArgs) {
  const [csrfToken, category] = await Promise.all([
    getCsrfToken(),
    getCategory(params.categoryID),
  ]);

  redirectIfNeeded(csrfToken.error, category.error);

  return { csrfToken: csrfToken.data!, category: category.data! };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  async function action() {
    const response = await deleteCategory(
      loaderData.category.id,
      loaderData.csrfToken,
    );
    if (response.error != null) {
      return { error: response.error };
    }

    return { error: null };
  }

  return (
    <>
      <title>Delete Category | Xpense</title>
      <DeletePage
        title="Delete Category"
        description={`Are you sure you want to delete the following category: ${loaderData.category.name}?`}
        action={action}
        redirectTo={`/books/${loaderData.category.bookID}/categories`}
      />
    </>
  );
}
