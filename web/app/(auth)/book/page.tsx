"use client";

import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { useSearchParams } from "next/navigation";
import { Suspense, useCallback, useState } from "react";
import { Temporal } from "temporal-polyfill";
import { useAuth } from "@/app/_components/AuthProvider";
import { fetchClient } from "@/app/_lib/query";
import type { Book } from "../_models/book";
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

const borrowErrorMessages = {
  CONFLICT: "この書籍は既に貸出中です",
  NOT_FOUND: "書籍が見つかりません",
} as const;

const returnErrorMessages = {
  CONFLICT: "この書籍は貸出中ではありません",
  NOT_FOUND: "書籍が見つかりません",
} as const;

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
        const message =
          (error.code &&
            borrowErrorMessages[
              error.code as keyof typeof borrowErrorMessages
            ]) ??
          error.message ??
          "借りることができませんでした";
        setActionState({ status: "error", message });
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
        const message =
          (error.code &&
            returnErrorMessages[
              error.code as keyof typeof returnErrorMessages
            ]) ??
          error.message ??
          "返却に失敗しました";
        setActionState({ status: "error", message });
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

const loadingBook: Book = {
  id: "",
  title: "読み込み中...",
  authors: [],
  thumbnailUrl: "",
  status: "available",
  createdAt: "",
};

function BookDetailInner() {
  const searchParams = useSearchParams();
  const bookId = searchParams.get("id");
  const authState = useAuth();

  if (!bookId) {
    return <div>書籍IDが指定されていません</div>;
  }

  if (authState.status !== "authenticated") {
    return (
      <BookDetailPage
        book={{ ...loadingBook, id: bookId }}
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
          book={{ ...loadingBook, id: bookId }}
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

export default function BookDetail() {
  return (
    <Suspense
      fallback={
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
      }
    >
      <BookDetailInner />
    </Suspense>
  );
}
