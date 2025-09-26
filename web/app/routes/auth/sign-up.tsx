import { Link, useNavigate } from "react-router";

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

import { signUp } from "~/lib/db/auth";
import { checkUsernameExists } from "~/lib/db/users";
import { passwordWithConfirmSchema, usernameSchema } from "~/lib/schemas/auth";

const formSchema = z.intersection(usernameSchema, passwordWithConfirmSchema);

export default function Page() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: "",
      password: "",
      confirmPassword: "",
    },
  });

  const {
    setError,
    formState: { isSubmitting, errors },
  } = form;

  let navigate = useNavigate();

  async function onSubmit(values: z.infer<typeof formSchema>) {
    const { data: exists, error: checkUsernameError } =
      await checkUsernameExists(values.username);

    if (checkUsernameError) {
      setError("root", {
        message: checkUsernameError,
      });
      return;
    }

    if (exists) {
      setError("username", {
        message: "Username is already taken",
      });
      return;
    }

    const { error } = await signUp(values.username, values.password);
    if (error) {
      setError("root", {
        message: error,
      });
      return;
    }
    navigate("/auth/sign-in");
  }

  return (
    <>
      <title>Sign Up | Xpense</title>
      <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10 bg-background">
        <div className="w-full max-w-sm">
          <div className="flex flex-col gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Sign up for an account</CardTitle>
                <CardDescription>
                  Enter your credentials below to create an account
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
                            <Input
                              placeholder="username_3_to_20_char"
                              {...field}
                            />
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
                              placeholder="P@assw0rd8to32char"
                              {...field}
                            />
                          </FormControl>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                    <FormField
                      control={form.control}
                      name="confirmPassword"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Confirm Password</FormLabel>
                          <FormControl>
                            <Input
                              type="password"
                              placeholder="P@assw0rd8to32char"
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
                      <div className="text-red-500 text-sm text-center">
                        {errors.root?.message}
                      </div>
                    )}
                    <div className="mt-4 text-center text-sm">
                      Already have an account?{" "}
                      <Link
                        to="/auth/sign-in"
                        className="underline underline-offset-4"
                      >
                        Sign in
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
