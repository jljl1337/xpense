import { redirect } from "react-router";
import type { Route } from "./+types/about";

import { isUnauthorizedError } from "~/lib/db/common";
import { getVersion } from "~/lib/db/version";

export async function clientLoader() {
  const version = await getVersion();

  if (version.error != null) {
    if (isUnauthorizedError(version.error)) {
      return redirect("/auth/sign-in");
    }
    return redirect("/error");
  }
  return { data: { version: version.data }, error: null };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  return (
    <>
      <title>About | Xpense</title>
      <div className="h-full flex items-center justify-center">
        <div className="h-full max-w-[90rem] flex-1 flex flex-col p-8 gap-4">
          <h1 className="text-4xl">About</h1>
          <p className="mb-2">Version: {loaderData.data.version}</p>
        </div>
      </div>
    </>
  );
}
