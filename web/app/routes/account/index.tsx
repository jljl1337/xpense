import { useState } from "react";
import { Link, redirect, useNavigate } from "react-router";
import type { Route } from "./+types";

import { LogOut, Trash2 } from "lucide-react";

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

import { getCsrfToken, signOut, signOutAll } from "~/lib/db/auth";
import { redirectIfNeeded } from "~/lib/db/common";
import { deleteMe, getMe } from "~/lib/db/users";

export async function clientLoader() {
  const [user, csrfToken] = await Promise.all([getMe(), getCsrfToken()]);

  redirectIfNeeded(user.error, csrfToken.error);

  return {
    user: user.data!,
    csrfToken: csrfToken.data!,
  };
}

export default function Page({ loaderData }: Route.ComponentProps) {
  const navigate = useNavigate();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function onSignOut() {
    setIsLoading(true);
    setError(null);

    const { error } = await signOut(loaderData.csrfToken);
    if (error) {
      setError(error);
      return;
    }

    setIsLoading(false);

    navigate("/auth/sign-in");
  }

  async function onSignOutAll() {
    setIsLoading(true);
    setError(null);

    const { error } = await signOutAll(loaderData.csrfToken);
    if (error) {
      setError(error);
      return;
    }

    setIsLoading(false);

    navigate("/auth/sign-in");
  }

  async function onDeleteAccount() {
    setIsLoading(true);
    setError(null);

    const { error } = await deleteMe(loaderData.csrfToken);
    if (error) {
      setError(error);
      return;
    }

    setIsLoading(false);

    navigate("/auth/sign-in");
  }

  return (
    <>
      <title>Account | Xpense</title>
      <div className="h-full flex items-center justify-center">
        <div className="h-full max-w-[90rem] flex-1 flex flex-col p-8 gap-4">
          <h1 className="text-4xl">Account</h1>
          <div>
            <p className="mb-2">User ID: {loaderData.user.id}</p>
            <p className="mb-2">Username: {loaderData.user.username}</p>
            <p className="mb-2">
              Created At: {new Date(loaderData.user.createdAt).toLocaleString()}
            </p>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>General</CardTitle>
              <CardDescription>Manage your account settings</CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <Separator />
              {/* Change Username */}
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label className="text-base">Change Username</Label>
                  <p className="text-muted-foreground text-sm">
                    Change your account username
                  </p>
                </div>
                <Button asChild>
                  <Link to="/account/change-username">Change Username</Link>
                </Button>
              </div>
              <Separator />
              {/* Change Password */}
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label className="text-base">Change Password</Label>
                  <p className="text-muted-foreground text-sm">
                    Change your account password
                  </p>
                </div>
                <Button asChild>
                  <Link to="/account/change-password">Change Password</Link>
                </Button>
              </div>
              <Separator />
              {/* Sign Out */}
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label className="text-base">Sign Out</Label>
                  <p className="text-muted-foreground text-sm">
                    Sign out of your account on this device
                  </p>
                </div>
                <Dialog>
                  <DialogTrigger asChild>
                    <Button variant="destructive" className="cursor-pointer">
                      <LogOut className="mr-2 h-4 w-4" />
                      Sign Out
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="sm:max-w-[425px]">
                    <DialogHeader>
                      <DialogTitle>Sign Out</DialogTitle>
                      <DialogDescription>
                        Are you sure you want to sign out of your account?
                      </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                      <DialogClose asChild>
                        <Button variant="outline" className="cursor-pointer">
                          Cancel
                        </Button>
                      </DialogClose>
                      <Button
                        variant="destructive"
                        className="cursor-pointer"
                        onClick={onSignOut}
                        disabled={isLoading}
                      >
                        Yes
                      </Button>
                      {error && !isLoading && (
                        <p className="text-destructive text-sm mt-2">{error}</p>
                      )}
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
              <Separator />
              {/* Sign Out (all devices) */}
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label className="text-base">Sign Out (all devices)</Label>
                  <p className="text-muted-foreground text-sm">
                    Sign out of your account on all devices
                  </p>
                </div>
                <Dialog>
                  <DialogTrigger asChild>
                    <Button variant="destructive" className="cursor-pointer">
                      <LogOut className="mr-2 h-4 w-4" />
                      Sign Out (all devices)
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="sm:max-w-[425px]">
                    <DialogHeader>
                      <DialogTitle>Sign Out</DialogTitle>
                      <DialogDescription>
                        Are you sure you want to sign out of your account on all
                        devices?
                      </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                      <DialogClose asChild>
                        <Button variant="outline" className="cursor-pointer">
                          Cancel
                        </Button>
                      </DialogClose>
                      <Button
                        variant="destructive"
                        className="cursor-pointer"
                        onClick={onSignOutAll}
                        disabled={isLoading}
                      >
                        Yes
                      </Button>
                      {error && !isLoading && (
                        <p className="text-destructive text-sm mt-2">{error}</p>
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
              <Separator />
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <Label className="text-base">Delete Account</Label>
                  <p className="text-muted-foreground text-sm">
                    Permanently delete your account and all data
                  </p>
                </div>
                <Dialog>
                  <DialogTrigger asChild>
                    <Button variant="destructive" className="cursor-pointer">
                      <Trash2 className="mr-2 h-4 w-4" />
                      Delete Account
                    </Button>
                  </DialogTrigger>
                  <DialogContent className="sm:max-w-[425px]">
                    <DialogHeader>
                      <DialogTitle>Delete Account</DialogTitle>
                      <DialogDescription>
                        Are you sure you want to delete your account? This
                        action is irreversible and will permanently delete all
                        your data.
                      </DialogDescription>
                    </DialogHeader>
                    <DialogFooter>
                      <DialogClose asChild>
                        <Button variant="outline" className="cursor-pointer">
                          Cancel
                        </Button>
                      </DialogClose>
                      <Button
                        variant="destructive"
                        className="cursor-pointer"
                        onClick={onDeleteAccount}
                        disabled={isLoading}
                      >
                        Yes
                      </Button>
                      {error && !isLoading && (
                        <p className="text-destructive text-sm mt-2">{error}</p>
                      )}
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </>
  );
}
