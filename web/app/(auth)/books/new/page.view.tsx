"use client";

import { Icon } from "@iconify/react";
import Link from "next/link";
import type { FC } from "react";
import { RegisteredBookCard } from "./_components/RegisteredBookCard";
import type { ScannerState } from "./_components/useScanner";

export type RegisteredBook = {
  id: string;
  title: string;
  authors: string[];
  thumbnailUrl?: string;
};

type InputMode = "scanner" | "manual";

type RegistrationState =
  | { status: "idle"; lastBook?: RegisteredBook }
  | { status: "registering"; lastBook?: RegisteredBook }
  | { status: "error"; message: string; lastBook?: RegisteredBook };

interface NewBookPageProps {
  code: string;
  registrationState: RegistrationState;
  scannerState: ScannerState;
  inputMode: InputMode;
  scannerId: string;
  onChangeCode: (code: string) => void;
  onChangeInputMode: (mode: InputMode) => void;
  onSubmit: (e: { preventDefault: () => void }) => void;
  onRetry: () => void;
}

export const NewBookPage: FC<NewBookPageProps> = ({
  code,
  registrationState,
  scannerState,
  inputMode,
  scannerId,
  onChangeCode,
  onChangeInputMode,
  onSubmit,
  onRetry,
}) => {
  return (
    <div
      className={["min-h-screen", "bg-slate-50", "dark:bg-slate-950"].join(" ")}
    >
      <header
        className={[
          "border-b",
          "border-slate-200",
          "bg-white",
          "dark:border-slate-800",
          "dark:bg-slate-900",
        ].join(" ")}
      >
        <div
          className={[
            "mx-auto",
            "flex",
            "max-w-4xl",
            "items-center",
            "gap-4",
            "px-4",
            "py-4",
          ].join(" ")}
        >
          <Link
            href="/"
            className={[
              "flex",
              "items-center",
              "gap-1",
              "text-slate-600",
              "transition-colors",
              "hover:text-slate-900",
              "dark:text-slate-400",
              "dark:hover:text-slate-100",
            ].join(" ")}
          >
            <Icon
              icon="material-symbols:arrow-back"
              className={["h-4", "w-4"].join(" ")}
            />
            戻る
          </Link>
          <h1
            className={[
              "text-xl",
              "font-bold",
              "text-slate-900",
              "dark:text-slate-100",
            ].join(" ")}
          >
            書籍登録
          </h1>
        </div>
      </header>

      <main className={["mx-auto", "max-w-md", "p-4", "space-y-4"].join(" ")}>
        <div
          className={[
            "rounded-xl",
            "border",
            "border-slate-200",
            "bg-white",
            "p-6",
            "dark:border-slate-800",
            "dark:bg-slate-900",
          ].join(" ")}
        >
          {registrationState.status === "error" ? (
            <ErrorView message={registrationState.message} onRetry={onRetry} />
          ) : (
            <ScanningView
              code={code}
              inputMode={inputMode}
              scannerId={scannerId}
              scannerState={scannerState}
              isRegistering={registrationState.status === "registering"}
              onChangeCode={onChangeCode}
              onChangeInputMode={onChangeInputMode}
              onSubmit={onSubmit}
            />
          )}
        </div>

        {registrationState.lastBook && (
          <RegisteredBookCard book={registrationState.lastBook} />
        )}
      </main>
    </div>
  );
};

const ErrorView: FC<{
  message: string;
  onRetry: () => void;
}> = ({ message, onRetry }) => (
  <div className={["space-y-4"].join(" ")}>
    <div
      className={[
        "rounded-lg",
        "bg-red-50",
        "px-4",
        "py-3",
        "text-sm",
        "text-red-600",
        "dark:bg-red-500/10",
        "dark:text-red-400",
      ].join(" ")}
    >
      {message}
    </div>

    <button
      type="button"
      onClick={onRetry}
      className={[
        "flex",
        "w-full",
        "items-center",
        "justify-center",
        "gap-1",
        "rounded-lg",
        "bg-sky-700",
        "px-4",
        "py-3",
        "font-semibold",
        "text-white",
        "transition-colors",
        "hover:bg-sky-800",
        "dark:bg-sky-600",
        "dark:hover:bg-sky-700",
      ].join(" ")}
    >
      もう一度試す
    </button>
  </div>
);

interface ScanningViewProps {
  code: string;
  inputMode: InputMode;
  scannerId: string;
  scannerState: ScannerState;
  isRegistering: boolean;
  onChangeCode: (code: string) => void;
  onChangeInputMode: (mode: InputMode) => void;
  onSubmit: (e: { preventDefault: () => void }) => void;
}

