"use client";

import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { Suspense, use, useCallback, useState } from "react";
import { Temporal } from "temporal-polyfill";
import { useAuth } from "@/app/_components/AuthProvider";
import { fetchClient } from "@/app/_lib/query";
import type { Book } from "../../_models/book";
import { BookDetailPage } from "./page.view";

type ActionState =
  | { status: "idle" }
  | { status: "selecting-date" }
  | { status: "loading"; action: "borrow" | "return" }
  | { status: "error"; message: string };

const getDefaultDueDate = (date = Temporal.Now.plainDateISO()) => {
  return date.add({ days: 7 }).toString();
};

const calculateDueDays = (dueDate: string) => {
  const today = Temporal.Now.plainDateISO();
  const due = Temporal.PlainDate.from(dueDate);
  return today.until(due).days;
};

const getBorrowErrorMessage = (error: { code: string; message?: string }) => {
  if (error.code === "CONFLICT") {
    return "この書籍は既に貸出中です";
  }
  if (error.code === "NOT_FOUND") {
    return "書籍が見つかりません";
  }
  return error.message ?? "借りることができませんでした";
};

const getReturnErrorMessage = (error: { code: string; message?: string }) => {
  if (error.code === "CONFLICT") {
    return "この書籍は貸出中ではありません";
  }
  if (error.code === "NOT_FOUND") {
    return "書籍が見つかりません";
  }
  return error.message ?? "返却に失敗しました";
};

function BookDetailContent({ bookId }: { bookId: string }) {
  const authState = useAuth();
  const queryClient = useQueryClient();
  const [actionState, setActionState] = useState<ActionState>({
    status: "idle",
  });

  const [dueDate, setDueDate] = useState<string>(getDefaultDueDate());

  if (authState.status !== "authenticated") {
    throw new Error("BookDetailContent requires authenticated user");
  }

  const { user } = authState;

  const { data: book } = useSuspenseQuery({
    queryKey: ["books", bookId, user.uid],
    queryFn: async () => {
      const token = await user.getIdToken();
      const { data, error } = await fetchClient.GET("/books/{bookId}", {
        params: { path: { bookId } },
        headers: { Authorization: `Bearer ${token}` },
      });
      if (error) {
        throw new Error(error.message ?? "書籍の取得に失敗しました");
      }
      return data as Book;
    },
  });

  const handleBorrowConfirm = useCallback(async () => {
    setActionState({ status: "loading", action: "borrow" });

    const dueDays = calculateDueDays(dueDate);
    const token = await user.getIdToken();

    try {
      const { error } = await fetchClient.POST("/books/{bookId}/borrow", {
        params: { path: { bookId } },
        body: { dueDays },
        headers: { Authorization: `Bearer ${token}` },
      });

      if (error) {
        setActionState({
          status: "error",
          message: getBorrowErrorMessage(error),
        });
        return;
      }
    } catch (_err) {
      setActionState({ status: "error", message: "通信エラーが発生しました" });
      return;
    }

    await queryClient.invalidateQueries({ queryKey: ["books"] });
    setActionState({ status: "idle" });
  }, [user, bookId, dueDate, queryClient]);

  const handleReturn = useCallback(async () => {
    setActionState({ status: "loading", action: "return" });

    const token = await user.getIdToken();

    try {
      const { error } = await fetchClient.POST("/books/{bookId}/return", {
        params: { path: { bookId } },
        headers: { Authorization: `Bearer ${token}` },
      });

      if (error) {
        setActionState({
          status: "error",
          message: getReturnErrorMessage(error),
        });
        return;
      }
    } catch (_err) {
      setActionState({ status: "error", message: "通信エラーが発生しました" });
      return;
    }

    await queryClient.invalidateQueries({ queryKey: ["books"] });
    setActionState({ status: "idle" });
  }, [user, bookId, queryClient]);

  return (
    <BookDetailPage
      book={book}
      actionState={actionState}
      currentUserId={user.uid}
      dueDate={dueDate}
      onBorrowClick={() => {
        setActionState({ status: "selecting-date" });
        setDueDate(getDefaultDueDate());
      }}
      onBorrowConfirm={handleBorrowConfirm}
      onBorrowCancel={() => setActionState({ status: "idle" })}
      onDueDateChange={setDueDate}
      onReturn={handleReturn}
    />
  );
}

export default function BookDetail({
  params,
}: {
  params: Promise<{ bookId: string }>;
}) {
  const { bookId } = use(params);
  const authState = useAuth();

  const loadingBook: Book = {
    id: bookId,
    title: "読み込み中...",
    authors: [],
    thumbnailUrl: "",
    status: "available",
    createdAt: "",
  };

  if (authState.status !== "authenticated") {
    return (
      <BookDetailPage
        book={loadingBook}
        actionState={{ status: "idle" }}
        currentUserId=""
        dueDate=""
        onBorrowClick={() => {}}
        onBorrowConfirm={() => {}}
        onBorrowCancel={() => {}}
        onDueDateChange={() => {}}
        onReturn={() => {}}
      />
    );
  }

  return (
    <Suspense
      fallback={
        <BookDetailPage
          book={loadingBook}
          actionState={{ status: "idle" }}
          currentUserId={authState.user.uid}
          dueDate=""
          onBorrowClick={() => {}}
          onBorrowConfirm={() => {}}
          onBorrowCancel={() => {}}
          onDueDateChange={() => {}}
          onReturn={() => {}}
        />
      }
    >
      <BookDetailContent bookId={bookId} />
    </Suspense>
  );
}
