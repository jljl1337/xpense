import { redirect } from "react-router";

import { getCsrfToken } from "~/lib/db/auth";
import { isUnauthorizedError } from "~/lib/db/common";

export async function clientLoader() {
  const csrfToken = await getCsrfToken();

  if (csrfToken.error != null && isUnauthorizedError(csrfToken.error)) {
    return redirect("/auth/sign-in");
  }
  return redirect("/books");
}
