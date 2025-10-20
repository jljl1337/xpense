import { redirect } from "react-router";
import type { Route } from "./+types/edit";

import type z from "zod";

import NameDescriptionPage from "~/components/pages/name-description-page";
import { getCsrfToken } from "~/lib/db/auth";
import { getBook, updateBook } from "~/lib/db/books";
import { redirectIfNeeded } from "~/lib/db/common";
import type { nameDescriptionSchema } from "~/lib/schemas/name-description";

export async function clientLoader({ params }: Route.LoaderArgs) {
  const [csrfToken, book] = await Promise.all([
    getCsrfToken(),
    getBook(params.bookID),
  ]);

  redirectIfNeeded(csrfToken.error, book.error);

  return { csrfToken: csrfToken.data!, book: book.data! };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  async function action(data: z.infer<typeof nameDescriptionSchema>) {
    const response = await updateBook(
      loaderData.book.id,
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
      <title>Edit Book | Xpense</title>
      <NameDescriptionPage
        title={"Edit Book"}
        description={"Enter the new name and description for your book"}
        nameValue={loaderData.book.name}
        descriptionValue={loaderData.book.description}
        nameFieldLabel={"Book Name"}
        descriptionFieldLabel={"Book Description"}
        nameFieldPlaceholder={"Trip ABC"}
        descriptionFieldPlaceholder={"Optional"}
        submitButtonLabel={"Update"}
        action={action}
      />
    </>
  );
}
