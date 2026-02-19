"use client";

import { Icon } from "@iconify/react";
import type { FC } from "react";
import { BookInfoCard } from "../_components/BookInfoCard";
import { BorrowForm } from "../_components/BorrowForm";
import type { ScannerState } from "../books/new/_components/useScanner";
import type { LendingStatus } from "./lendingStatus";

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

interface LendingPageProps {
  scannerId: string;
  scannerState: ScannerState;
  scanState: ScanState;
  actionState: ActionState;
  dueDate: string;
  onBorrowClick: () => void;
  onBorrowConfirm: () => void;
  onBorrowCancel: () => void;
  onDueDateChange: (date: string) => void;
  onReturn: () => void;
  onReset: () => void;
}

export const LendingPage: FC<LendingPageProps> = ({
  scannerId,
  scannerState,
  scanState,
  actionState,
  dueDate,
  onBorrowClick,
  onBorrowConfirm,
  onBorrowCancel,
  onDueDateChange,
  onReturn,
  onReset,
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
            "px-4",
            "py-4",
          ].join(" ")}
        >
          <h1
            className={[
              "text-xl",
              "font-bold",
              "text-slate-900",
              "dark:text-slate-100",
            ].join(" ")}
          >
            貸出・返却
          </h1>
        </div>
      </header>

      <main className={["mx-auto", "max-w-md", "space-y-4", "p-4"].join(" ")}>
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
          <p
            className={[
              "mb-4",
              "text-sm",
              "font-medium",
              "text-slate-700",
              "dark:text-slate-300",
            ].join(" ")}
          >
            書籍のバーコードをスキャンしてください
          </p>
          <ScannerArea
            scannerId={scannerId}
            scannerState={scannerState}
            scanState={scanState}
          />
        </div>

        {scanState.status !== "idle" && (
          <ResultArea
            scanState={scanState}
            actionState={actionState}
            dueDate={dueDate}
            onBorrowClick={onBorrowClick}
            onBorrowConfirm={onBorrowConfirm}
            onBorrowCancel={onBorrowCancel}
            onDueDateChange={onDueDateChange}
            onReturn={onReturn}
            onReset={onReset}
          />
        )}
      </main>
    </div>
  );
};

