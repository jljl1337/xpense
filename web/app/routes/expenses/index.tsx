import { Link, redirect } from "react-router";
import type { Route } from "./+types";

import type { ColumnDef } from "@tanstack/react-table";

import { Button } from "~/components/ui/button";

import Pagination from "~/components/pagination";
import { DataTable } from "~/components/tables/data-table";
import TableRowDropdown from "~/components/tables/dropdown";
import { getCategories } from "~/lib/db/categories";
import { isUnauthorizedError } from "~/lib/db/common";
import {
  getExpensesByBookID,
  getExpensesCountByBookID,
  type Expense,
} from "~/lib/db/expenses";
import { getPaymentMethods } from "~/lib/db/payment-methods";

export async function clientLoader({
  params,
  request,
}: Route.ClientLoaderArgs) {
  const searchParams = new URL(request.url).searchParams;
  let page = parseInt(searchParams.get("page") ?? "1");
  if (isNaN(page) || page < 1) {
    page = 1;
  }

  const [count, categoryList, paymentMethodList, expenseList] =
    await Promise.all([
      getExpensesCountByBookID(params.bookID),
      getCategories(params.bookID),
      getPaymentMethods(params.bookID),
      getExpensesByBookID(params.bookID, page, 20),
    ]);

  if (count.error != null) {
    if (isUnauthorizedError(count.error)) {
      return redirect("/auth/sign-in");
    }
    return redirect("/error");
  }
  if (categoryList.error != null) {
    if (isUnauthorizedError(categoryList.error)) {
      return redirect("/auth/sign-in");
    }
    return redirect("/error");
  }
  if (paymentMethodList.error != null) {
    if (isUnauthorizedError(paymentMethodList.error)) {
      return redirect("/auth/sign-in");
    }
    return redirect("/error");
  }
  if (expenseList.error != null) {
    if (isUnauthorizedError(expenseList.error)) {
      return redirect("/auth/sign-in");
    }
    return redirect("/error");
  }

  return {
    count: count.data,
    page,
    expenses: expenseList.data,
    categoryList: categoryList.data,
    paymentMethodList: paymentMethodList.data,
  };
}

export default function Page({ params, loaderData }: Route.ComponentProps) {
  const categoryList = loaderData.categoryList;
  const paymentMethodList = loaderData.paymentMethodList;

  const columns: ColumnDef<Expense>[] = [
    {
      accessorKey: "category_id",
      header: "Category",
      cell: ({ row }) => {
        return (
          <Link
            to={`/books/${row.original.bookID}/categories/${row.original.categoryID}/edit`}
            className="underline"
          >
            {
              categoryList.find(
                (category) => category.id === row.original.categoryID,
              )!.name
            }
          </Link>
        );
      },
    },
    {
      accessorKey: "payment_method_id",
      header: "Payment Method",
      cell: ({ row }) => {
        return (
          <Link
            to={`/books/${row.original.bookID}/payment-methods/${row.original.paymentMethodID}/edit`}
            className="underline"
          >
            {
              paymentMethodList.find(
                (method) => method.id === row.original.paymentMethodID,
              )!.name
            }
          </Link>
        );
      },
    },
    {
      accessorKey: "date",
      header: "Date",
    },
    {
      accessorKey: "amount",
      header: "Amount",
    },
    {
      accessorKey: "remark",
      header: "Remark",
    },
    {
      accessorKey: "createdAt",
      header: "Created At",
      cell: ({ row }) => {
        return new Date(row.original.createdAt).toLocaleString();
      },
    },
    {
      accessorKey: "updatedAt",
      header: "Updated At",
      cell: ({ row }) => {
        return new Date(row.original.updatedAt).toLocaleString();
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        return (
          <div className="text-right">
            <TableRowDropdown
              editUrl={`/books/${row.original.bookID}/expenses/${row.original.id}/edit`}
              deleteUrl={`/books/${row.original.bookID}/expenses/${row.original.id}/delete`}
            />
          </div>
        );
      },
    },
  ];

  const totalPages = Math.ceil(loaderData.count / 20);
  const firstPageUrl = `/books/${params.bookID}/expenses`;
  const lastPageUrl = `/books/${params.bookID}/expenses?page=${totalPages}`;
  const previousPageUrl =
    loaderData.page > 1
      ? `/books/${params.bookID}/expenses?page=${loaderData.page - 1}`
      : "";
  const nextPageUrl =
    loaderData.page < totalPages
      ? `/books/${params.bookID}/expenses?page=${loaderData.page + 1}`
      : "";

  return (
    <>
      <title>Expenses | Xpense</title>
      <div className="h-full flex items-center justify-center">
        <div className="h-full max-w-[90rem] flex-1 flex flex-col p-8 gap-4">
          <h1 className="text-4xl">Expenses</h1>
          <Button className="w-24" asChild>
            <Link to={`/books/${params.bookID}/expenses/create`}>Create</Link>
          </Button>
          <DataTable columns={columns} data={loaderData.expenses} />
          <div className="self-end">
            {loaderData.count > 0 && (
              <Pagination
                page={loaderData.page}
                totalPages={totalPages}
                firstPageUrl={firstPageUrl}
                lastPageUrl={lastPageUrl}
                previousPageUrl={previousPageUrl}
                nextPageUrl={nextPageUrl}
              />
            )}
          </div>
        </div>
      </div>
    </>
  );
}
