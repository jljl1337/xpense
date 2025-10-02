import { redirect } from "react-router";
import type { Route } from "./+types/edit";

import type z from "zod";

import NameDescriptionPage from "~/components/pages/name-description-page";
import { getCsrfToken } from "~/lib/db/auth";
import { getCategory, updateCategory } from "~/lib/db/categories";
import type { nameDescriptionSchema } from "~/lib/schemas/name-description";

export async function clientLoader({ params }: Route.LoaderArgs) {
  const [csrfToken, category] = await Promise.all([
    getCsrfToken(),
    getCategory(params.categoryID),
  ]);

  if (csrfToken.error != null) {
    return redirect("/error");
  }

  if (category.error != null) {
    return redirect("/error");
  }

  return { csrfToken: csrfToken.data, category: category.data };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  async function action(data: z.infer<typeof nameDescriptionSchema>) {
    const response = await updateCategory(
      loaderData.category.id,
      data.name,
      data.description,
      loaderData.csrfToken,
    );

    if (response.error != null) {
      return { error: response.error };
    }

    return { error: null };
  }

  return (
    <>
      <title>Edit Category | Xpense</title>
      <NameDescriptionPage
        title={"Edit Category"}
        description={"Enter the new name and description for your category"}
        nameValue={loaderData.category.name}
        descriptionValue={loaderData.category.description}
        nameFieldLabel={"Category Name"}
        descriptionFieldLabel={"Category Description"}
        nameFieldPlaceholder={"Groceries"}
        descriptionFieldPlaceholder={"Optional"}
        submitButtonLabel={"Update"}
        action={action}
      />
    </>
  );
}
