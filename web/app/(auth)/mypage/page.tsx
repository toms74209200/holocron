"use client";

import { useSuspenseQuery } from "@tanstack/react-query";
import { Suspense } from "react";
import { useAuth } from "@/app/_components/AuthProvider";
import { fetchClient } from "@/app/_lib/query";
import type { Borrowing } from "../_models/borrowing";
import { MyPage } from "./page.view";

function MyPageContent() {
  const authState = useAuth();

  if (authState.status !== "authenticated") {
    throw new Error("MyPageContent requires authenticated user");
  }

  const { user } = authState;

  const { data: borrowings } = useSuspenseQuery({
    queryKey: ["borrowings", user.uid],
    queryFn: async () => {
      const token = await user.getIdToken();
      const { data, error } = await fetchClient.GET("/users/me/borrowings", {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (error) {
        throw new Error("借りている本の取得に失敗しました");
      }
      return (data?.items ?? []) as Borrowing[];
    },
  });

  return <MyPage borrowings={borrowings} />;
}

const emptyBorrowings: Borrowing[] = [];

function MyPageInner() {
  const authState = useAuth();

  if (authState.status !== "authenticated") {
    return <MyPage borrowings={emptyBorrowings} isLoading />;
  }

  return (
    <Suspense fallback={<MyPage borrowings={emptyBorrowings} isLoading />}>
      <MyPageContent />
    </Suspense>
  );
}

export default function MyPageRoute() {
  return (
    <Suspense fallback={<MyPage borrowings={emptyBorrowings} isLoading />}>
      <MyPageInner />
    </Suspense>
  );
}
