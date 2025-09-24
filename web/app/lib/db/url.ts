export function getUrl(path: string) {
  // If in development, prefix the log with http://localhost:8080
  if (import.meta.env.DEV) {
    return `http://localhost:8080${path}`;
  }
  return path;
}
