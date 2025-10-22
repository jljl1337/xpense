import { redirect } from "react-router";
import type { Route } from "./+types/delete";

import DeletePage from "~/components/pages/delete-page";
import { getCsrfToken } from "~/lib/db/auth";
import { redirectIfNeeded } from "~/lib/db/common";
import { getExpensesCountByBookID } from "~/lib/db/expenses";
import {
  deletePaymentMethod,
  getPaymentMethod,
} from "~/lib/db/payment-methods";

export async function clientLoader({ params }: Route.LoaderArgs) {
  const [csrfToken, paymentMethod] = await Promise.all([
    getCsrfToken(),
    getPaymentMethod(params.paymentMethodID),
  ]);

  redirectIfNeeded(csrfToken.error, paymentMethod.error);

  const expenseCount = await getExpensesCountByBookID(
    paymentMethod.data!.bookID,
    undefined,
    paymentMethod.data!.id,
  );

  redirectIfNeeded(expenseCount.error);

  return {
    csrfToken: csrfToken.data!,
    paymentMethod: paymentMethod.data!,
    expenseCount: expenseCount.data!,
  };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  async function action() {
    const response = await deletePaymentMethod(
      loaderData.paymentMethod.id,
      loaderData.csrfToken,
    );
    if (response.error != null) {
      return { error: response.error };
    }

    return { error: null };
  }

  let description = `Are you sure you want to delete the following payment method: ${loaderData.paymentMethod.name}?`;

  if (loaderData.expenseCount > 0) {
    description += ` The associated ${loaderData.expenseCount} expense(s) will also be deleted. This action cannot be undone.`;
  }

  return (
    <>
      <title>Delete Payment Method | Xpense</title>
      <DeletePage
        title="Delete Payment Method"
        description={description}
        action={action}
        redirectTo={`/books/${loaderData.paymentMethod.bookID}/payment-methods`}
      />
    </>
  );
}
