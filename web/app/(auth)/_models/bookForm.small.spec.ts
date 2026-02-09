import { expect, test } from "vitest";
import { parseBookForm } from "./bookForm";

const randomString = (length: number) =>
  Array.from({ length }, () =>
    String.fromCharCode(33 + Math.floor(Math.random() * 94)),
  ).join("");

const randomDate = () => {
  const year = 2000 + Math.floor(Math.random() * 25);
  const month = String(1 + Math.floor(Math.random() * 12)).padStart(2, "0");
  const day = String(1 + Math.floor(Math.random() * 28)).padStart(2, "0");
  return `${year}-${month}-${day}`;
};

const randomWhitespace = () => {
  const chars = [" ", "\t", "\n"];
  const length = 1 + Math.floor(Math.random() * 5);
  return Array.from(
    { length },
    () => chars[Math.floor(Math.random() * chars.length)],
  ).join("");
};

test("when parseBookForm with valid title and author then returns BookForm", () => {
  const title = randomString(10);
  const author = randomString(10);

  const result = parseBookForm({
    title,
    authors: [{ id: crypto.randomUUID(), value: author }],
    publisher: "",
    publishedDate: "",
    thumbnailUrl: "",
  });

  expect("error" in result).toBe(false);
  if ("error" in result) {
    return;
  }
  expect(result.title).toBe(title);
  expect(result.authors).toEqual([author]);
  expect(result.publisher).toBeUndefined();
  expect(result.publishedDate).toBeUndefined();
  expect(result.thumbnailUrl).toBeUndefined();
});

test("when parseBookForm with multiple authors then returns all authors", () => {
  const title = randomString(10);
  const authors = Array.from({ length: 3 }, () => randomString(10));

  const result = parseBookForm({
    title,
    authors: authors.map((a) => ({ id: crypto.randomUUID(), value: a })),
    publisher: "",
    publishedDate: "",
    thumbnailUrl: "",
  });

  expect("error" in result).toBe(false);
  if ("error" in result) {
    return;
  }
  expect(result.authors).toEqual(authors);
});

test("when parseBookForm with whitespace around title then trims title", () => {
  const coreTitle = randomString(10);
  const title = randomWhitespace() + coreTitle + randomWhitespace();
  const author = randomString(10);

  const result = parseBookForm({
    title,
    authors: [{ id: crypto.randomUUID(), value: author }],
    publisher: "",
    publishedDate: "",
    thumbnailUrl: "",
  });

  expect("error" in result).toBe(false);
  if ("error" in result) {
    return;
  }
  expect(result.title).toBe(coreTitle);
});

test("when parseBookForm with whitespace around author then trims author", () => {
  const title = randomString(10);
  const coreAuthor = randomString(10);
  const author = randomWhitespace() + coreAuthor + randomWhitespace();

  const result = parseBookForm({
    title,
    authors: [{ id: crypto.randomUUID(), value: author }],
    publisher: "",
    publishedDate: "",
    thumbnailUrl: "",
  });

  expect("error" in result).toBe(false);
  if ("error" in result) {
    return;
  }
  expect(result.authors).toEqual([coreAuthor]);
});

test("when parseBookForm with empty authors mixed then filters empty authors", () => {
  const title = randomString(10);
  const validAuthor = randomString(10);

  const result = parseBookForm({
    title,
    authors: [
      { id: crypto.randomUUID(), value: "" },
      { id: crypto.randomUUID(), value: validAuthor },
      { id: crypto.randomUUID(), value: randomWhitespace() },
    ],
    publisher: "",
    publishedDate: "",
    thumbnailUrl: "",
  });

  expect("error" in result).toBe(false);
  if ("error" in result) {
    return;
  }
  expect(result.authors).toEqual([validAuthor]);
});

test("when parseBookForm with all fields then preserves all fields", () => {
  const title = randomString(10);
  const author = randomString(10);
  const publisher = randomString(10);
  const publishedDate = randomDate();
  const thumbnailUrl = `https://example.com/${randomString(10)}.jpg`;

  const result = parseBookForm({
    title,
    authors: [{ id: crypto.randomUUID(), value: author }],
    publisher,
    publishedDate,
    thumbnailUrl,
  });

  expect("error" in result).toBe(false);
  if ("error" in result) {
    return;
  }
  expect(result.title).toBe(title);
  expect(result.authors).toEqual([author]);
  expect(result.publisher).toBe(publisher);
  expect(result.publishedDate).toBe(publishedDate);
  expect(result.thumbnailUrl).toBe(thumbnailUrl);
});

test("when parseBookForm with whitespace around optional fields then trims optional fields", () => {
  const title = randomString(10);
  const author = randomString(10);
  const corePublisher = randomString(10);
  const publisher = randomWhitespace() + corePublisher + randomWhitespace();
  const coreThumbnail = `https://example.com/${randomString(10)}.jpg`;
  const thumbnailUrl = randomWhitespace() + coreThumbnail + randomWhitespace();

  const result = parseBookForm({
    title,
    authors: [{ id: crypto.randomUUID(), value: author }],
    publisher,
    publishedDate: "",
    thumbnailUrl,
  });

  expect("error" in result).toBe(false);
  if ("error" in result) {
    return;
  }
  expect(result.publisher).toBe(corePublisher);
  expect(result.thumbnailUrl).toBe(coreThumbnail);
});

test.each([
  { title: "", expectedError: "タイトルは必須です" },
  { title: randomWhitespace(), expectedError: "タイトルは必須です" },
])(
  "when parseBookForm with invalid title '$title' then returns error",
  ({ title, expectedError }) => {
    const result = parseBookForm({
      title,
      authors: [{ id: crypto.randomUUID(), value: randomString(10) }],
      publisher: "",
      publishedDate: "",
      thumbnailUrl: "",
    });

    expect(result).toEqual({ error: expectedError });
  },
);

test("when parseBookForm with all empty authors then returns error", () => {
  const title = randomString(10);

  const result = parseBookForm({
    title,
    authors: [
      { id: crypto.randomUUID(), value: "" },
      { id: crypto.randomUUID(), value: randomWhitespace() },
    ],
    publisher: "",
    publishedDate: "",
    thumbnailUrl: "",
  });

  expect(result).toEqual({ error: "著者は1人以上必須です" });
});

test("when parseBookForm with empty authors array then returns error", () => {
  const result = parseBookForm({
    title: randomString(10),
    authors: [],
    publisher: "",
    publishedDate: "",
    thumbnailUrl: "",
  });

  expect(result).toEqual({ error: "著者は1人以上必須です" });
});

test.each([
  { publishedDate: "2021/03/27" },
  { publishedDate: "27-03-2021" },
  { publishedDate: "invalid" },
])(
  "when parseBookForm with invalid publishedDate '$publishedDate' then returns error",
  ({ publishedDate }) => {
    const result = parseBookForm({
      title: randomString(10),
      authors: [{ id: crypto.randomUUID(), value: randomString(10) }],
      publisher: "",
      publishedDate,
      thumbnailUrl: "",
    });

    expect(result).toEqual({
      error: "出版日の形式が正しくありません（YYYY-MM-DD）",
    });
  },
);
