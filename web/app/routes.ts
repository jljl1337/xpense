import {
  index,
  layout,
  prefix,
  route,
  type RouteConfig,
} from "@react-router/dev/routes";

export default [
  layout("layouts/theme.tsx", [
    index("routes/home.tsx"),

    ...prefix("auth", [
      route("login", "routes/auth/login.tsx"),
      route("sign-up", "routes/auth/sign-up.tsx"),
    ]),

    route("account", "routes/account.tsx"),
  ]),
] satisfies RouteConfig;
