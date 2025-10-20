import { redirect } from "react-router";
import type { Route } from "./+types/create";

import type z from "zod";

import NameDescriptionPage from "~/components/pages/name-description-page";
import { getCsrfToken } from "~/lib/db/auth";
import { createCategory } from "~/lib/db/categories";
import { redirectIfNeeded } from "~/lib/db/common";
import type { nameDescriptionSchema } from "~/lib/schemas/name-description";

export async function clientLoader() {
  const csrfToken = await getCsrfToken();
  redirectIfNeeded(csrfToken.error);

  return { csrfToken: csrfToken.data! };
}

export default function Page({ loaderData, params }: Route.ComponentProps) {
  async function action(data: z.infer<typeof nameDescriptionSchema>) {
    const response = await createCategory(
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
      <title>New Category | Xpense</title>
      <NameDescriptionPage
        title={"Create new category"}
        description={"Enter the name and description for your new category"}
        nameValue={""}
        descriptionValue={""}
        nameFieldLabel={"Category Name"}
        descriptionFieldLabel={"Category Description"}
        nameFieldPlaceholder={"Groceries"}
        descriptionFieldPlaceholder={"Optional"}
        submitButtonLabel={"Create"}
        action={action}
      />
    </>
  );
}
