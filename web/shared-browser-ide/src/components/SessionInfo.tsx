import { useState } from "react";
import { Eye, EyeOff, Lock, Clipboard, Check } from "lucide-react";

type SessionProps = {
  sessionId: string;
  password: string;
};

export default function SessionInfo({ sessionId, password }: SessionProps) {
  const [showPassword, setShowPassword] = useState(false);
  const [copied, setCopied] = useState<null | "id" | "password">(null);

  const copyToClipboard = async (text: string, field: "id" | "password") => {
    try {
      await navigator.clipboard.writeText(text);
      setCopied(field);
      setTimeout(() => setCopied(null), 1500);
    } catch {
      // optionally handle error here
    }
  };

  return (
    <div className="p-4 bg-gray-100 rounded-xl shadow inline-flex">
      {/* Info column */}
      <div className="flex flex-col gap-2 flex-grow">
        {/* Session ID with copy */}
        <div className="text-sm text-gray-700 flex items-center gap-2">
          <strong>Session ID:</strong> {sessionId}
          <button
            onClick={() => copyToClipboard(sessionId, "id")}
            aria-label="Copy session ID"
            className="text-black p-0 copy-button hover:text-gray-700 transition"
            type="button"
          >
            {copied === "id" ? (
              <Check className="w-4 h-4" />
            ) : (
              <Clipboard className="w-4 h-4" />
            )}
          </button>
        </div>

        {/* Password with icon and copy */}
        <div className="text-sm text-gray-700 flex items-center gap-2">
          <Lock className="w-4 h-4 text-gray-600" />
          <span className="font-semibold">Password:</span>
          {showPassword ? (
            <span className="ml-1 text-blue-700 font-mono">{password}</span>
          ) : (
            <span className="ml-1 text-gray-400 italic">hidden</span>
          )}
          <button
            onClick={() => copyToClipboard(password, "password")}
            aria-label="Copy password"
            className="text-black p-0 copy-button hover:text-gray-700 transition bg-color-white-900"
            type="button"
          >
            {copied === "password" ? (
              <Check className="w-4 h-4" />
            ) : (
              <Clipboard className="w-4 h-4" />
            )}
          </button>
        </div>
      </div>

      {/* Show/hide password button */}
      <button
        onClick={() => setShowPassword(!showPassword)}
        className="ml-4 self-stretch flex items-center text-white-700 hover:text-blue-900 transition"
        aria-label="Toggle password visibility"
        type="button"
      >
        {showPassword ? (
          <EyeOff className="w-5 h-5" />
        ) : (
          <Eye className="w-5 h-5" />
        )}
      </button>
    </div>
  );
}
