"use client";

import { Icon } from "@iconify/react";
import type { FC } from "react";
import { BookInfoCard } from "../../_components/BookInfoCard";
import type { Book } from "../../_models/book";

type DeleteReason = "transfer" | "disposal" | "lost" | "other";

type DeleteState =
  | { status: "idle" }
  | { status: "deleting" }
  | { status: "error"; message: string };

interface DeleteBookPageProps {
  book: Book;
  deleteState: DeleteState;
  deleteReason: DeleteReason;
  deleteMemo: string;
  onDeleteReasonChange: (reason: DeleteReason) => void;
  onDeleteMemoChange: (memo: string) => void;
  onDelete: () => void;
  onCancel: () => void;
}

export const DeleteBookPage: FC<DeleteBookPageProps> = ({
  book,
  deleteState,
  deleteReason,
  deleteMemo,
  onDeleteReasonChange,
  onDeleteMemoChange,
  onDelete,
  onCancel,
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
          <button
            type="button"
            onClick={onCancel}
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
          </button>
          <h1
            className={[
              "text-xl",
              "font-bold",
              "text-slate-900",
              "dark:text-slate-100",
            ].join(" ")}
          >
            書籍を削除
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
          <BookInfoCard
            title={book.title}
            authors={book.authors}
            status={book.status}
            borrower={book.borrower}
            thumbnailUrl={book.thumbnailUrl}
          />
        </div>

        {deleteState.status === "error" && (
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
            {deleteState.message}
          </div>
        )}

        {book.status === "borrowed" && (
          <div
            className={[
              "rounded-lg",
              "bg-amber-50",
              "px-4",
              "py-3",
              "text-sm",
              "text-amber-800",
              "dark:bg-amber-500/10",
              "dark:text-amber-400",
            ].join(" ")}
          >
            この書籍は現在貸出中のため、削除できません。
          </div>
        )}

        <div
          className={[
            "rounded-xl",
            "border",
            "border-slate-200",
            "bg-white",
            "p-6",
            "dark:border-slate-800",
            "dark:bg-slate-900",
            "space-y-4",
          ].join(" ")}
        >
          <div
            className={[
              "flex",
              "items-start",
              "gap-3",
              "rounded-lg",
              "bg-red-50",
              "p-4",
              "dark:bg-red-950",
            ].join(" ")}
          >
            <Icon
              icon="material-symbols:warning"
              className={[
                "h-5",
                "w-5",
                "text-red-600",
                "dark:text-red-400",
              ].join(" ")}
            />
            <div>
              <p
                className={[
                  "text-sm",
                  "font-semibold",
                  "text-red-800",
                  "dark:text-red-300",
                ].join(" ")}
              >
                この操作は取り消せません
              </p>
              <p
                className={[
                  "mt-1",
                  "text-xs",
                  "text-red-700",
                  "dark:text-red-400",
                ].join(" ")}
              >
                削除された書籍の情報は復元できません。
              </p>
            </div>
          </div>

          <div>
            <label
              htmlFor="deleteReason"
              className={[
                "mb-2",
                "block",
                "text-sm",
                "font-medium",
                "text-slate-700",
                "dark:text-slate-300",
              ].join(" ")}
            >
              削除理由
            </label>
            <select
              id="deleteReason"
              value={deleteReason}
              onChange={(e) =>
                onDeleteReasonChange(e.target.value as DeleteReason)
              }
              disabled={deleteState.status === "deleting"}
              className={[
                "w-full",
                "rounded-lg",
                "border",
                "border-slate-200",
                "bg-white",
                "px-4",
                "py-3",
                "text-slate-900",
                "focus:border-sky-600",
                "focus:outline-none",
                "focus:ring-1",
                "focus:ring-sky-600",
                "disabled:cursor-not-allowed",
                "disabled:bg-slate-100",
                "disabled:text-slate-500",
                "dark:border-slate-700",
                "dark:bg-slate-800",
                "dark:text-slate-100",
                "dark:focus:border-sky-500",
                "dark:focus:ring-sky-500",
                "dark:disabled:bg-slate-900",
              ].join(" ")}
            >
              <option value="transfer">譲渡（他の人に譲った）</option>
              <option value="disposal">破棄（廃棄処分した）</option>
              <option value="lost">紛失（紛失した）</option>
              <option value="other">その他</option>
            </select>
          </div>

          <div>
            <label
              htmlFor="deleteMemo"
              className={[
                "mb-2",
                "block",
                "text-sm",
                "font-medium",
                "text-slate-700",
                "dark:text-slate-300",
              ].join(" ")}
            >
              メモ（任意）
            </label>
            <textarea
              id="deleteMemo"
              value={deleteMemo}
              onChange={(e) => onDeleteMemoChange(e.target.value)}
              placeholder="削除の詳細な理由や備考を入力できます"
              maxLength={500}
              rows={4}
              disabled={deleteState.status === "deleting"}
              className={[
                "w-full",
                "rounded-lg",
                "border",
                "border-slate-200",
                "bg-white",
                "px-4",
                "py-3",
                "text-slate-900",
                "focus:border-sky-600",
                "focus:outline-none",
                "focus:ring-1",
                "focus:ring-sky-600",
                "disabled:cursor-not-allowed",
                "disabled:bg-slate-100",
                "disabled:text-slate-500",
                "dark:border-slate-700",
                "dark:bg-slate-800",
                "dark:text-slate-100",
                "dark:focus:border-sky-500",
                "dark:focus:ring-sky-500",
                "dark:disabled:bg-slate-900",
              ].join(" ")}
            />
            <p
              className={[
                "mt-1",
                "text-xs",
                "text-slate-500",
                "dark:text-slate-400",
              ].join(" ")}
            >
              {deleteMemo.length}/500
            </p>
          </div>

          <button
            type="button"
            onClick={onDelete}
            disabled={
              deleteState.status === "deleting" || book.status === "borrowed"
            }
            className={[
              "flex",
              "w-full",
              "items-center",
              "justify-center",
              "gap-1",
              "rounded-lg",
              "bg-red-700",
              "px-4",
              "py-3",
              "font-semibold",
              "text-white",
              "transition-colors",
              "hover:bg-red-800",
              "disabled:cursor-not-allowed",
              "disabled:bg-slate-300",
              "disabled:text-slate-500",
              "dark:bg-red-600",
              "dark:hover:bg-red-700",
              "dark:disabled:bg-slate-700",
              "dark:disabled:text-slate-400",
            ].join(" ")}
          >
            {deleteState.status === "deleting" ? (
              <>
                <Icon
                  icon="material-symbols:progress-activity"
                  className="h-5 w-5 animate-spin"
                />
                削除中...
              </>
            ) : (
              <>
                <Icon icon="material-symbols:delete" className="h-5 w-5" />
                削除する
              </>
            )}
          </button>
        </div>
      </main>
    </div>
  );
};
