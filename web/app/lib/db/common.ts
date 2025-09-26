export function isUnauthorizedError(error: string): boolean {
  return error.trim() === "Unauthorized";
}
