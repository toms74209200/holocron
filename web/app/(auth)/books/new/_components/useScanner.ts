import { Html5Qrcode } from "html5-qrcode";
import { useEffect, useRef, useState } from "react";

export type ScannerState =
  | { status: "idle" }
  | { status: "initializing" }
  | { status: "scanning" }
  | { status: "cooldown" }
  | { status: "error"; message: string };

const COOLDOWN_MS = 1000;
const INITIALIZATION_DELAY_MS = 100;

const CAMERA_CONSTRAINTS = { facingMode: "environment" } as const;

const SCANNER_CONFIG = {
  fps: 10,
  qrbox: { width: 250, height: 100 },
  aspectRatio: 1.0,
} as const;

const ERROR_MESSAGES = {
  ELEMENT_NOT_FOUND: "スキャナー要素が見つかりません",
  CAMERA_FAILED: "カメラの起動に失敗しました",
} as const;

export function useScanner(
  elementId: string,
  enabled: boolean,
  onScan: (code: string) => void,
) {
  const scannerRef = useRef<Html5Qrcode | null>(null);
  const [state, setState] = useState<ScannerState>({ status: "idle" });

  useEffect(() => {
    if (state.status !== "cooldown") {
      return;
    }

    const timer = setTimeout(() => {
      setState({ status: "scanning" });
    }, COOLDOWN_MS);

    return () => clearTimeout(timer);
  }, [state.status]);

  useEffect(() => {
    if (!enabled) {
      return;
    }

    // Flag to cancel async operations after cleanup
    // Prevents double initialization in React Strict Mode
    let cancelled = false;

    const startScanner = async () => {
      setState({ status: "initializing" });

      await new Promise((resolve) =>
        setTimeout(resolve, INITIALIZATION_DELAY_MS),
      );

      if (cancelled) {
        return;
      }

      const element = document.getElementById(elementId);
      if (!element) {
        if (cancelled) {
          return;
        }
        setState({
          status: "error",
          message: ERROR_MESSAGES.ELEMENT_NOT_FOUND,
        });
        return;
      }

      element.innerHTML = "";

      const scanner = new Html5Qrcode(elementId);
      scannerRef.current = scanner;

      try {
        await scanner.start(
          CAMERA_CONSTRAINTS,
          SCANNER_CONFIG,
          (decodedText) => {
            setState((prev) => {
              if (prev.status !== "scanning") {
                return prev;
              }

              onScan(decodedText);

              return { status: "cooldown" };
            });
          },
          () => {},
        );

        if (cancelled) {
          return;
        }

        setState({ status: "scanning" });
      } catch (err) {
        if (cancelled) {
          return;
        }
        const errorMessage =
          err instanceof Error ? err.message : ERROR_MESSAGES.CAMERA_FAILED;
        setState({ status: "error", message: errorMessage });
      }
    };

    startScanner();

    return () => {
      cancelled = true;
      const scanner = scannerRef.current;
      if (scanner) {
        const element = document.getElementById(elementId);
        if (element) {
          element.innerHTML = "";
        }
        scanner
          .stop()
          .then(() => scanner.clear())
          .catch(() => {});
        scannerRef.current = null;
      }
    };
  }, [elementId, enabled, onScan]);

  return state;
}