const ScanningView: FC<ScanningViewProps> = ({
  code,
  inputMode,
  scannerId,
  scannerState,
  isRegistering,
  onChangeCode,
  onChangeInputMode,
  onSubmit,
}) => (
  <>
    <div className={["mb-4", "flex", "gap-2"].join(" ")}>
      <button
        type="button"
        onClick={() => onChangeInputMode("scanner")}
        disabled={isRegistering}
        className={[
          "flex",
          "flex-1",
          "items-center",
          "justify-center",
          "gap-1",
          "rounded-lg",
          "px-3",
          "py-2",
          "text-sm",
          "font-medium",
          "transition-colors",
          inputMode === "scanner"
            ? "bg-sky-700 text-white dark:bg-sky-600"
            : "bg-slate-100 text-slate-600 hover:bg-slate-200 dark:bg-slate-800 dark:text-slate-400 dark:hover:bg-slate-700",
        ].join(" ")}
      >
        <Icon
          icon="material-symbols:qr-code-scanner"
          className={["size-4"].join(" ")}
        />
        スキャン
      </button>
      <button
        type="button"
        onClick={() => onChangeInputMode("manual")}
        disabled={isRegistering}
        className={[
          "flex",
          "flex-1",
          "items-center",
          "justify-center",
          "gap-1",
          "rounded-lg",
          "px-3",
          "py-2",
          "text-sm",
          "font-medium",
          "transition-colors",
          inputMode === "manual"
            ? "bg-sky-700 text-white dark:bg-sky-600"
            : "bg-slate-100 text-slate-600 hover:bg-slate-200 dark:bg-slate-800 dark:text-slate-400 dark:hover:bg-slate-700",
        ].join(" ")}
      >
        <Icon
          icon="material-symbols:keyboard"
          className={["size-4"].join(" ")}
        />
        手入力
      </button>
    </div>

    {inputMode === "scanner" ? (
      <div className={["space-y-4"].join(" ")}>
        {scannerState.status === "error" ? (
          <div
            className={[
              "flex",
              "flex-col",
              "items-center",
              "justify-center",
              "rounded-lg",
              "border",
              "border-slate-200",
              "bg-slate-100",
              "p-8",
              "text-center",
              "dark:border-slate-700",
              "dark:bg-slate-800",
            ].join(" ")}
          >
            <p
              className={[
                "text-sm",
                "text-slate-600",
                "dark:text-slate-400",
              ].join(" ")}
            >
              {scannerState.message}
            </p>
          </div>
        ) : (
          <>
            <div className={["relative"].join(" ")}>
              <div
                id={scannerId}
                className={[
                  "overflow-hidden",
                  "rounded-lg",
                  "border",
                  "border-slate-200",
                  "dark:border-slate-700",
                  scannerState.status === "cooldown" ? "opacity-50" : "",
                ].join(" ")}
              />
              {scannerState.status === "cooldown" && (
                <div
                  className={[
                    "absolute",
                    "inset-0",
                    "rounded-lg",
                    "bg-slate-900/50",
                    "dark:bg-slate-950/70",
                  ].join(" ")}
                />
              )}
            </div>
            <p
              className={[
                "text-center",
                "text-sm",
                "text-slate-500",
                "dark:text-slate-400",
              ].join(" ")}
            >
              {scannerState.status === "idle" ||
              scannerState.status === "initializing"
                ? "カメラを起動中..."
                : isRegistering
                  ? "登録中..."
                  : scannerState.status === "cooldown"
                    ? "少々お待ちください..."
                    : "書籍のバーコードをカメラに映してください"}
            </p>
          </>
        )}
      </div>
    ) : (
      <form onSubmit={onSubmit} className={["space-y-4"].join(" ")}>
        <div>
          <label
            htmlFor="code"
            className={[
              "mb-2",
              "block",
              "text-sm",
              "font-medium",
              "text-slate-700",
              "dark:text-slate-300",
            ].join(" ")}
          >
            ISBN
          </label>
          <input
            id="code"
            type="text"
            value={code}
            onChange={(e) => onChangeCode(e.target.value)}
            placeholder="978-4-87311-752-2"
            disabled={isRegistering}
            className={[
              "w-full",
              "rounded-lg",
              "border",
              "border-slate-200",
              "bg-white",
              "px-4",
              "py-3",
              "text-slate-900",
              "placeholder:text-slate-400",
              "focus:border-sky-600",
              "focus:outline-none",
              "focus:ring-1",
              "focus:ring-sky-600",
              "dark:border-slate-700",
              "dark:bg-slate-800",
              "dark:text-slate-100",
              "dark:placeholder:text-slate-500",
              "dark:focus:border-sky-500",
              "dark:focus:ring-sky-500",
              "disabled:opacity-50",
            ].join(" ")}
          />
        </div>

        <button
          type="submit"
          disabled={!code.trim() || isRegistering}
          className={[
            "flex",
            "w-full",
            "items-center",
            "justify-center",
            "gap-1",
            "rounded-lg",
            "bg-sky-700",
            "px-4",
            "py-3",
            "font-semibold",
            "text-white",
            "transition-colors",
            "hover:bg-sky-800",
            "disabled:cursor-not-allowed",
            "disabled:bg-slate-300",
            "disabled:text-slate-500",
            "dark:bg-sky-600",
            "dark:hover:bg-sky-700",
            "dark:disabled:bg-slate-700",
            "dark:disabled:text-slate-400",
          ].join(" ")}
        >
          {isRegistering ? (
            <>
              <Icon
                icon="material-symbols:progress-activity"
                className="h-5 w-5 animate-spin"
              />
              登録中...
            </>
          ) : (
            <>
              <Icon icon="material-symbols:add" className="h-5 w-5" />
              登録する
            </>
          )}
        </button>
      </form>
    )}
  </>
);
