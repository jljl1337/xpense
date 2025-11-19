import { useNavigate } from "react-router";
import type { Route } from "./+types/change-username";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import type z from "zod";

import { Button } from "~/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "~/components/ui/form";
import { Input } from "~/components/ui/input";

import { getCsrfToken } from "~/lib/db/auth";
import { redirectIfNeeded } from "~/lib/db/common";
import { updateUsername } from "~/lib/db/users";
import { usernameSchema } from "~/lib/schemas/auth";

export async function clientLoader() {
  const csrfToken = await getCsrfToken();
  redirectIfNeeded(csrfToken.error);

  return { csrfToken: csrfToken.data! };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  const form = useForm<z.infer<typeof usernameSchema>>({
    resolver: zodResolver(usernameSchema),
    defaultValues: {
      username: "",
    },
  });

  const {
    setError,
    formState: { isSubmitting, errors },
  } = form;

  const navigate = useNavigate();

  async function onSubmit(values: z.infer<typeof usernameSchema>) {
    const { error } = await updateUsername(
      values.username,
      loaderData.csrfToken,
    );
    if (error) {
      setError("root", {
        message: error,
      });
      return;
    }

    navigate("/account");
  }

  return (
    <>
      <title>Change Username | Xpense</title>
      <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10 bg-background">
        <div className="w-full max-w-sm">
          <div className="flex flex-col gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Change Username</CardTitle>
                <CardDescription>Change your account username</CardDescription>
              </CardHeader>
              <CardContent>
                <Form {...form}>
                  <form
                    onSubmit={form.handleSubmit(onSubmit)}
                    className="space-y-4"
                  >
                    <FormField
                      control={form.control}
                      name="username"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Username</FormLabel>
                          <FormControl>
                            <Input placeholder="your_new_username" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <Button
                      type="submit"
                      className="w-full cursor-pointer"
                      disabled={isSubmitting}
                    >
                      Save
                    </Button>
                    {errors.root?.message && !isSubmitting && (
                      <div className="text-destructive text-sm text-center">
                        {errors.root?.message}
                      </div>
                    )}
                  </form>
                </Form>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </>
  );
}
