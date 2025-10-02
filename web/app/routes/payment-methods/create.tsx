import { redirect } from "react-router";
import type { Route } from "./+types/create";

import type z from "zod";

import NameDescriptionPage from "~/components/pages/name-description-page";
import { getCsrfToken } from "~/lib/db/auth";
import { createPaymentMethod } from "~/lib/db/payment-methods";
import type { nameDescriptionSchema } from "~/lib/schemas/name-description";

export async function clientLoader() {
  const csrfToken = await getCsrfToken();
  if (csrfToken.error != null) {
    return redirect("/error");
  }

  return { csrfToken: csrfToken.data };
}

export default function Page({ loaderData, params }: Route.ComponentProps) {
  async function action(data: z.infer<typeof nameDescriptionSchema>) {
    const response = await createPaymentMethod(
      params.bookID,
      data.name,
      data.description,
      loaderData.csrfToken,
    );

    if (response.error != null) {
      return { error: response.error };
    }

    return { error: null };
  }

  return (
    <>
      <title>New Payment Method | Xpense</title>
      <NameDescriptionPage
        title={"Create new payment method"}
        description={
          "Enter the name and description for your new payment method"
        }
        nameValue={""}
        descriptionValue={""}
        nameFieldLabel={"Payment Method Name"}
        descriptionFieldLabel={"Payment Method Description"}
        nameFieldPlaceholder={"Credit Card"}
        descriptionFieldPlaceholder={"Optional"}
        submitButtonLabel={"Create"}
        action={action}
      />
    </>
  );
}
