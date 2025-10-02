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
}

export default function DeletePage({
  title,
  description,
  action,
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

    navigate(-1);
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
            <div className="flex justify-end gap-2">
              <Button
                variant={"outline"}
                className="cursor-pointer"
                onClick={() => navigate(-1)}
              >
                Cancel
              </Button>
              <Button
                variant={"destructive"}
                className="cursor-pointer"
                onClick={onDelete}
                disabled={isLoading}
              >
                Delete
              </Button>
            </div>
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
