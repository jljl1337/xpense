import { useNavigate } from "react-router";

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

import { nameDescriptionSchema } from "~/lib/schemas/name-description";

interface NameDescriptionPageProps {
  title: string;
  description: string;
  nameValue: string;
  descriptionValue: string;
  nameFieldLabel: string;
  descriptionFieldLabel: string;
  nameFieldPlaceholder: string;
  descriptionFieldPlaceholder: string;
  submitButtonLabel: string;
  action: (
    data: z.infer<typeof nameDescriptionSchema>,
  ) => Promise<{ error: string | null }>;
  redirectTo: string;
}

export default function NameDescriptionPage({
  title,
  description,
  nameValue,
  descriptionValue,
  nameFieldLabel,
  descriptionFieldLabel,
  nameFieldPlaceholder,
  descriptionFieldPlaceholder,
  submitButtonLabel,
  action,
  redirectTo,
}: NameDescriptionPageProps) {
  const form = useForm<z.infer<typeof nameDescriptionSchema>>({
    resolver: zodResolver(nameDescriptionSchema),
    defaultValues: {
      name: nameValue,
      description: descriptionValue,
    },
  });

  const {
    setError,
    formState: { isSubmitting, errors },
  } = form;

  const navigate = useNavigate();

  async function onSubmit(data: z.infer<typeof nameDescriptionSchema>) {
    const response = await action(data);

    if (response.error != null) {
      setError("root", {
        message: response.error,
      });
    }

    navigate(redirectTo);
  }

  return (
    <div className="flex h-full w-full items-center justify-center">
      <Card className="m-4 w-full max-w-sm">
        <CardHeader>
          <CardTitle className="text-2xl">{title}</CardTitle>
          <CardDescription>{description}</CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <div className="min-h-30">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>{nameFieldLabel}</FormLabel>
                      <FormControl>
                        <Input
                          placeholder={nameFieldPlaceholder}
                          disabled={isSubmitting}
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <FormField
                control={form.control}
                name="description"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{descriptionFieldLabel}</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={descriptionFieldPlaceholder}
                        disabled={isSubmitting}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <div className="h-12 flex items-center justify-center">
                {errors.root?.message && !isSubmitting && (
                  <div className="text-red-500 text-sm text-center">
                    {errors.root?.message}
                  </div>
                )}
              </div>
              <div className="flex justify-end">
                <Button
                  className="cursor-pointer"
                  type="submit"
                  disabled={isSubmitting}
                >
                  {submitButtonLabel}
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
