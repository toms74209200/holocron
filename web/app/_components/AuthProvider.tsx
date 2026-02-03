"use client";

import {
  onAuthStateChanged,
  signInWithCustomToken,
  type User,
} from "firebase/auth";
import { createContext, useContext, useEffect, useState } from "react";
import { getFirebaseAuth } from "@/app/_lib/firebase";
import { fetchClient } from "@/app/_lib/query";

type AuthState =
  | { status: "initializing" }
  | { status: "signing_in" }
  | { status: "authenticated"; user: User };

const AuthContext = createContext<AuthState>({ status: "initializing" });

export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [authState, setAuthState] = useState<AuthState>({
    status: "initializing",
  });

  useEffect(() => {
    const auth = getFirebaseAuth();

    const unsubscribe = onAuthStateChanged(auth, async (user) => {
      if (user) {
        setAuthState({ status: "authenticated", user });
      } else {
        setAuthState({ status: "signing_in" });
        const { data, error } = await fetchClient.POST("/users", {
          body: {},
        });
        if (error || !data) {
          console.error("Failed to create user:", error);
          return;
        }
        await signInWithCustomToken(auth, data.customToken);
      }
    });

    return () => unsubscribe();
  }, []);

  return (
    <AuthContext.Provider value={authState}>{children}</AuthContext.Provider>
  );
};
