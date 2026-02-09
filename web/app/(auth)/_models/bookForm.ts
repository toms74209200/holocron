export type BookForm = {
  title: string;
  authors: string[];
  publisher?: string;
  publishedDate?: string;
  thumbnailUrl?: string;
};

export const parseBookForm = (form: {
  title: string;
  authors: Array<{ id: string; value: string }>;
  publisher: string;
  publishedDate: string;
  thumbnailUrl: string;
}): BookForm | { error: string } => {
  if (!form.title.trim()) {
    return { error: "タイトルは必須です" };
  }

  const validAuthors = form.authors
    .filter((a) => a.value.trim())
    .map((a) => a.value.trim());

  if (validAuthors.length === 0) {
    return { error: "著者は1人以上必須です" };
  }

  if (form.publishedDate && !/^\d{4}-\d{2}-\d{2}$/.test(form.publishedDate)) {
    return { error: "出版日の形式が正しくありません（YYYY-MM-DD）" };
  }

  return {
    title: form.title.trim(),
    authors: validAuthors,
    publisher: form.publisher.trim() || undefined,
    publishedDate: form.publishedDate.trim() || undefined,
    thumbnailUrl: form.thumbnailUrl.trim() || undefined,
  };
};
