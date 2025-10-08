import { redirect } from "react-router";

import { isUnauthorizedError } from "~/lib/db/common";
import { getMe } from "~/lib/db/users";

export async function clientLoader() {
  const me = await getMe();

  if (me.error != null && isUnauthorizedError(me.error)) {
    return redirect("/auth/sign-in");
  }
  return redirect("/books");
}
