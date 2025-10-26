import { Link, redirect } from "react-router";
import type { Route } from "./+types";

import type { ColumnDef } from "@tanstack/react-table";

import { Button } from "~/components/ui/button";

import ExpenseFilterForm from "~/components/expense-filter-form";
import Pagination from "~/components/pagination";
import { DataTable } from "~/components/tables/data-table";
import TableRowDropdown from "~/components/tables/dropdown";
import { getCategories } from "~/lib/db/categories";
import { redirectIfNeeded } from "~/lib/db/common";
import {
  getExpensesByBookID,
  getExpensesCountByBookID,
  type Expense,
} from "~/lib/db/expenses";
import { getPaymentMethods } from "~/lib/db/payment-methods";
import { YYYYMMDDToLocaleDateString } from "~/lib/format/date";

export async function clientLoader({
  params,
  request,
}: Route.ClientLoaderArgs) {
  const searchParams = new URL(request.url).searchParams;
  let page = parseInt(searchParams.get("page") ?? "1");
  if (isNaN(page) || page < 1) {
    page = 1;
  }

  const categoryID = searchParams.get("category-id") ?? undefined;
  const paymentMethodID = searchParams.get("payment-method-id") ?? undefined;
  const remark = searchParams.get("remark") ?? undefined;

  const [count, categoryList, paymentMethodList, expenseList] =
    await Promise.all([
      getExpensesCountByBookID(
        params.bookID,
        categoryID,
        paymentMethodID,
        remark,
      ),
      getCategories(params.bookID),
      getPaymentMethods(params.bookID),
      getExpensesByBookID(
        params.bookID,
        page,
        20,
        categoryID,
        paymentMethodID,
        remark,
      ),
    ]);

  redirectIfNeeded(
    count.error,
    categoryList.error,
    paymentMethodList.error,
    expenseList.error,
  );

  return {
    count: count.data!,
    page,
    expenses: expenseList.data!,
    categoryList: categoryList.data!,
    paymentMethodList: paymentMethodList.data!,
    categoryID,
    paymentMethodID,
    remark,
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
      cell: ({ row }) => {
        return YYYYMMDDToLocaleDateString(row.original.date);
      },
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

  const filterParams = new URLSearchParams();
  if (loaderData.categoryID) {
    filterParams.append("category-id", loaderData.categoryID);
  }
  if (loaderData.paymentMethodID) {
    filterParams.append("payment-method-id", loaderData.paymentMethodID);
  }
  if (loaderData.remark) {
    filterParams.append("remark", loaderData.remark);
  }

  const totalPages = Math.ceil(loaderData.count / 20);

  const firstPageParam = new URLSearchParams(filterParams);
  firstPageParam.append("page", "1");
  const lastPageParam = new URLSearchParams(filterParams);
  lastPageParam.append("page", totalPages.toString());
  const previousPageParam = new URLSearchParams(filterParams);
  previousPageParam.append("page", (loaderData.page - 1).toString());
  const nextPageParam = new URLSearchParams(filterParams);
  nextPageParam.append("page", (loaderData.page + 1).toString());

  const firstPageUrl = `/books/${params.bookID}/expenses?${firstPageParam.toString()}`;
  const lastPageUrl = `/books/${params.bookID}/expenses?${lastPageParam.toString()}`;
  const previousPageUrl =
    loaderData.page > 1
      ? `/books/${params.bookID}/expenses?${previousPageParam.toString()}`
      : "";
  const nextPageUrl =
    loaderData.page < totalPages
      ? `/books/${params.bookID}/expenses?${nextPageParam.toString()}`
      : "";

  return (
    <>
      <title>Expenses | Xpense</title>
      <div className="h-full flex items-center justify-center">
        <div className="h-full max-w-[90rem] flex-1 flex flex-col p-8 gap-4">
          <h1 className="text-4xl">Expenses</h1>
          <div className="flex justify-between items-center">
            <Button className="w-24" asChild>
              <Link to={`/books/${params.bookID}/expenses/create`}>Create</Link>
            </Button>
            <ExpenseFilterForm
              bookId={params.bookID}
              categories={categoryList}
              paymentMethods={paymentMethodList}
              categoryId={loaderData.categoryID}
              paymentMethodId={loaderData.paymentMethodID}
              remark={loaderData.remark}
            />
          </div>
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
