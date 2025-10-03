import { customFetch } from "~/lib/db/fetch";

type ExpenseCount = {
  count: number;
};

export type Expense = {
  id: string;
  bookID: string;
  categoryID: string;
  paymentMethodID: string;
  date: string;
  amount: number;
  remark: string;
  createdAt: string;
  updatedAt: string;
};

export async function createExpense(
  bookID: string,
  categoryID: string,
  paymentMethodID: string,
  date: string,
  amount: number,
  remark: string,
  csrfToken: string,
) {
  const response = await customFetch(
    "/api/expenses",
    "POST",
    {
      bookID,
      categoryID,
      paymentMethodID,
      date,
      amount,
      remark,
    },
    csrfToken,
  );

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}

export async function getExpensesCountByBookID(bookID: string) {
  const response = await customFetch(
    `/api/expenses/count?book-id=${bookID}`,
    "GET",
  );

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: ExpenseCount = await response.json();
  return { data: data.count, error: null };
}

export async function getExpensesByBookID(
  bookID: string,
  page: number,
  pageSize: number,
) {
  const response = await customFetch(
    `/api/expenses?book-id=${bookID}&page=${page}&page-size=${pageSize}`,
    "GET",
  );

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: Expense[] = await response.json();
  return { data, error: null };
}

export async function getExpense(expenseID: string) {
  const response = await customFetch(`/api/expenses/${expenseID}`, "GET");

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: Expense = await response.json();
  return { data, error: null };
}

export async function updateExpense(
  expenseID: string,
  categoryID: string,
  paymentMethodID: string,
  date: string,
  amount: number,
  remark: string,
  csrfToken: string,
) {
  const response = await customFetch(
    `/api/expenses/${expenseID}`,
    "PUT",
    {
      categoryID,
      paymentMethodID,
      date,
      amount,
      remark,
    },
    csrfToken,
  );

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}
export async function deleteExpense(expenseID: string, csrfToken: string) {
  const response = await customFetch(
    `/api/expenses/${expenseID}`,
    "DELETE",
    null,
    csrfToken,
  );

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}
