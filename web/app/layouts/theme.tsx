import { Outlet } from "react-router";

import { ThemeProvider } from "~/components/theme-provider";

export default function ThemeLayout() {
  return (
    <ThemeProvider>
      <Outlet />
    </ThemeProvider>
  );
}
