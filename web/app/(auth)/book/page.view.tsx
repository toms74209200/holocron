"use client";

import { Icon } from "@iconify/react";
import Link from "next/link";
import type { FC } from "react";
import { BookInfoCard } from "../_components/BookInfoCard";
import { BorrowForm } from "../_components/BorrowForm";
import type { Book } from "../_models/book";

type ActionState =
  | { status: "idle" }
  | { status: "selecting-date" }
  | { status: "loading"; action: "borrow" | "return" }
  | { status: "error"; message: string };

interface BookDetailPageProps {
  book: Book;
  actionState: ActionState;
  currentUserId: string;
  dueDate: string;
  onBorrowClick: () => void;
  onBorrowConfirm: () => void;
  onBorrowCancel: () => void;
  onDueDateChange: (date: string) => void;
  onReturn: () => void;
}

export const BookDetailPage: FC<BookDetailPageProps> = ({
  book,
  actionState,
  currentUserId,
  dueDate,
  onBorrowClick,
  onBorrowConfirm,
  onBorrowCancel,
  onDueDateChange,
  onReturn,
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
            書籍詳細
          </h1>
        </div>
      </header>

      <main className={["mx-auto", "max-w-md", "p-4", "space-y-4"].join(" ")}>
        <div
          className={[
            "relative",
            "rounded-xl",
            "border",
            "border-slate-200",
            "bg-white",
            "p-6",
            "dark:border-slate-800",
            "dark:bg-slate-900",
          ].join(" ")}
        >
          {book.status === "available" && (
            <Link
              href={`/book/delete?id=${book.id}`}
              className={[
                "absolute",
                "right-4",
                "top-4",
                "flex",
                "items-center",
                "justify-center",
                "rounded-lg",
                "p-2",
                "text-red-600",
                "transition-colors",
                "hover:bg-red-50",
                "dark:text-red-400",
                "dark:hover:bg-red-950",
              ].join(" ")}
              aria-label="削除"
            >
              <Icon
                icon="material-symbols:delete-outline"
                className="h-5 w-5"
              />
            </Link>
          )}

          <div className={book.status === "available" ? "pr-10" : ""}>
            <BookInfoCard
              title={book.title}
              authors={book.authors}
              status={book.status}
              borrower={book.borrower}
              thumbnailUrl={book.thumbnailUrl}
            />
          </div>

          {(book.publisher || book.publishedDate || book.code) && (
            <div
              className={[
                "mt-4",
                "pt-4",
                "border-t",
                "border-slate-200",
                "dark:border-slate-800",
                "space-y-2",
              ].join(" ")}
            >
              {book.publisher && (
                <div>
                  <p
                    className={[
                      "text-xs",
                      "text-slate-500",
                      "dark:text-slate-400",
                    ].join(" ")}
                  >
                    出版社
                  </p>
                  <p
                    className={[
                      "text-sm",
                      "text-slate-900",
                      "dark:text-slate-100",
                    ].join(" ")}
                  >
                    {book.publisher}
                  </p>
                </div>
              )}
              {book.publishedDate && (
                <div>
                  <p
                    className={[
                      "text-xs",
                      "text-slate-500",
                      "dark:text-slate-400",
                    ].join(" ")}
                  >
                    出版日
                  </p>
                  <p
                    className={[
                      "text-sm",
                      "text-slate-900",
                      "dark:text-slate-100",
                    ].join(" ")}
                  >
                    {book.publishedDate}
                  </p>
                </div>
              )}
              {book.code && (
                <div>
                  <p
                    className={[
                      "text-xs",
                      "text-slate-500",
                      "dark:text-slate-400",
                    ].join(" ")}
                  >
                    コード
                  </p>
                  <p
                    className={[
                      "text-sm",
                      "text-slate-900",
                      "dark:text-slate-100",
                    ].join(" ")}
                  >
                    {book.code}
                  </p>
                </div>
              )}
            </div>
          )}
        </div>

        {actionState.status === "error" && (
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
            {actionState.message}
          </div>
        )}

        {book.status === "available" && (
          <BorrowForm
            isSelectingDate={actionState.status === "selecting-date"}
            isLoading={
              actionState.status === "loading" &&
              actionState.action === "borrow"
            }
            dueDate={dueDate}
            onBorrowClick={onBorrowClick}
            onBorrowConfirm={onBorrowConfirm}
            onBorrowCancel={onBorrowCancel}
            onDueDateChange={onDueDateChange}
          />
        )}

        {book.status === "borrowed" && book.borrower?.id === currentUserId && (
          <button
            type="button"
            onClick={onReturn}
            disabled={actionState.status === "loading"}
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
            {actionState.status === "loading" &&
            actionState.action === "return" ? (
              <>
                <Icon
                  icon="material-symbols:progress-activity"
                  className="h-5 w-5 animate-spin"
                />
                返却中...
              </>
            ) : (
              <>
                <Icon
                  icon="material-symbols:keyboard-return"
                  className="h-5 w-5"
                />
                返却する
              </>
            )}
          </button>
        )}
      </main>
    </div>
  );
};
