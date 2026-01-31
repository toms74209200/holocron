"use client";

import {
  onAuthStateChanged,
  signInAnonymously,
  type User,
} from "firebase/auth";
import { createContext, useContext, useEffect, useState } from "react";
import { getFirebaseAuth } from "@/app/_lib/firebase";

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
        await signInAnonymously(auth);
      }
    });

    return () => unsubscribe();
  }, []);

  return (
    <AuthContext.Provider value={authState}>{children}</AuthContext.Provider>
  );
};
