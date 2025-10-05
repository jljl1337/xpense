import { z } from "zod";

export const expenseFilterSchema = z.object({
  categoryId: z.string().trim().optional(),
  paymentMethodId: z.string().trim().optional(),
  remark: z.string().trim().optional(),
});
