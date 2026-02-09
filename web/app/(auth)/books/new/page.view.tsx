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
  registrationState: RegistrationState;
  scannerState: ScannerState;
  inputMode: InputMode;
  bookForm: {
    title: string;
    authors: Array<{ id: string; value: string }>;
    publisher: string;
    publishedDate: string;
    thumbnailUrl: string;
  };
  scannerId: string;
  onChangeInputMode: (mode: InputMode) => void;
  onChangeBookForm: (form: {
    title: string;
    authors: Array<{ id: string; value: string }>;
    publisher: string;
    publishedDate: string;
    thumbnailUrl: string;
  }) => void;
  onManualSubmit: (e: { preventDefault: () => void }) => void;
  onRetry: () => void;
}

export const NewBookPage: FC<NewBookPageProps> = ({
  registrationState,
  scannerState,
  inputMode,
  bookForm,
  scannerId,
  onChangeInputMode,
  onChangeBookForm,
  onManualSubmit,
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
            <InputModeView
              inputMode={inputMode}
              bookForm={bookForm}
              scannerId={scannerId}
              scannerState={scannerState}
              isRegistering={registrationState.status === "registering"}
              onChangeInputMode={onChangeInputMode}
              onChangeBookForm={onChangeBookForm}
              onManualSubmit={onManualSubmit}
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

interface InputModeViewProps {
  inputMode: InputMode;
  bookForm: {
    title: string;
    authors: Array<{ id: string; value: string }>;
    publisher: string;
    publishedDate: string;
    thumbnailUrl: string;
  };
  scannerId: string;
  scannerState: ScannerState;
  isRegistering: boolean;
  onChangeInputMode: (mode: InputMode) => void;
  onChangeBookForm: (form: {
    title: string;
    authors: Array<{ id: string; value: string }>;
    publisher: string;
    publishedDate: string;
    thumbnailUrl: string;
  }) => void;
  onManualSubmit: (e: { preventDefault: () => void }) => void;
}

const InputModeView: FC<InputModeViewProps> = ({
  inputMode,
  bookForm,
  scannerId,
  scannerState,
  isRegistering,
  onChangeInputMode,
  onChangeBookForm,
  onManualSubmit,
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
      <DetailInputForm
        bookForm={bookForm}
        isRegistering={isRegistering}
        onChangeBookForm={onChangeBookForm}
        onSubmit={onManualSubmit}
      />
    )}
  </>
);

interface DetailInputFormProps {
  bookForm: {
    title: string;
    authors: Array<{ id: string; value: string }>;
    publisher: string;
    publishedDate: string;
    thumbnailUrl: string;
  };
  isRegistering: boolean;
  onChangeBookForm: (form: {
    title: string;
    authors: Array<{ id: string; value: string }>;
    publisher: string;
    publishedDate: string;
    thumbnailUrl: string;
  }) => void;
  onSubmit: (e: { preventDefault: () => void }) => void;
}

const DetailInputForm: FC<DetailInputFormProps> = ({
  bookForm,
  isRegistering,
  onChangeBookForm,
  onSubmit,
}) => {
  const labelClasses = [
    "mb-2",
    "block",
    "text-sm",
    "font-medium",
    "text-slate-700",
    "dark:text-slate-300",
  ].join(" ");

  const inputClasses = [
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
  ].join(" ");

  const handleAuthorChange = (index: number, value: string) => {
    const newAuthors = [...bookForm.authors];
    newAuthors[index] = { ...newAuthors[index], value };
    onChangeBookForm({ ...bookForm, authors: newAuthors });
  };

  const handleAddAuthor = () => {
    onChangeBookForm({
      ...bookForm,
      authors: [...bookForm.authors, { id: crypto.randomUUID(), value: "" }],
    });
  };

  const handleRemoveAuthor = (index: number) => {
    if (bookForm.authors.length === 1) {
      return;
    }
    onChangeBookForm({
      ...bookForm,
      authors: bookForm.authors.filter((_, i) => i !== index),
    });
  };

  const isFormValid =
    !!bookForm.title.trim() && bookForm.authors.some((a) => !!a.value.trim());

  return (
    <form onSubmit={onSubmit} className={["space-y-4"].join(" ")}>
      <div>
        <label htmlFor="title" className={labelClasses}>
          タイトル <span className="text-red-500">*</span>
        </label>
        <input
          id="title"
          type="text"
          value={bookForm.title}
          onChange={(e) =>
            onChangeBookForm({ ...bookForm, title: e.target.value })
          }
          placeholder="詳解システム・パフォーマンス"
          disabled={isRegistering}
          className={inputClasses}
        />
      </div>

      <div>
        <div className={labelClasses}>
          著者 <span className="text-red-500">*</span>
        </div>
        <div className={["space-y-2"].join(" ")}>
          {bookForm.authors.map((author, index) => (
            <div key={author.id} className={["flex", "gap-2"].join(" ")}>
              <input
                type="text"
                value={author.value}
                onChange={(e) => handleAuthorChange(index, e.target.value)}
                placeholder={`著者${index + 1}`}
                aria-label={`著者${index + 1}`}
                disabled={isRegistering}
                className={inputClasses}
              />
              {bookForm.authors.length > 1 && (
                <button
                  type="button"
                  onClick={() => handleRemoveAuthor(index)}
                  aria-label="著者を削除"
                  disabled={isRegistering}
                  className={[
                    "flex",
                    "items-center",
                    "justify-center",
                    "rounded-lg",
                    "border",
                    "border-slate-200",
                    "bg-white",
                    "px-3",
                    "text-slate-600",
                    "transition-colors",
                    "hover:bg-slate-100",
                    "dark:border-slate-700",
                    "dark:bg-slate-800",
                    "dark:text-slate-400",
                    "dark:hover:bg-slate-700",
                    "disabled:opacity-50",
                  ].join(" ")}
                >
                  <Icon icon="material-symbols:close" className="size-5" />
                </button>
              )}
            </div>
          ))}
          <button
            type="button"
            onClick={handleAddAuthor}
            disabled={isRegistering}
            className={[
              "flex",
              "w-full",
              "items-center",
              "justify-center",
              "gap-1",
              "rounded-lg",
              "border",
              "border-dashed",
              "border-slate-300",
              "bg-slate-50",
              "px-4",
              "py-2",
              "text-sm",
              "text-slate-600",
              "transition-colors",
              "hover:bg-slate-100",
              "dark:border-slate-700",
              "dark:bg-slate-900",
              "dark:text-slate-400",
              "dark:hover:bg-slate-800",
              "disabled:opacity-50",
            ].join(" ")}
          >
            <Icon icon="material-symbols:add" className="size-4" />
            著者を追加
          </button>
        </div>
      </div>

      <div>
        <label htmlFor="publisher" className={labelClasses}>
          出版社
        </label>
        <input
          id="publisher"
          type="text"
          value={bookForm.publisher}
          onChange={(e) =>
            onChangeBookForm({ ...bookForm, publisher: e.target.value })
          }
          placeholder="オライリージャパン"
          disabled={isRegistering}
          className={inputClasses}
        />
      </div>

      <div>
        <label htmlFor="publishedDate" className={labelClasses}>
          出版日
        </label>
        <input
          id="publishedDate"
          type="date"
          value={bookForm.publishedDate}
          onChange={(e) =>
            onChangeBookForm({ ...bookForm, publishedDate: e.target.value })
          }
          disabled={isRegistering}
          className={inputClasses}
        />
      </div>

      <div>
        <label htmlFor="thumbnailUrl" className={labelClasses}>
          サムネイルURL
        </label>
        <input
          id="thumbnailUrl"
          type="url"
          value={bookForm.thumbnailUrl}
          onChange={(e) =>
            onChangeBookForm({ ...bookForm, thumbnailUrl: e.target.value })
          }
          placeholder="https://example.com/book-cover.jpg"
          disabled={isRegistering}
          className={inputClasses}
        />
      </div>

      <button
        type="submit"
        disabled={!isFormValid || isRegistering}
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
  );
};
