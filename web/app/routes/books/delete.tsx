import { redirect } from "react-router";
import type { Route } from "./+types/delete";

import DeletePage from "~/components/pages/delete-page";
import { getCsrfToken } from "~/lib/db/auth";
import { deleteBook, getBook } from "~/lib/db/books";

export async function clientLoader({ params }: Route.LoaderArgs) {
  const [csrfToken, book] = await Promise.all([
    getCsrfToken(),
    getBook(params.bookID),
  ]);

  if (csrfToken.error != null) {
    return redirect("/error");
  }

  if (book.error != null) {
    return redirect("/error");
  }

  return { csrfToken: csrfToken.data, book: book.data };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  async function action() {
    const response = await deleteBook(loaderData.book.id, loaderData.csrfToken);
    if (response.error != null) {
      return { error: response.error };
    }

    return { error: null };
  }

  return (
    <>
      <title>Delete Book | Xpense</title>
      <DeletePage
        title="Delete Book"
        description={`Are you sure you want to delete the following book: ${loaderData.book.name}?`}
        action={action}
      />
    </>
  );
}
