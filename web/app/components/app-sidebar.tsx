import { Link, useLocation } from "react-router";

import {
  Book,
  DollarSign,
  Edit,
  Info,
  Laptop,
  List,
  Moon,
  Sun,
  Trash,
  User,
  Wallet,
} from "lucide-react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "~/components/ui/sidebar";

import { useTheme } from "~/components/theme-provider";

export function AppSidebar() {
  const { setTheme, theme } = useTheme();

  const { pathname } = useLocation();

  const bookIDRegex = /(?<=\/books\/)[0-9A-HJKMNP-TV-Z]{26}/;
  const bookIDMatch = pathname.match(bookIDRegex);

  const bookID = bookIDMatch != null ? bookIDMatch[0] : null;

  return (
    <Sidebar variant="inset">
      {bookID == null ? (
        <>
          <SidebarContent>
            <SidebarGroup>
              <SidebarGroupLabel>Contents</SidebarGroupLabel>
              <SidebarMenu>
                <SidebarMenuItem key="books">
                  <SidebarMenuButton asChild isActive={pathname === "/books"}>
                    <Link to={"/books"}>
                      <Book />
                      Books
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroup>
          </SidebarContent>
        </>
      ) : (
        <SidebarContent>
          <SidebarGroup>
            <SidebarGroupLabel>Contents</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem key="expense">
                  <SidebarMenuButton
                    asChild
                    isActive={pathname.startsWith(`/books/${bookID}/expenses`)}
                  >
                    <Link to={`/books/${bookID}/expenses`}>
                      <DollarSign />
                      Expenses
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem key="category">
                  <SidebarMenuButton
                    asChild
                    isActive={pathname.startsWith(
                      `/books/${bookID}/categories`,
                    )}
                  >
                    <Link to={`/books/${bookID}/categories`}>
                      <List />
                      Categories
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem key="payment-method">
                  <SidebarMenuButton
                    asChild
                    isActive={pathname.startsWith(
                      `/books/${bookID}/payment-methods`,
                    )}
                  >
                    <Link to={`/books/${bookID}/payment-methods`}>
                      <Wallet />
                      Payment Methods
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
          <SidebarGroup>
            <SidebarGroupLabel>Books</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                <SidebarMenuItem key="edit">
                  <SidebarMenuButton
                    asChild
                    isActive={pathname === `/books/${bookID}/edit`}
                  >
                    <Link to={`/books/${bookID}/edit`}>
                      <Edit />
                      Edit
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem key="delete">
                  <SidebarMenuButton
                    asChild
                    isActive={pathname === `/books/${bookID}/delete`}
                  >
                    <Link to={`/books/${bookID}/delete`}>
                      <Trash />
                      Delete
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
                <SidebarMenuItem key="books">
                  <SidebarMenuButton asChild>
                    <Link to="/books">
                      <Book />
                      All Books
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        </SidebarContent>
      )}
      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton className="cursor-pointer">
                  {theme === "light" ? (
                    <Sun />
                  ) : theme === "dark" ? (
                    <Moon />
                  ) : (
                    <Laptop />
                  )}
                  Theme
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <DropdownMenuRadioGroup value={theme}>
                  <DropdownMenuRadioItem
                    className="cursor-pointer"
                    value="system"
                    onSelect={() => setTheme("system")}
                  >
                    System
                  </DropdownMenuRadioItem>
                  <DropdownMenuRadioItem
                    className="cursor-pointer"
                    value="light"
                    onSelect={() => setTheme("light")}
                  >
                    Light
                  </DropdownMenuRadioItem>
                  <DropdownMenuRadioItem
                    className="cursor-pointer"
                    value="dark"
                    onSelect={() => setTheme("dark")}
                  >
                    Dark
                  </DropdownMenuRadioItem>
                </DropdownMenuRadioGroup>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
          <SidebarMenuItem key="account">
            <SidebarMenuButton asChild isActive={pathname === "/account"}>
              <Link to={"/account"}>
                <User />
                Account
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem key="about">
            <SidebarMenuButton asChild isActive={pathname === "/about"}>
              <Link to={"/about"}>
                <Info />
                About
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
    </Sidebar>
  );
}
