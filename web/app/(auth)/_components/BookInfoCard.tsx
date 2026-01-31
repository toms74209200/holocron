// biome-ignore-all lint/performance/noImgElement: Image content can be external

import type { FC } from "react";

export interface BookInfoCardProps {
  title: string;
  authors: string[];
  status: "available" | "borrowed";
  borrower?: {
    id: string;
    name: string;
    borrowedAt: string;
  };
  thumbnailUrl: string;
}

export const BookInfoCard: FC<BookInfoCardProps> = ({
  title,
  authors,
  status,
  borrower,
  thumbnailUrl,
}) => {
  return (
    <div className={["flex", "gap-4"].join(" ")}>
      <div
        className={[
          "h-28",
          "w-20",
          "shrink-0",
          "overflow-hidden",
          "rounded-lg",
          "bg-slate-100",
          "dark:bg-slate-800",
        ].join(" ")}
      >
        <img
          src={thumbnailUrl}
          alt={title}
          className={["h-full", "w-full", "object-cover"].join(" ")}
        />
      </div>

      <div
        className={["flex", "flex-1", "flex-col", "justify-between"].join(" ")}
      >
        <div>
          <h2
            className={[
              "font-bold",
              "text-slate-900",
              "line-clamp-2",
              "dark:text-slate-100",
            ].join(" ")}
          >
            {title}
          </h2>
          <p
            className={[
              "mt-1",
              "text-sm",
              "text-slate-500",
              "dark:text-slate-400",
            ].join(" ")}
          >
            {authors.join(", ")}
          </p>
        </div>

        <div className={["mt-2"].join(" ")}>
          {status === "available" ? (
            <span
              className={[
                "text-sm",
                "text-emerald-600",
                "dark:text-emerald-400",
              ].join(" ")}
            >
              貸出可能
            </span>
          ) : (
            <span
              className={[
                "text-sm",
                "text-amber-600",
                "dark:text-amber-400",
              ].join(" ")}
            >
              貸出中: {borrower?.name}
            </span>
          )}
        </div>
      </div>
    </div>
  );
};
