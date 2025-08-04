import { useParams, useLocation } from "react-router-dom";
import SessionInfo from "@/components/SessionInfo";
import CodeEditor from "@monaco-editor/react";

export default function Editor() {
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const { name, password } = location.state;

  return (
    <div className="min-h-screen p-6 bg-gray-1000">
      <SessionInfo sessionId={id ?? ""} password={password} />
      {/* Your editor component or iframe goes here */}
      <div className="mt-4">
        <CodeEditor
          height="80vh"
          defaultLanguage="java"
          defaultValue={`// Welcome ${name}`}
        />
      </div>
    </div>
  );
}
