import { useState } from "react";

import type { Route } from "./+types/account";
import { LogOut, Trash2 } from "lucide-react";
import { redirect, useNavigate } from "react-router";

import { Button } from "~/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "~/components/ui/dialog";
import { Label } from "~/components/ui/label";
import { Separator } from "~/components/ui/separator";

import { getCsrfToken, logout, logoutAll } from "~/lib/db/auth";
import { deleteMe, getMe } from "~/lib/db/user";

export async function clientLoader() {
  const [user, csrfToken] = await Promise.all([getMe(), getCsrfToken()]);

  if (user.error != null) {
    if (user.error.trim() == "Unauthorized") {
      return redirect("/auth/login");
    }
    return { data: null, error: user.error };
  }
  if (csrfToken.error != null) {
    if (csrfToken.error.trim() == "Unauthorized") {
      return redirect("/auth/login");
    }
    return { data: null, error: csrfToken.error };
  }

  return {
    data: { user: user.data, csrfToken: csrfToken.data },
    error: null,
  };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onLogout() {
    if (loaderData.error != null) {
      return;
    }

    setIsLoading(true);
    setError(null);

    const { error } = await logout(loaderData.data.csrfToken);
    if (error) {
      console.error("Failed to log out:", error);
      setError(error);
      return;
    }

    setIsLoading(false);

    navigate("/auth/login");
  }

  async function onLogoutAll() {
    if (loaderData.error != null) {
      return;
    }

    setIsLoading(true);
    setError(null);

    const { error } = await logoutAll(loaderData.data.csrfToken);
    if (error) {
      console.error("Failed to log out of all devices:", error);
      setError(error);
      return;
    }

    setIsLoading(false);

    navigate("/auth/login");
  }

  async function onDeleteAccount() {
    if (loaderData.error != null) {
      return;
    }

    setIsLoading(true);
    setError(null);

    const { error } = await deleteMe(loaderData.data.csrfToken);
    if (error) {
      console.error("Failed to delete account:", error);
      setError(error);
      return;
    }

    setIsLoading(false);

    navigate("/auth/login");
  }

  return (
    <div className="h-full flex items-center justify-center">
      <div className="h-full max-w-[90rem] flex-1 flex flex-col p-8 gap-4">
        <h1 className="text-2xl font-bold mb-4">Account</h1>
        {loaderData?.data ? (
          <div>
            <p className="mb-2">User ID: {loaderData.data.user.id}</p>
            <p className="mb-2">Username: {loaderData.data.user.username}</p>
            <p className="mb-2">
              Created At:{" "}
              {new Date(loaderData.data.user.createdAt).toLocaleString()}
            </p>
          </div>
        ) : (
          <p className="text-red-500">
            Error: {loaderData?.error || "Failed to load user data."}
          </p>
        )}

        <Card>
          <CardHeader>
            <CardTitle>General</CardTitle>
            <CardDescription>Manage your account settings</CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            {/* Log Out */}
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <Label className="text-base">Log Out</Label>
                <p className="text-muted-foreground text-sm">
                  Log out of your account on this device
                </p>
              </div>
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="destructive">
                    <LogOut className="mr-2 h-4 w-4" />
                    Log Out
                  </Button>
                </DialogTrigger>
                <DialogContent className="sm:max-w-[425px]">
                  <DialogHeader>
                    <DialogTitle>Log Out</DialogTitle>
                    <DialogDescription>
                      Are you sure you want to log out of your account?
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter>
                    <DialogClose asChild>
                      <Button variant="outline">Cancel</Button>
                    </DialogClose>
                    <Button
                      variant="destructive"
                      onClick={onLogout}
                      disabled={isLoading}
                    >
                      Yes
                    </Button>
                    {error && !isLoading && (
                      <p className="text-red-500 text-sm mt-2">{error}</p>
                    )}
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </div>
            <Separator />
            {/* Log Out (all devices) */}
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <Label className="text-base">Log Out (all devices)</Label>
                <p className="text-muted-foreground text-sm">
                  Log out of your account on all devices
                </p>
              </div>
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="destructive">
                    <LogOut className="mr-2 h-4 w-4" />
                    Log Out (all devices)
                  </Button>
                </DialogTrigger>
                <DialogContent className="sm:max-w-[425px]">
                  <DialogHeader>
                    <DialogTitle>Log Out</DialogTitle>
                    <DialogDescription>
                      Are you sure you want to log out of your account on all
                      devices?
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter>
                    <DialogClose asChild>
                      <Button variant="outline">Cancel</Button>
                    </DialogClose>
                    <Button
                      variant="destructive"
                      onClick={onLogoutAll}
                      disabled={isLoading}
                    >
                      Yes
                    </Button>
                    {error && !isLoading && (
                      <p className="text-red-500 text-sm mt-2">{error}</p>
                    )}
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </div>
          </CardContent>
        </Card>

        <Card className="border-destructive/50">
          <CardHeader>
            <CardTitle className="text-destructive">Danger Zone</CardTitle>
            <CardDescription>
              Irreversible and destructive actions
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-6">
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <Label className="text-base">Delete Account</Label>
                <p className="text-muted-foreground text-sm">
                  Permanently delete your account and all data
                </p>
              </div>
              <Dialog>
                <DialogTrigger asChild>
                  <Button variant="destructive">
                    <Trash2 className="mr-2 h-4 w-4" />
                    Delete Account
                  </Button>
                </DialogTrigger>
                <DialogContent className="sm:max-w-[425px]">
                  <DialogHeader>
                    <DialogTitle>Delete Account</DialogTitle>
                    <DialogDescription>
                      Are you sure you want to delete your account? This action
                      is irreversible and will permanently delete all your data.
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter>
                    <DialogClose asChild>
                      <Button variant="outline">Cancel</Button>
                    </DialogClose>
                    <Button
                      variant="destructive"
                      onClick={onDeleteAccount}
                      disabled={isLoading}
                    >
                      Yes
                    </Button>
                    {error && !isLoading && (
                      <p className="text-red-500 text-sm mt-2">{error}</p>
                    )}
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
