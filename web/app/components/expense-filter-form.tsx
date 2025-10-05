import { useEffect } from "react";
import { useNavigate } from "react-router";

import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "~/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
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

import type { Category } from "~/lib/db/categories";
import type { PaymentMethod } from "~/lib/db/payment-methods";
import { expenseFilterSchema } from "~/lib/schemas/filter";

interface ExpenseFilterFormClientProps {
  bookId: string;
  categories: Category[];
  paymentMethods: PaymentMethod[];
  categoryId?: string;
  paymentMethodId?: string;
  remark?: string;
}

export default function ExpenseFilterFormClient({
  bookId,
  categories,
  paymentMethods,
  categoryId,
  paymentMethodId,
  remark,
}: ExpenseFilterFormClientProps) {
  const navigate = useNavigate();
  const form = useForm<z.infer<typeof expenseFilterSchema>>({
    resolver: zodResolver(expenseFilterSchema),
    defaultValues: {
      categoryId: categoryId,
      paymentMethodId: paymentMethodId,
      remark: remark,
    },
  });

  async function onSubmit(data: z.infer<typeof expenseFilterSchema>) {
    const params = new URLSearchParams();
    if (data.categoryId) {
      params.append("category-id", data.categoryId);
    }
    if (data.paymentMethodId) {
      params.append("payment-method-id", data.paymentMethodId);
    }
    if (data.remark) {
      params.append("remark", data.remark);
    }

    const searchParams = params.toString();
    const url = `/books/${bookId}/expenses?${searchParams}`;
    navigate(url);
  }

  function onReset() {
    const url = `/books/${bookId}/expenses`;
    navigate(url);
  }

  useEffect(() => {
    form.reset({
      categoryId: categoryId,
      paymentMethodId: paymentMethodId,
      remark: remark,
    });
  }, [categoryId, paymentMethodId, remark, form]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <div className="flex gap-4">
          <FormField
            control={form.control}
            name="categoryId"
            render={({ field }) => (
              <FormItem>
                <Select
                  onValueChange={field.onChange}
                  value={field.value || ""}
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Category" />
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
            name="paymentMethodId"
            render={({ field }) => (
              <FormItem>
                <Select
                  onValueChange={field.onChange}
                  value={field.value || ""}
                >
                  <FormControl>
                    <SelectTrigger className="w-full">
                      <SelectValue placeholder="Payment Method" />
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
            name="remark"
            render={({ field }) => (
              <FormItem>
                <FormControl>
                  <Input
                    placeholder="Remark"
                    {...field}
                    value={field.value || ""} // Explicitly handle undefined values
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button variant={"secondary"} type="submit">
            Filter
          </Button>
          {(categoryId || paymentMethodId || remark) && (
            <Button variant={"destructive"} type="button" onClick={onReset}>
              Reset
            </Button>
          )}
        </div>
      </form>
    </Form>
  );
}
