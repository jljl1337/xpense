import { Link, redirect } from "react-router";
import type { Route } from "./+types";

import type { ColumnDef } from "@tanstack/react-table";

import { Button } from "~/components/ui/button";

import Pagination from "~/components/pagination";
import { DataTable } from "~/components/tables/data-table";
import TableRowDropdown from "~/components/tables/dropdown";
import { getBooks, getBooksCount, type Book } from "~/lib/db/books";
import { isUnauthorizedError } from "~/lib/db/common";

export async function clientLoader({ request }: Route.ClientLoaderArgs) {
  const searchParams = new URL(request.url).searchParams;
  let page = parseInt(searchParams.get("page") ?? "1");
  if (isNaN(page) || page < 1) {
    page = 1;
  }

  const [count, bookList] = await Promise.all([
    getBooksCount(),
    getBooks(page, 20),
  ]);

  if (count.error != null) {
    if (isUnauthorizedError(count.error)) {
      return redirect("/auth/sign-in");
    }
    return redirect("/error");
  }
  if (bookList.error != null) {
    if (isUnauthorizedError(bookList.error)) {
      return redirect("/auth/sign-in");
    }
    return redirect("/error");
  }

  return {
    count: count.data,
    page,
    books: bookList.data,
  };
}

const columns: ColumnDef<Book>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => {
      return (
        <Link to={`/books/${row.original.id}/expenses`} className="underline">
          {row.original.name}
        </Link>
      );
    },
  },
  {
    accessorKey: "description",
    header: "Description",
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
            editUrl={`/books/${row.original.id}/edit`}
            deleteUrl={`/books/${row.original.id}/delete`}
          />
        </div>
      );
    },
  },
];

export default function Page({ loaderData }: Route.ComponentProps) {
  const totalPages = Math.ceil(loaderData.count / 20);
  const firstPageUrl = `/books?page=1`;
  const lastPageUrl = `/books?page=${totalPages}`;
  const previousPageUrl =
    loaderData.page > 1 ? `/books?page=${loaderData.page - 1}` : "";
  const nextPageUrl =
    loaderData.page < totalPages ? `/books?page=${loaderData.page + 1}` : "";

  return (
    <>
      <title>Books | Xpense</title>
      <div className="h-full flex items-center justify-center">
        <div className="h-full max-w-[90rem] flex-1 flex flex-col p-8 gap-4">
          <h1 className="text-4xl">Books</h1>
          <Button className="w-24" asChild>
            <Link to="/books/create">Create</Link>
          </Button>
          <DataTable columns={columns} data={loaderData.books} />
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
