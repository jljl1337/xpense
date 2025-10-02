import { redirect } from "react-router";
import type { Route } from "./+types/delete";

import DeletePage from "~/components/pages/delete-page";
import { getCsrfToken } from "~/lib/db/auth";
import {
  deletePaymentMethod,
  getPaymentMethod,
} from "~/lib/db/payment-methods";

export async function clientLoader({ params }: Route.LoaderArgs) {
  const [csrfToken, paymentMethod] = await Promise.all([
    getCsrfToken(),
    getPaymentMethod(params.paymentMethodID),
  ]);

  if (csrfToken.error != null) {
    return redirect("/error");
  }

  if (paymentMethod.error != null) {
    return redirect("/error");
  }

  return { csrfToken: csrfToken.data, paymentMethod: paymentMethod.data };
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

  return (
    <>
      <title>Delete Payment Method | Xpense</title>
      <DeletePage
        title="Delete Payment Method"
        description={`Are you sure you want to delete the following payment method: ${loaderData.paymentMethod.name}?`}
        action={action}
        redirectTo={`/books/${loaderData.paymentMethod.bookID}/payment-methods`}
      />
    </>
  );
}
