"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useCallback, useId, useState } from "react";
import { useAuth } from "@/app/_components/AuthProvider";
import { fetchClient } from "@/app/_lib/query";
import { parseIsbn } from "../../_models/isbn";
import { useScanner } from "./_components/useScanner";
import { NewBookPage, type RegisteredBook } from "./page.view";

type InputMode = "scanner" | "manual";

type RegistrationState =
  | { status: "idle"; lastBook?: RegisteredBook }
  | { status: "registering"; lastBook?: RegisteredBook }
  | { status: "error"; message: string; lastBook?: RegisteredBook };

export default function NewBook() {
  const authState = useAuth();
  const queryClient = useQueryClient();
  const elementId = useId().replace(/:/g, "");

  const [code, setCode] = useState("");
  const [inputMode, setInputMode] = useState<InputMode>("scanner");
  const [registrationState, setRegistrationState] = useState<RegistrationState>(
    { status: "idle" },
  );

  const user = authState.status === "authenticated" ? authState.user : null;

  const registerBook = useCallback(
    async (isbn: string) => {
      if (!user) {
        return;
      }

      setRegistrationState((prev) => ({
        status: "registering",
        lastBook: prev.status !== "error" ? prev.lastBook : undefined,
      }));

      const token = await user.getIdToken();
      const { data, error: apiError } = await fetchClient.POST("/books/code", {
        body: { code: isbn },
        headers: { Authorization: `Bearer ${token}` },
      });

      if (apiError || !data) {
        setRegistrationState((prev) => ({
          status: "error",
          message: apiError?.message ?? "書籍登録に失敗しました",
          lastBook: prev.status !== "error" ? prev.lastBook : undefined,
        }));
      } else {
        queryClient.invalidateQueries({ queryKey: ["books"] });
        setRegistrationState({ status: "idle", lastBook: data });
      }
    },
    [user, queryClient],
  );

  const handleScan = useCallback(
    (decodedText: string) => {
      const isbn = parseIsbn(decodedText);
      if (!isbn) {
        setRegistrationState({
          status: "error",
          message: "読み取ったコードがISBNの形式ではありません",
        });
      } else {
        registerBook(isbn);
      }
    },
    [registerBook],
  );

  const scannerState = useScanner(
    elementId,
    inputMode === "scanner" && !!user,
    handleScan,
  );

  const handleSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    const isbn = parseIsbn(code);
    if (!isbn) {
      setRegistrationState({
        status: "error",
        message: "ISBNの形式が正しくありません",
      });
      return;
    }
    await registerBook(isbn);
  };

  const handleRetry = () => {
    setCode("");
    setRegistrationState({ status: "idle" });
  };

  if (!user) {
    return (
      <NewBookPage
        code=""
        registrationState={{ status: "idle" }}
        scannerState={{ status: "idle" }}
        inputMode="scanner"
        scannerId=""
        onChangeCode={() => {}}
        onChangeInputMode={() => {}}
        onSubmit={() => {}}
        onRetry={() => {}}
      />
    );
  }

  return (
    <NewBookPage
      code={code}
      registrationState={registrationState}
      scannerState={scannerState}
      inputMode={inputMode}
      scannerId={elementId}
      onChangeCode={setCode}
      onChangeInputMode={setInputMode}
      onSubmit={handleSubmit}
      onRetry={handleRetry}
    />
  );
}
