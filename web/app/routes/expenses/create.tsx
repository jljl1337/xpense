import { redirect } from "react-router";
import type { Route } from "./+types/create";

import type z from "zod";

import ExpensePage from "~/components/pages/expense-page";
import { getCsrfToken } from "~/lib/db/auth";
import { getCategories } from "~/lib/db/categories";
import { createExpense } from "~/lib/db/expenses";
import { getPaymentMethods } from "~/lib/db/payment-methods";
import { dateToYYYYMMDD } from "~/lib/format/date";
import type { expenseSchema } from "~/lib/schemas/expense";

export async function clientLoader({ params }: Route.ClientLoaderArgs) {
  const [csrfToken, categoryList, paymentMethodList] = await Promise.all([
    getCsrfToken(),
    getCategories(params.bookID),
    getPaymentMethods(params.bookID),
  ]);

  if (csrfToken.error != null) {
    return redirect("/error");
  }

  if (categoryList.error != null) {
    return redirect("/error");
  }

  if (paymentMethodList.error != null) {
    return redirect("/error");
  }

  if (categoryList.data.length === 0) {
    return redirect(`/books/${params.bookID}/categories/create`);
  }

  if (paymentMethodList.data.length === 0) {
    return redirect(`/books/${params.bookID}/payment-methods/create`);
  }

  return {
    csrfToken: csrfToken.data,
    categoryList: categoryList.data,
    paymentMethodList: paymentMethodList.data,
  };
}

export default function Page({ loaderData, params }: Route.ComponentProps) {
  const categoryList = loaderData.categoryList;
  const paymentMethodList = loaderData.paymentMethodList;

  const defaultCategoryID = categoryList[0].id;
  const defaultPaymentMethodID = paymentMethodList[0].id;
  const defaultDate = dateToYYYYMMDD(new Date());

  async function action(data: z.infer<typeof expenseSchema>) {
    const response = await createExpense(
      params.bookID,
      data.categoryID,
      data.paymentMethodID,
      data.date,
      data.amount,
      data.remark,
      loaderData.csrfToken,
    );

    if (response.error != null) {
      return { error: response.error };
    }

    return { error: null };
  }

  return (
    <>
      <title>New Expense | Xpense</title>
      <ExpensePage
        categories={categoryList}
        paymentMethods={paymentMethodList}
        title="New Expense"
        description="Create a new expense"
        categoryIdValue={defaultCategoryID}
        paymentMethodIdValue={defaultPaymentMethodID}
        dateValue={defaultDate}
        remarkValue=""
        submitButtonLabel="Create"
        action={action}
      />
    </>
  );
}