const ScannerArea: FC<{
  scannerId: string;
  scannerState: ScannerState;
  scanState: ScanState;
}> = ({ scannerId, scannerState, scanState }) => {
  const isOverlay =
    scanState.status === "searching" ||
    scanState.status === "found" ||
    scannerState.status === "cooldown";

  if (scannerState.status === "error") {
    return (
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
        <Icon
          icon="material-symbols:camera-off"
          className={["mb-2", "size-8", "text-slate-400"].join(" ")}
        />
        <p
          className={["text-sm", "text-slate-600", "dark:text-slate-400"].join(
            " ",
          )}
        >
          {scannerState.message}
        </p>
      </div>
    );
  }

  return (
    <div className={["space-y-3"].join(" ")}>
      <div className={["relative"].join(" ")}>
        <div
          id={scannerId}
          className={[
            "overflow-hidden",
            "rounded-lg",
            "border",
            "border-slate-200",
            "dark:border-slate-700",
            isOverlay ? "opacity-50" : "",
          ].join(" ")}
        />
        {isOverlay && (
          <div
            className={[
              "absolute",
              "inset-0",
              "flex",
              "items-center",
              "justify-center",
              "rounded-lg",
              "bg-slate-900/50",
            ].join(" ")}
          >
            {scanState.status === "searching" && (
              <Icon
                icon="material-symbols:progress-activity"
                className={["size-8", "animate-spin", "text-white"].join(" ")}
              />
            )}
          </div>
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
          : scanState.status === "searching"
            ? "書籍を検索中..."
            : scanState.status === "found"
              ? "書籍が見つかりました"
              : "バーコードをカメラに映してください"}
      </p>
    </div>
  );
};

const ResultArea: FC<{
  scanState: ScanState;
  actionState: ActionState;
  dueDate: string;
  onBorrowClick: () => void;
  onBorrowConfirm: () => void;
  onBorrowCancel: () => void;
  onDueDateChange: (date: string) => void;
  onReturn: () => void;
  onReset: () => void;
}> = ({
  scanState,
  actionState,
  dueDate,
  onBorrowClick,
  onBorrowConfirm,
  onBorrowCancel,
  onDueDateChange,
  onReturn,
  onReset,
}) => {
  switch (scanState.status) {
    case "searching":
      return null;
    case "not_found":
      return (
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
          <p
            className={[
              "mb-4",
              "text-sm",
              "text-slate-600",
              "dark:text-slate-400",
            ].join(" ")}
          >
            この書籍は登録されていません。
          </p>
          <ResetButton onClick={onReset} />
        </div>
      );
    case "error":
      return (
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
          <div
            className={[
              "mb-4",
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
            {scanState.message}
          </div>
          <ResetButton onClick={onReset} />
        </div>
      );
    case "found": {
      const { lendingStatus } = scanState;
      const { book } = lendingStatus;
      const isLoading = actionState.status === "loading";
      return (
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
          <div className={["mb-6"].join(" ")}>
            <BookInfoCard
              title={book.title}
              authors={book.authors}
              status={book.status}
              borrower={book.borrower}
              thumbnailUrl={book.thumbnailUrl}
            />
          </div>

          {actionState.status === "error" && (
            <div
              className={[
                "mb-4",
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
              {actionState.message}
            </div>
          )}

          <div className={["space-y-2"].join(" ")}>
            <LendingActionButton
              lendingStatus={lendingStatus}
              isLoading={isLoading}
              actionState={actionState}
              dueDate={dueDate}
              onBorrowClick={onBorrowClick}
              onBorrowConfirm={onBorrowConfirm}
              onBorrowCancel={onBorrowCancel}
              onDueDateChange={onDueDateChange}
              onReturn={onReturn}
            />
            <ResetButton onClick={onReset} disabled={isLoading} />
          </div>
        </div>
      );
    }
  }
};

const LendingActionButton: FC<{
  lendingStatus: LendingStatus;
  isLoading: boolean;
  actionState: ActionState;
  dueDate: string;
  onBorrowClick: () => void;
  onBorrowConfirm: () => void;
  onBorrowCancel: () => void;
  onDueDateChange: (date: string) => void;
  onReturn: () => void;
}> = ({
  lendingStatus,
  isLoading,
  actionState,
  dueDate,
  onBorrowClick,
  onBorrowConfirm,
  onBorrowCancel,
  onDueDateChange,
  onReturn,
}) => {
  switch (lendingStatus.kind) {
    case "available":
      return (
        <BorrowForm
          isSelectingDate={actionState.status === "selecting-date"}
          isLoading={
            actionState.status === "loading" && actionState.action === "borrow"
          }
          dueDate={dueDate}
          onBorrowClick={onBorrowClick}
          onBorrowConfirm={onBorrowConfirm}
          onBorrowCancel={onBorrowCancel}
          onDueDateChange={onDueDateChange}
        />
      );
    case "borrowed_by_me":
      return (
        <button
          type="button"
          onClick={onReturn}
          disabled={isLoading}
          className={[
            "flex",
            "w-full",
            "items-center",
            "justify-center",
            "gap-2",
            "rounded-lg",
            "bg-emerald-700",
            "px-4",
            "py-3",
            "font-semibold",
            "text-white",
            "transition-colors",
            "hover:bg-emerald-800",
            "disabled:cursor-not-allowed",
            "disabled:bg-slate-300",
            "disabled:text-slate-500",
            "dark:bg-emerald-600",
            "dark:hover:bg-emerald-700",
            "dark:disabled:bg-slate-700",
            "dark:disabled:text-slate-400",
          ].join(" ")}
        >
          {actionState.status === "loading" &&
          actionState.action === "return" ? (
            <>
              <Icon
                icon="material-symbols:progress-activity"
                className="size-5 animate-spin"
              />
              返却中...
            </>
          ) : (
            <>
              <Icon
                icon="material-symbols:assignment-return"
                className="size-5"
              />
              返却する
            </>
          )}
        </button>
      );
    case "borrowed_by_other":
      return (
        <button
          type="button"
          onClick={onReturn}
          disabled={isLoading}
          className={[
            "flex",
            "w-full",
            "items-center",
            "justify-center",
            "gap-2",
            "rounded-lg",
            "border",
            "border-slate-300",
            "bg-white",
            "px-4",
            "py-3",
            "font-semibold",
            "text-slate-700",
            "transition-colors",
            "hover:bg-slate-50",
            "disabled:cursor-not-allowed",
            "disabled:bg-slate-100",
            "disabled:text-slate-400",
            "dark:border-slate-600",
            "dark:bg-slate-800",
            "dark:text-slate-300",
            "dark:hover:bg-slate-700",
            "dark:disabled:bg-slate-900",
            "dark:disabled:text-slate-600",
          ].join(" ")}
        >
          {actionState.status === "loading" &&
          actionState.action === "return" ? (
            <>
              <Icon
                icon="material-symbols:progress-activity"
                className="size-5 animate-spin"
              />
              返却中...
            </>
          ) : (
            <>
              <Icon
                icon="material-symbols:assignment-return"
                className="size-5"
              />
              代理返却する
            </>
          )}
        </button>
      );
  }
};

const ResetButton: FC<{ onClick: () => void; disabled?: boolean }> = ({
  onClick,
  disabled,
}) => (
  <button
    type="button"
    onClick={onClick}
    disabled={disabled}
    className={[
      "flex",
      "w-full",
      "items-center",
      "justify-center",
      "gap-1",
      "rounded-lg",
      "px-4",
      "py-2",
      "text-sm",
      "text-slate-600",
      "transition-colors",
      "hover:bg-slate-100",
      "disabled:cursor-not-allowed",
      "disabled:text-slate-400",
      "dark:text-slate-400",
      "dark:hover:bg-slate-800",
      "dark:disabled:text-slate-600",
    ].join(" ")}
  >
    <Icon icon="material-symbols:qr-code-scanner" className="size-4" />
    別の書籍をスキャン
  </button>
);
