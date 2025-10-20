import { redirect } from "react-router";
import type { Route } from "./+types/delete";

import DeletePage from "~/components/pages/delete-page";
import { getCsrfToken } from "~/lib/db/auth";
import { redirectIfNeeded } from "~/lib/db/common";
import { deleteExpense, getExpense } from "~/lib/db/expenses";

export async function clientLoader({ params }: Route.LoaderArgs) {
  const [csrfToken, expense] = await Promise.all([
    getCsrfToken(),
    getExpense(params.expenseID),
  ]);

  redirectIfNeeded(csrfToken.error, expense.error);

  return { csrfToken: csrfToken.data!, expense: expense.data! };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  async function action() {
    const response = await deleteExpense(
      loaderData.expense.id,
      loaderData.csrfToken,
    );
    if (response.error != null) {
      return { error: response.error };
    }

    return { error: null };
  }

  return (
    <>
      <title>Delete Expense | Xpense</title>
      <DeletePage
        title="Delete Expense"
        description={`Are you sure you want to delete this expense of amount: ${loaderData.expense.amount}?`}
        action={action}
        redirectTo={`/books/${loaderData.expense.bookID}/expenses`}
      />
    </>
  );
}
