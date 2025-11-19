import { Link, redirect, useNavigate } from "react-router";
import type { Route } from "./+types/sign-in";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

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

import { createPreSession, getCsrfToken, signIn } from "~/lib/db/auth";
import { getMe } from "~/lib/db/users";

const formSchema = z.object({
  username: z.string().trim().min(1, "Username is required"),
  password: z.string().min(1, "Password is required"),
});

export async function clientLoader() {
  const me = await getMe();

  if (me.data != null) {
    return redirect("/books");
  }

  // Return the CSRF token of the existing pre-session if it is valid
  const existingPreSessionCSRFToken = await getCsrfToken();
  if (existingPreSessionCSRFToken.error == null) {
    console.log("Using existing pre-session CSRF token");
    return {
      data: { preSessionCSRFToken: existingPreSessionCSRFToken.data },
      error: null,
    };
  }

  const preSessionCSRFToken = await createPreSession();
  if (preSessionCSRFToken.error != null) {
    return redirect("/error");
  }

  return {
    data: { preSessionCSRFToken: preSessionCSRFToken.data },
    error: null,
  };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      password: "",
    },
  });

  const {
    setError,
    formState: { isSubmitting, errors },
  } = form;

  const navigate = useNavigate();

  async function onSubmit(values: z.infer<typeof formSchema>) {
    const { error } = await signIn(
      values.username,
      values.password,
      loaderData.data.preSessionCSRFToken,
    );
    if (error) {
      setError("root", {
        message: error,
      });
      return;
    }
    navigate("/books");
  }

  return (
    <>
      <title>Sign In | Xpense</title>
      <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10 bg-background">
        <div className="w-full max-w-sm">
          <div className="flex flex-col gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Sign in to your account</CardTitle>
                <CardDescription>
                  Enter your credentials below to sign in
                </CardDescription>
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
                            <Input placeholder="your_username" {...field} />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name="password"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Password</FormLabel>
                          <FormControl>
                            <Input
                              type="password"
                              placeholder="yourVerySecureP@ssw0rd!"
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <Button
                      type="submit"
                      className="w-full"
                      disabled={isSubmitting}
                    >
                      Submit
                    </Button>
                    {errors.root?.message && !isSubmitting && (
                      <div className="text-destructive text-sm text-center">
                        {errors.root?.message}
                      </div>
                    )}
                    <div className="mt-4 text-center text-sm">
                      Don&apos;t have an account?{" "}
                      <Link
                        to="/auth/sign-up"
                        className="underline underline-offset-4"
                      >
                        Sign up
                      </Link>
                    </div>
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
