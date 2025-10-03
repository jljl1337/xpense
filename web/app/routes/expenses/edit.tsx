import { redirect } from "react-router";
import type { Route } from "./+types/edit";

import type z from "zod";

import ExpensePage from "~/components/pages/expense-page";
import { getCsrfToken } from "~/lib/db/auth";
import { getCategories } from "~/lib/db/categories";
import { createExpense, getExpense, updateExpense } from "~/lib/db/expenses";
import { getPaymentMethods } from "~/lib/db/payment-methods";
import { dateToYYYYMMDD } from "~/lib/format/date";
import type { expenseSchema } from "~/lib/schemas/expense";

export async function clientLoader({ params }: Route.ClientLoaderArgs) {
  const [csrfToken, categoryList, paymentMethodList, expense] =
    await Promise.all([
      getCsrfToken(),
      getCategories(params.bookID),
      getPaymentMethods(params.bookID),
      getExpense(params.expenseID),
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

  if (expense.error != null) {
    return redirect("/error");
  }

  return {
    csrfToken: csrfToken.data,
    categoryList: categoryList.data,
    paymentMethodList: paymentMethodList.data,
    expense: expense.data,
  };
}

export default function Page({ loaderData, params }: Route.ComponentProps) {
  const categoryList = loaderData.categoryList;
  const paymentMethodList = loaderData.paymentMethodList;
  const expense = loaderData.expense;

  const defaultCategoryID = expense.categoryID;
  const defaultPaymentMethodID = expense.paymentMethodID;
  const defaultDate = expense.date;

  async function action(data: z.infer<typeof expenseSchema>) {
    const response = await updateExpense(
      params.expenseID,
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
      <title>Edit Expense | Xpense</title>
      <ExpensePage
        categories={categoryList}
        paymentMethods={paymentMethodList}
        title="Edit Expense"
        description="Edit an existing expense"
        categoryIdValue={defaultCategoryID}
        paymentMethodIdValue={defaultPaymentMethodID}
        dateValue={defaultDate}
        amountValue={expense.amount}
        remarkValue={expense.remark}
        submitButtonLabel="Update"
        action={action}
      />
    </>
  );
}
