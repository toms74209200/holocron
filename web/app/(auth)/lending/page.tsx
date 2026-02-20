"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useCallback, useEffect, useId, useState } from "react";
import { Temporal } from "temporal-polyfill";
import { useAuth } from "@/app/_components/AuthProvider";
import { fetchClient } from "@/app/_lib/query";
import type { Book } from "../_models/book";
import { useScanner } from "../books/new/_components/useScanner";
import type { LendingStatus } from "./lendingStatus";
import { parseLendingStatus } from "./lendingStatus";
import { LendingPage } from "./page.view";

type ScanState =
  | { status: "idle" }
  | { status: "searching" }
  | { status: "found"; lendingStatus: LendingStatus }
  | { status: "not_found" }
  | { status: "error"; message: string };

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

export default function Lending() {
  const authState = useAuth();
  const queryClient = useQueryClient();
  const elementId = useId().replace(/:/g, "");

  const [scanState, setScanState] = useState<ScanState>({ status: "idle" });
  const [actionState, setActionState] = useState<ActionState>({
    status: "idle",
  });
  const [dueDate, setDueDate] = useState<string>(getDefaultDueDate());

  const user = authState.status === "authenticated" ? authState.user : null;

  const scannerEnabled =
    !!user && scanState.status === "idle" && actionState.status === "idle";

  const { state: scannerState, scannedCode } = useScanner(
    elementId,
    scannerEnabled,
  );

  const handleCodeScanned = useCallback(
    async (code: string) => {
      if (!user) {
        return;
      }

      setScanState({ status: "searching" });

      const token = await user.getIdToken();
      const { data, error } = await fetchClient.GET("/books", {
        params: { query: { code } },
        headers: { Authorization: `Bearer ${token}` },
      });

      if (error || !data) {
        setScanState({ status: "error", message: "書籍の検索に失敗しました" });
        return;
      }

      if (data.items.length === 0) {
        setScanState({ status: "not_found" });
        return;
      }

      setScanState({
        status: "found",
        lendingStatus: parseLendingStatus(data.items[0] as Book, user.uid),
      });
    },
    [user],
  );

  const [lastProcessedCode, setLastProcessedCode] = useState<string | null>(
    null,
  );

  useEffect(() => {
    if (
      scannedCode &&
      scannedCode !== lastProcessedCode &&
      scanState.status === "idle"
    ) {
      setLastProcessedCode(scannedCode);
      handleCodeScanned(scannedCode);
    }
  }, [scannedCode, lastProcessedCode, scanState.status, handleCodeScanned]);

  useEffect(() => {
    if (!scannedCode) {
      setLastProcessedCode(null);
    }
  }, [scannedCode]);

  const handleBorrowClick = useCallback(() => {
    setActionState({ status: "selecting-date" });
  }, []);

  const handleBorrowConfirm = useCallback(async () => {
    if (!user || scanState.status !== "found") {
      return;
    }

    const bookId = scanState.lendingStatus.book.id;
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
    } catch {
      setActionState({ status: "error", message: "通信エラーが発生しました" });
      return;
    }

    await queryClient.invalidateQueries({ queryKey: ["books"] });
    setActionState({ status: "idle" });
  }, [user, scanState, dueDate, queryClient]);

  const handleBorrowCancel = useCallback(() => {
    setActionState({ status: "idle" });
  }, []);

  const handleReturn = useCallback(async () => {
    if (!user || scanState.status !== "found") {
      return;
    }

    const bookId = scanState.lendingStatus.book.id;
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
    } catch {
      setActionState({ status: "error", message: "通信エラーが発生しました" });
      return;
    }

    await queryClient.invalidateQueries({ queryKey: ["books"] });
    setActionState({ status: "idle" });
  }, [user, scanState, queryClient]);

  const handleReset = useCallback(() => {
    setScanState({ status: "idle" });
    setActionState({ status: "idle" });
  }, []);

  return (
    <LendingPage
      scannerId={elementId}
      scannerState={scannerState}
      scanState={scanState}
      actionState={actionState}
      dueDate={dueDate}
      onBorrowClick={handleBorrowClick}
      onBorrowConfirm={handleBorrowConfirm}
      onBorrowCancel={handleBorrowCancel}
      onDueDateChange={setDueDate}
      onReturn={handleReturn}
      onReset={handleReset}
    />
  );
}
