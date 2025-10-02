import { Link, redirect } from "react-router";
import type { Route } from "./+types";

import type { ColumnDef } from "@tanstack/react-table";

import { Button } from "~/components/ui/button";

import { DataTable } from "~/components/tables/data-table";
import TableRowDropdown from "~/components/tables/dropdown";
import {
  getPaymentMethods,
  type PaymentMethod,
} from "~/lib/db/payment-methods";

export async function clientLoader({ params }: Route.ClientLoaderArgs) {
  const paymentMethodList = await getPaymentMethods(params.bookID);

  if (paymentMethodList.error != null) {
    return redirect("/error");
  }

  return {
    paymentMethods: paymentMethodList.data,
  };
}

const columns: ColumnDef<PaymentMethod>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => {
      return (
        <Link
          to={`/books/${row.original.bookID}/payment-methods/${row.original.id}/edit`}
          className="underline"
        >
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
            editUrl={`/books/${row.original.bookID}/payment-methods/${row.original.id}/edit`}
            deleteUrl={`/books/${row.original.bookID}/payment-methods/${row.original.id}/delete`}
          />
        </div>
      );
    },
  },
];

export default function Page({ params, loaderData }: Route.ComponentProps) {
  return (
    <>
      <title>Payment Methods | Xpense</title>
      <div className="h-full flex items-center justify-center">
        <div className="h-full max-w-[90rem] flex-1 flex flex-col p-8 gap-4">
          <h1 className="text-4xl">Payment Methods</h1>
          <Button className="w-24" asChild>
            <Link to={`/books/${params.bookID}/payment-methods/create`}>
              Create
            </Link>
          </Button>
          <DataTable columns={columns} data={loaderData.paymentMethods} />
        </div>
      </div>
    </>
  );
}
