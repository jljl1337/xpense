import { customFetch } from "~/lib/db/fetch";

export type PaymentMethod = {
  id: string;
  bookID: string;
  name: string;
  description: string;
  createdAt: string;
  updatedAt: string;
};

export async function createPaymentMethod(
  bookID: string,
  name: string,
  description: string,
  csrfToken: string,
) {
  const response = await customFetch(
    "/api/payment-methods",
    "POST",
    {
      bookID,
      name,
      description,
    },
    csrfToken,
  );

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}

export async function getPaymentMethods(bookID: string) {
  const response = await customFetch(
    `/api/payment-methods?book-id=${bookID}`,
    "GET",
  );

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: PaymentMethod[] = await response.json();
  return { data, error: null };
}

export async function getPaymentMethod(paymentMethodID: string) {
  const response = await customFetch(
    `/api/payment-methods/${paymentMethodID}`,
    "GET",
  );

  if (!response.ok) {
    const error = await response.text();
    return { data: null, error };
  }

  const data: PaymentMethod = await response.json();
  return { data, error: null };
}

export async function updatePaymentMethod(
  paymentMethodID: string,
  name: string,
  description: string,
  csrfToken: string,
) {
  const response = await customFetch(
    `/api/payment-methods/${paymentMethodID}`,
    "PUT",
    {
      name,
      description,
    },
    csrfToken,
  );

  if (!response.ok) {
    const error = await response.text();
    return { error };
  }

  return { error: null };
}

export async function deletePaymentMethod(
  paymentMethodID: string,
  csrfToken: string,
) {
  const response = await customFetch(
    `/api/payment-methods/${paymentMethodID}`,
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
