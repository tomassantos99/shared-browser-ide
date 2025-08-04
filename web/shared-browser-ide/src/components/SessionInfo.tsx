import { useState } from "react";
import { Eye, EyeOff, Lock } from "lucide-react";

type SessionProps = {
  sessionId: string;
  password: string;
};

export default function SessionInfo({ sessionId, password }: SessionProps) {
  const [showPassword, setShowPassword] = useState(false);

  return (
    <div className="p-4 bg-gray-100 rounded-xl shadow flex max-w-md">
      <div className="flex flex-col gap-1 flex-grow">
        <div className="text-sm text-gray-700">
          <strong>Session ID:</strong> {sessionId}
        </div>

        <div className="text-sm text-gray-700 flex items-center gap-2">
          <Lock className="w-4 h-4 text-gray-600" />
          <span className="font-semibold">Password:</span>
          {showPassword ? (
            <span className="ml-1 text-blue-700 font-mono">{password}</span>
          ) : (
            <span className="ml-1 text-gray-400 italic">hidden</span>
          )}
        </div>
      </div>

      <button
        onClick={() => setShowPassword(!showPassword)}
        className="ml-4 self-center text-blue-700 hover:text-blue-900 transition"
        aria-label="Toggle password visibility"
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
