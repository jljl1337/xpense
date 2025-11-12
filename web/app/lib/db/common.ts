import { redirect } from "react-router";

/**
 * Redirects the user based on the provided error messages.
 * If any error is "Unauthorized", redirects to the sign-in page.
 * For any other error, redirects to a generic error page.
 *
 * @param errors - An array of error messages (strings or null).
 */
export function redirectIfNeeded(...errors: (string | null)[]) {
  for (const error of errors) {
    if (error != null) {
      if (isUnauthorizedError(error)) {
        throw redirect("/auth/sign-in");
      }
      throw redirect("/error");
    }
  }
}

export function isUnauthorizedError(error: string): boolean {
  return error.trim().toLowerCase() === "unauthorized";
}
