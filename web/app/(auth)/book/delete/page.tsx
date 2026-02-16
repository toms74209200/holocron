"use client";

import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useCallback, useState } from "react";
import { useAuth } from "@/app/_components/AuthProvider";
import { fetchClient } from "@/app/_lib/query";
import type { Book } from "../../_models/book";
import { DeleteBookPage } from "./page.view";

type DeleteReason = "transfer" | "disposal" | "lost" | "other";

type DeleteState =
  | { status: "idle" }
  | { status: "deleting" }
  | { status: "error"; message: string };

const deleteErrorMessages = {
  CONFLICT: "貸出中の書籍は削除できません",
  NOT_FOUND: "書籍が見つかりません",
} as const;

function DeleteBookContent({ bookId }: { bookId: string }) {
  const authState = useAuth();
  const queryClient = useQueryClient();
  const router = useRouter();
  const [deleteState, setDeleteState] = useState<DeleteState>({
    status: "idle",
  });
  const [deleteReason, setDeleteReason] = useState<DeleteReason>("transfer");
  const [deleteMemo, setDeleteMemo] = useState<string>("");

  if (authState.status !== "authenticated") {
    throw new Error("DeleteBookContent requires authenticated user");
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

  const handleDelete = useCallback(async () => {
    setDeleteState({ status: "deleting" });

    const token = await user.getIdToken();

    try {
      const { error } = await fetchClient.DELETE("/books/{bookId}", {
        params: { path: { bookId } },
        body: { reason: deleteReason, memo: deleteMemo || undefined },
        headers: { Authorization: `Bearer ${token}` },
      });

      if (error) {
        const message =
          (error.code &&
            deleteErrorMessages[
              error.code as keyof typeof deleteErrorMessages
            ]) ??
          error.message ??
          "削除に失敗しました";
        setDeleteState({ status: "error", message });
        return;
      }
    } catch (_err) {
      setDeleteState({ status: "error", message: "通信エラーが発生しました" });
      return;
    }

    await queryClient.invalidateQueries({ queryKey: ["books"] });
    router.push("/");
  }, [user, bookId, deleteReason, deleteMemo, queryClient, router]);

  const handleCancel = useCallback(() => {
    router.back();
  }, [router]);

  return (
    <DeleteBookPage
      book={book}
      deleteState={deleteState}
      deleteReason={deleteReason}
      deleteMemo={deleteMemo}
      onDeleteReasonChange={setDeleteReason}
      onDeleteMemoChange={setDeleteMemo}
      onDelete={handleDelete}
      onCancel={handleCancel}
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

function DeleteBookInner() {
  const searchParams = useSearchParams();
  const bookId = searchParams.get("id");
  const authState = useAuth();

  if (!bookId) {
    return <div>書籍IDが指定されていません</div>;
  }

  if (authState.status !== "authenticated") {
    return (
      <DeleteBookPage
        book={{ ...loadingBook, id: bookId }}
        deleteState={{ status: "idle" }}
        deleteReason="transfer"
        deleteMemo=""
        onDeleteReasonChange={() => {}}
        onDeleteMemoChange={() => {}}
        onDelete={() => {}}
        onCancel={() => {}}
      />
    );
  }

  return (
    <Suspense
      fallback={
        <DeleteBookPage
          book={{ ...loadingBook, id: bookId }}
          deleteState={{ status: "idle" }}
          deleteReason="transfer"
          deleteMemo=""
          onDeleteReasonChange={() => {}}
          onDeleteMemoChange={() => {}}
          onDelete={() => {}}
          onCancel={() => {}}
        />
      }
    >
      <DeleteBookContent bookId={bookId} />
    </Suspense>
  );
}

export default function DeleteBook() {
  return (
    <Suspense
      fallback={
        <DeleteBookPage
          book={loadingBook}
          deleteState={{ status: "idle" }}
          deleteReason="transfer"
          deleteMemo=""
          onDeleteReasonChange={() => {}}
          onDeleteMemoChange={() => {}}
          onDelete={() => {}}
          onCancel={() => {}}
        />
      }
    >
      <DeleteBookInner />
    </Suspense>
  );
}
