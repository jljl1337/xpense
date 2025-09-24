import type { Route } from "./+types/account";

import { me } from "~/lib/db/user";

export async function clientLoader() {
  return await me();
}

export default function Page({ loaderData }: Route.ComponentProps) {
  return (
    <div>
      <h1 className="text-2xl font-bold mb-4">Account</h1>
      {loaderData?.data ? (
        <div>
          <p className="mb-2">User ID: {loaderData.data.id}</p>
          <p className="mb-2">Username: {loaderData.data.username}</p>
          <p className="mb-2">
            Created At: {new Date(loaderData.data.createdAt).toLocaleString()}
          </p>
        </div>
      ) : (
        <p className="text-red-500">
          Error: {loaderData?.error || "Failed to load user data."}
        </p>
      )}
    </div>
  );
}
