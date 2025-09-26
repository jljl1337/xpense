import { Outlet } from "react-router";

import {
  SidebarInset,
  SidebarProvider,
  SidebarTrigger,
} from "~/components/ui/sidebar";

import { AppSidebar } from "~/components/app-sidebar";

export default function Layout() {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <div className="mx-2 my-2 flex items-center">
          <SidebarTrigger />
        </div>
        <main className="h-full w-full">
          <Outlet />
        </main>
      </SidebarInset>
    </SidebarProvider>
  );
}
