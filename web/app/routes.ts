import {
  index,
  layout,
  prefix,
  route,
  type RouteConfig,
} from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),

  ...prefix("auth", [
    route("login", "routes/auth/login.tsx"),
    route("sign-up", "routes/auth/sign-up.tsx"),
  ]),

  layout("layouts/sidebar.tsx", [
    route("account", "routes/account.tsx"),

    route("books", "routes/books/index.tsx"),
  ]),

  route("*", "routes/not-found.tsx"),
] satisfies RouteConfig;
