"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useCallback, useEffect, useId, useState } from "react";
import { useAuth } from "@/app/_components/AuthProvider";
import { fetchClient } from "@/app/_lib/query";
import { parseIsbn } from "../../_models/isbn";
import { useScanner } from "./_components/useScanner";
import { NewBookPage, type RegisteredBook } from "./page.view";

type InputMode = "scanner" | "manual";

type State =
  | { type: "idle"; lastBook?: RegisteredBook }
  | { type: "registering" }
  | { type: "success"; book: RegisteredBook }
  | { type: "error"; message: string };

export default function NewBook() {
  const authState = useAuth();
  const queryClient = useQueryClient();
  const elementId = useId().replace(/:/g, "");

  const [code, setCode] = useState("");
  const [inputMode, setInputMode] = useState<InputMode>("scanner");
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

  const handleSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    await handleBookRegistration(code);
  };

  const handleRetry = () => {
    setCode("");
    setState((prev) => ({
      type: "idle",
      lastBook: prev.type === "success" ? prev.book : undefined,
    }));
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
