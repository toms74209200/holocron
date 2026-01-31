"use client";

import { useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useAuth } from "@/app/_components/AuthProvider";
import { fetchClient } from "@/app/_lib/query";
import { parseIsbn } from "../../_models/isbn";
import { NewBookPage } from "./page.view";

function NewBookContent() {
  const authState = useAuth();
  const router = useRouter();
  const queryClient = useQueryClient();
  const [code, setCode] = useState("");
  const [error, setError] = useState<string | null>(null);

  if (authState.status !== "authenticated") {
    throw new Error("NewBookContent requires authenticated user");
  }

  const { user } = authState;

  const handleSubmit = async (e: { preventDefault: () => void }) => {
    e.preventDefault();
    const isbn = parseIsbn(code);
    if (!isbn) {
      setError("ISBNの形式が正しくありません");
      return;
    }

    setError(null);

    router.push("/");

    const token = await user.getIdToken();
    const { error: apiError } = await fetchClient.POST("/books/code", {
      body: { code: isbn },
      headers: { Authorization: `Bearer ${token}` },
    });

    if (apiError) {
      console.error("書籍登録に失敗しました:", apiError.message);
    } else {
      queryClient.invalidateQueries({ queryKey: ["books"] });
    }
  };

  return (
    <NewBookPage
      code={code}
      error={error}
      onChangeCode={setCode}
      onSubmit={handleSubmit}
    />
  );
}

export default function NewBook() {
  const authState = useAuth();

  if (authState.status !== "authenticated") {
    return (
      <NewBookPage
        code=""
        error={null}
        onChangeCode={() => {}}
        onSubmit={() => {}}
      />
    );
  }

  return <NewBookContent />;
}
