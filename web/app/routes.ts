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
    route("sign-in", "routes/auth/sign-in.tsx"),
    route("sign-up", "routes/auth/sign-up.tsx"),
  ]),

  layout("layouts/sidebar.tsx", [
    route("account", "routes/account.tsx"),
    route("about", "routes/about.tsx"),

    route("books", "routes/books/index.tsx"),
    route("books/create", "routes/books/create.tsx"),
    route("books/:bookID/edit", "routes/books/edit.tsx"),
    route("books/:bookID/delete", "routes/books/delete.tsx"),

    route("books/:bookID/categories", "routes/categories/index.tsx"),
    route("books/:bookID/categories/create", "routes/categories/create.tsx"),
    route(
      "books/:bookID/categories/:categoryID/edit",
      "routes/categories/edit.tsx",
    ),
    route(
      "books/:bookID/categories/:categoryID/delete",
      "routes/categories/delete.tsx",
    ),

    route("books/:bookID/payment-methods", "routes/payment-methods/index.tsx"),
    route(
      "books/:bookID/payment-methods/create",
      "routes/payment-methods/create.tsx",
    ),
    route(
      "books/:bookID/payment-methods/:paymentMethodID/edit",
      "routes/payment-methods/edit.tsx",
    ),
    route(
      "books/:bookID/payment-methods/:paymentMethodID/delete",
      "routes/payment-methods/delete.tsx",
    ),

    route("books/:bookID/expenses", "routes/expenses/index.tsx"),
    route("books/:bookID/expenses/create", "routes/expenses/create.tsx"),
    route("books/:bookID/expenses/:expenseID/edit", "routes/expenses/edit.tsx"),
    route(
      "books/:bookID/expenses/:expenseID/delete",
      "routes/expenses/delete.tsx",
    ),
  ]),

  route("error", "routes/error.tsx"),
  route("*", "routes/not-found.tsx"),
] satisfies RouteConfig;
