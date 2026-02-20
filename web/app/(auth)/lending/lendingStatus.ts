import type { Book } from "../_models/book";

export type LendingStatus =
  | { kind: "available"; book: Book }
  | { kind: "borrowed_by_me"; book: Book }
  | { kind: "borrowed_by_other"; book: Book };

export function parseLendingStatus(
  book: Book,
  currentUserId: string,
): LendingStatus {
  if (book.status !== "borrowed") {
    return { kind: "available", book };
  }
  if (book.borrower?.id === currentUserId) {
    return { kind: "borrowed_by_me", book };
  }
  return { kind: "borrowed_by_other", book };
}
