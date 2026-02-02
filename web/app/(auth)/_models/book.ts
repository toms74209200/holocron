export type Book = {
  id: string;
  code?: string;
  title: string;
  authors: string[];
  publisher?: string;
  publishedDate?: string;
  thumbnailUrl: string;
  status: "available" | "borrowed";
  borrower?: {
    id: string;
    name: string;
    borrowedAt: string;
  };
  createdAt: string;
};
