import { useNavigate } from "react-router";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import * as z from "zod";

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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "~/components/ui/select";

import { DatePicker } from "~/components/date-picker";
import type { Category } from "~/lib/db/categories";
import type { PaymentMethod } from "~/lib/db/payment-methods";
import { dateToYYYYMMDD, YYYYMMDDToDate } from "~/lib/format/date";
import { expenseSchema } from "~/lib/schemas/expense";

interface ExpensePageProps {
  categories: Category[];
  paymentMethods: PaymentMethod[];
  title: string;
  description: string;
  categoryIdValue: string;
  paymentMethodIdValue: string;
  amountValue?: number;
  dateValue: string;
  remarkValue: string;
  submitButtonLabel: string;
  action: (
    data: z.infer<typeof expenseSchema>,
  ) => Promise<{ error: string | null }>;
}

export default function ExpensePage({
  categories,
  paymentMethods,
  title,
  description,
  categoryIdValue,
  paymentMethodIdValue,
  amountValue,
  dateValue,
  remarkValue,
  submitButtonLabel,
  action,
}: ExpensePageProps) {
  const navigate = useNavigate();

  const form = useForm<z.infer<typeof expenseSchema>>({
    resolver: zodResolver(expenseSchema),
    defaultValues: {
      categoryID: categoryIdValue,
      paymentMethodID: paymentMethodIdValue,
      date: dateValue,
      amount: amountValue,
      remark: remarkValue,
    },
  });

  const {
    setError,
    formState: { isSubmitting, errors },
  } = form;

  async function onSubmit(data: z.infer<typeof expenseSchema>) {
    const response = await action(data);

    if (response.error != null) {
      setError("root", {
        message: response.error,
      });
    }

    navigate(-1);
  }

  return (
    <div className="flex h-full w-full items-center justify-center">
      <Card className="m-4">
        <CardHeader>
          <CardTitle className="text-2xl">{title}</CardTitle>
          <CardDescription>{description}</CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              <FormField
                control={form.control}
                name="date"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Date</FormLabel>
                    <FormControl>
                      <DatePicker
                        setDate={(date) => {
                          if (date === undefined) {
                            return;
                          }
                          field.onChange(dateToYYYYMMDD(date));
                        }}
                        date={YYYYMMDDToDate(field.value)}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="categoryID"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Category</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger className="w-full">
                          <SelectValue placeholder="Select a category" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {categories.map((category) => (
                          <SelectItem key={category.id} value={category.id}>
                            {category.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="paymentMethodID"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Payment Method</FormLabel>
                    <Select
                      onValueChange={field.onChange}
                      defaultValue={field.value}
                    >
                      <FormControl>
                        <SelectTrigger className="w-full">
                          <SelectValue placeholder="Select a payment method" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent>
                        {paymentMethods.map((paymentMethod) => (
                          <SelectItem
                            key={paymentMethod.id}
                            value={paymentMethod.id}
                          >
                            {paymentMethod.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="amount"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Amount</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        disabled={isSubmitting}
                        value={field.value}
                        onChange={(e) => {
                          if (e.target.value.trim() === "") {
                            return field.onChange(undefined);
                          }
                          return field.onChange(Number(e.target.value));
                        }}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="remark"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Remark</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="Optional"
                        disabled={isSubmitting}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <Button type="submit" className="w-full" disabled={isSubmitting}>
                {submitButtonLabel}
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
  );
}
