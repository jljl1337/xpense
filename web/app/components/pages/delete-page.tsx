import { useState } from "react";
import { useNavigate } from "react-router";

import { Button } from "~/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";

interface DeletePageProps {
  title: string;
  description: string;
  action: () => Promise<{ error: string | null }>;
  redirectTo: string;
}

export default function DeletePage({
  title,
  description,
  action,
  redirectTo,
}: DeletePageProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  async function onDelete() {
    setIsLoading(true);
    setError(null);

    const response = await action();

    if (response.error != null) {
      setError(response.error);
      setIsLoading(false);
      return;
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
          <div className="flex flex-col gap-4">
            <Button
              variant={"destructive"}
              className="w-full cursor-pointer"
              onClick={onDelete}
              disabled={isLoading}
            >
              {isLoading ? "Deleting..." : "Delete"}
            </Button>
            {error && !isLoading && (
              <div className="text-destructive text-sm text-center">
                {error}
              </div>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
