import * as z from "zod";

export const expenseSchema = z.object({
  categoryID: z.string().trim(),
  paymentMethodID: z.string().trim(),
  date: z.string().trim(),
  amount: z.number("Amount is required"),
  remark: z.string().trim(),
});
