"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useCallback, useEffect, useId, useState } from "react";
import { useAuth } from "@/app/_components/AuthProvider";
import { fetchClient } from "@/app/_lib/query";
import { parseIsbn } from "../../_models/isbn";
import { useScanner } from "./_components/useScanner";
import { NewBookPage, type RegisteredBook } from "./page.view";

type InputMode = "scanner" | "manual";

type BookDetailForm = {
  title: string;
  authors: Array<{ id: string; value: string }>;
  publisher: string;
  publishedDate: string;
  thumbnailUrl: string;
};

type State =
  | { type: "idle"; lastBook?: RegisteredBook }
  | { type: "registering" }
  | { type: "success"; book: RegisteredBook }
  | { type: "error"; message: string };

const validateBookForm = (form: BookDetailForm): string | null => {
  if (!form.title.trim()) {
    return "タイトルは必須です";
  }

  const validAuthors = form.authors.filter((a) => a.value.trim());
  if (validAuthors.length === 0) {
    return "著者は1人以上必須です";
  }

  if (form.publishedDate && !/^\d{4}-\d{2}-\d{2}$/.test(form.publishedDate)) {
    return "出版日の形式が正しくありません（YYYY-MM-DD）";
  }

  return null;
};

export default function NewBook() {
  const authState = useAuth();
  const queryClient = useQueryClient();
  const elementId = useId().replace(/:/g, "");

  const [inputMode, setInputMode] = useState<InputMode>("scanner");
  const [bookForm, setBookForm] = useState<BookDetailForm>({
    title: "",
    authors: [{ id: crypto.randomUUID(), value: "" }],
    publisher: "",
    publishedDate: "",
    thumbnailUrl: "",
  });
  const [state, setState] = useState<State>({ type: "idle" });

  const user = authState.status === "authenticated" ? authState.user : null;

  const handleBookRegistration = useCallback(
    async (rawCode: string) => {
      if (!user) {
        return;
      }
      if (state.type !== "idle") {
        return;
      }

      const isbn = parseIsbn(rawCode);
      if (!isbn) {
        setState({ type: "error", message: "ISBNの形式が正しくありません" });
        return;
      }

      setState({ type: "registering" });

      const token = await user.getIdToken();
      const { data, error: apiError } = await fetchClient.POST("/books/code", {
        body: { code: isbn },
        headers: { Authorization: `Bearer ${token}` },
      });

      if (apiError || !data) {
        setState({
          type: "error",
          message: apiError?.message ?? "書籍登録に失敗しました",
        });
      } else {
        queryClient.invalidateQueries({ queryKey: ["books"] });
        setState({ type: "idle", lastBook: data });
      }
    },
    [user, queryClient, state.type],
  );

  const { state: scannerState, scannedCode } = useScanner(
    elementId,
    inputMode === "scanner" && !!user,
  );

  const [lastProcessedCode, setLastProcessedCode] = useState<string | null>(
    null,
  );

  useEffect(() => {
    if (
      scannedCode &&
      scannedCode !== lastProcessedCode &&
      state.type === "idle"
    ) {
      setLastProcessedCode(scannedCode);
      handleBookRegistration(scannedCode);
    }
  }, [scannedCode, lastProcessedCode, state.type, handleBookRegistration]);

  useEffect(() => {
    if (!scannedCode) {
      setLastProcessedCode(null);
    }
  }, [scannedCode]);

  const handleRetry = () => {
    setState((prev) => ({
      type: "idle",
      lastBook: prev.type === "success" ? prev.book : undefined,
    }));
  };

  const handleManualBookRegistration = useCallback(
    async (form: BookDetailForm) => {
      if (!user) {
        return;
      }
      if (state.type !== "idle") {
        return;
      }

      const error = validateBookForm(form);
      if (error) {
        setState({ type: "error", message: error });
        return;
      }

      setState({ type: "registering" });

      const validAuthors = form.authors
        .filter((a) => a.value.trim())
        .map((a) => a.value.trim());

      const token = await user.getIdToken();
      const { data, error: apiError } = await fetchClient.POST("/books", {
        body: {
          title: form.title.trim(),
          authors: validAuthors,
          publisher: form.publisher.trim() || undefined,
          publishedDate: form.publishedDate.trim() || undefined,
          thumbnailUrl: form.thumbnailUrl.trim() || undefined,
        },
        headers: { Authorization: `Bearer ${token}` },
      });

      if (apiError || !data) {
        setState({
          type: "error",
          message: apiError?.message ?? "書籍登録に失敗しました",
        });
      } else {
        queryClient.invalidateQueries({ queryKey: ["books"] });
        setState({ type: "idle", lastBook: data });
        setBookForm({
          title: "",
          authors: [{ id: crypto.randomUUID(), value: "" }],
          publisher: "",
          publishedDate: "",
          thumbnailUrl: "",
        });
      }
    },
    [user, queryClient, state.type],
  );

  const handleManualSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    await handleManualBookRegistration(bookForm);
  };

  if (!user) {
    return (
      <NewBookPage
        registrationState={{ status: "idle" }}
        scannerState={{ status: "idle" }}
        inputMode="scanner"
        bookForm={{
          title: "",
          authors: [{ id: "temp", value: "" }],
          publisher: "",
          publishedDate: "",
          thumbnailUrl: "",
        }}
        scannerId=""
        onChangeInputMode={() => {}}
        onChangeBookForm={() => {}}
        onManualSubmit={() => {}}
        onRetry={() => {}}
      />
    );
  }

  const registrationState =
    state.type === "idle"
      ? { status: "idle" as const, lastBook: state.lastBook }
      : state.type === "registering"
        ? { status: "registering" as const }
        : state.type === "success"
          ? { status: "idle" as const, lastBook: state.book }
          : { status: "error" as const, message: state.message };

  return (
    <NewBookPage
      registrationState={registrationState}
      scannerState={scannerState}
      inputMode={inputMode}
      bookForm={bookForm}
      scannerId={elementId}
      onChangeInputMode={setInputMode}
      onChangeBookForm={setBookForm}
      onManualSubmit={handleManualSubmit}
      onRetry={handleRetry}
    />
  );
}
