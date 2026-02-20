import { expect, test } from "vitest";
import type { Book } from "../_models/book";
import { parseLendingStatus } from "./lendingStatus";

test("when parseLendingStatus with available book then returns available", () => {
  const currentUserId = crypto.randomUUID();
  const book: Book = {
    id: crypto.randomUUID(),
    title: crypto.randomUUID(),
    authors: [crypto.randomUUID()],
    thumbnailUrl: "",
    status: "available",
    createdAt: new Date().toISOString(),
  };

  const result = parseLendingStatus(book, currentUserId);

  expect(result.kind).toBe("available");
  expect(result.book).toBe(book);
});

test("when parseLendingStatus with book borrowed by me then returns borrowed_by_me", () => {
  const currentUserId = crypto.randomUUID();
  const book: Book = {
    id: crypto.randomUUID(),
    title: crypto.randomUUID(),
    authors: [crypto.randomUUID()],
    thumbnailUrl: "",
    status: "borrowed",
    borrower: {
      id: currentUserId,
      name: crypto.randomUUID(),
      borrowedAt: new Date().toISOString(),
    },
    createdAt: new Date().toISOString(),
  };

  const result = parseLendingStatus(book, currentUserId);

  expect(result.kind).toBe("borrowed_by_me");
  expect(result.book).toBe(book);
});

test("when parseLendingStatus with book borrowed by other then returns borrowed_by_other", () => {
  const currentUserId = crypto.randomUUID();
  const book: Book = {
    id: crypto.randomUUID(),
    title: crypto.randomUUID(),
    authors: [crypto.randomUUID()],
    thumbnailUrl: "",
    status: "borrowed",
    borrower: {
      id: crypto.randomUUID(),
      name: crypto.randomUUID(),
      borrowedAt: new Date().toISOString(),
    },
    createdAt: new Date().toISOString(),
  };

  const result = parseLendingStatus(book, currentUserId);

  expect(result.kind).toBe("borrowed_by_other");
  expect(result.book).toBe(book);
});
