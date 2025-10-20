import { redirect } from "react-router";
import type { Route } from "./+types/edit";

import type z from "zod";

import NameDescriptionPage from "~/components/pages/name-description-page";
import { getCsrfToken } from "~/lib/db/auth";
import { redirectIfNeeded } from "~/lib/db/common";
import {
  getPaymentMethod,
  updatePaymentMethod,
} from "~/lib/db/payment-methods";
import type { nameDescriptionSchema } from "~/lib/schemas/name-description";

export async function clientLoader({ params }: Route.LoaderArgs) {
  const [csrfToken, paymentMethod] = await Promise.all([
    getCsrfToken(),
    getPaymentMethod(params.paymentMethodID),
  ]);

  redirectIfNeeded(csrfToken.error, paymentMethod.error);

  return { csrfToken: csrfToken.data!, paymentMethod: paymentMethod.data! };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  async function action(data: z.infer<typeof nameDescriptionSchema>) {
    const response = await updatePaymentMethod(
      loaderData.paymentMethod.id,
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
      <title>Edit Payment Method | Xpense</title>
      <NameDescriptionPage
        title={"Edit Payment Method"}
        description={
          "Enter the new name and description for your payment method"
        }
        nameValue={loaderData.paymentMethod.name}
        descriptionValue={loaderData.paymentMethod.description}
        nameFieldLabel={"Payment Method Name"}
        descriptionFieldLabel={"Payment Method Description"}
        nameFieldPlaceholder={"Credit Card"}
        descriptionFieldPlaceholder={"Optional"}
        submitButtonLabel={"Update"}
        action={action}
      />
    </>
  );
}
