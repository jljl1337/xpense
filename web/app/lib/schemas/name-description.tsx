import { z } from "zod";

import { idSchema } from "~/lib/schemas/id";

export const nameDescriptionSchema = z.object({
  name: z
    .string()
    .trim()
    .min(1, "Name must contain at least 1 non-whitespace character"),
  description: z.string().trim(),
});

export const idNameDescriptionSchema = z.intersection(
  nameDescriptionSchema,
  idSchema,
);
