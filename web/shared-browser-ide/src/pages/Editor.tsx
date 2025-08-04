import { useParams, useLocation } from "react-router-dom";
import { useEffect, useState } from "react";
import SessionInfo from "@/components/SessionInfo";
import CodeEditor from "@monaco-editor/react";

export default function Editor() {
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const { name, password } = location.state || {};
  const [socket, setSocket] = useState<WebSocket | null>(null);

  function sendMessage(message: string | undefined, _: any) {
    if (socket?.readyState === socket?.OPEN) {
      socket?.send(message ?? "failed");
    }
  }

  //Connect and test if message is sent on editor change
  useEffect(() => {
    const ws = new WebSocket(
      `ws://localhost:8080/session/${id}/ws?name=${name}`
    );

    setSocket(ws);
    ws.onopen = () => {
      console.log("WebSocket connection established.");
    };

    ws.onclose = () => {
      console.log("WebSocket connection closed.");
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    ws.onmessage = (message) => console.log(message);

    return () => {
      ws.close(); // clean up on unmount
    };
  }, [id, name]);

  return (
    <div className="min-h-screen p-6 bg-gray-1000">
      <SessionInfo sessionId={id ?? ""} password={password} />
      {/* Your editor component or iframe goes here */}
      <div className="mt-4">
        <CodeEditor
          height="80vh"
          defaultLanguage="java"
          defaultValue={`// Welcome ${name}`}
          onChange={sendMessage}
        />
      </div>
    </div>
  );
}
