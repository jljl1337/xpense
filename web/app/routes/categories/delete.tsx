import { redirect } from "react-router";
import type { Route } from "./+types/delete";

import DeletePage from "~/components/pages/delete-page";
import { getCsrfToken } from "~/lib/db/auth";
import { deleteCategory, getCategory } from "~/lib/db/categories";
import { redirectIfNeeded } from "~/lib/db/common";
import { getExpensesCountByBookID } from "~/lib/db/expenses";

export async function clientLoader({ params }: Route.LoaderArgs) {
  const [csrfToken, category] = await Promise.all([
    getCsrfToken(),
    getCategory(params.categoryID),
  ]);

  redirectIfNeeded(csrfToken.error, category.error);

  const expenseCount = await getExpensesCountByBookID(
    category.data!.bookID,
    category.data!.id,
  );

  redirectIfNeeded(expenseCount.error);

  return {
    csrfToken: csrfToken.data!,
    category: category.data!,
    expenseCount: expenseCount.data!,
  };
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

  let description = `Are you sure you want to delete the following category: ${loaderData.category.name}?`;

  if (loaderData.expenseCount > 0) {
    description += ` The associated ${loaderData.expenseCount} expense(s) will also be deleted. This action cannot be undone.`;
  }

  return (
    <>
      <title>Delete Category | Xpense</title>
      <DeletePage
        title="Delete Category"
        description={description}
        action={action}
        redirectTo={`/books/${loaderData.category.bookID}/categories`}
      />
    </>
  );
}
