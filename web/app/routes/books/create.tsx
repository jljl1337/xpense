import { redirect } from "react-router";
import type { Route } from "./+types/create";

import type z from "zod";

import NameDescriptionPage from "~/components/pages/name-description-page";
import { getCsrfToken } from "~/lib/db/auth";
import { createBook } from "~/lib/db/books";
import type { nameDescriptionSchema } from "~/lib/schemas/name-description";

export async function clientLoader() {
  const csrfToken = await getCsrfToken();
  if (csrfToken.error != null) {
    return redirect("/error");
  }

  return { csrfToken: csrfToken.data };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  async function action(data: z.infer<typeof nameDescriptionSchema>) {
    const response = await createBook(
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
      <title>New Book | Xpense</title>
      <NameDescriptionPage
        title={"Create new book"}
        description={"Enter the name and description for your new book"}
        nameValue={""}
        descriptionValue={""}
        nameFieldLabel={"Book Name"}
        descriptionFieldLabel={"Book Description"}
        nameFieldPlaceholder={"Trip ABC"}
        descriptionFieldPlaceholder={"Optional"}
        submitButtonLabel={"Create"}
        action={action}
      />
    </>
  );
}
