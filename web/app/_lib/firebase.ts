import { initializeApp } from "firebase/app";
import { connectAuthEmulator, getAuth } from "firebase/auth";

const projectId = process.env.NEXT_PUBLIC_PROJECT_ID || "holocron";
const apiKey = process.env.NEXT_PUBLIC_API_KEY || "dummy";
const authDomain = process.env.NEXT_PUBLIC_AUTH_DOMAIN;
const authEmulatorHost = process.env.NEXT_PUBLIC_FIREBASE_AUTH_EMULATOR_HOST;

initializeApp({
  projectId,
  apiKey,
  authDomain,
});

const auth = getAuth();

if (authEmulatorHost) {
  connectAuthEmulator(auth, `http://${authEmulatorHost}`);
}

export const getFirebaseAuth = () => auth;
