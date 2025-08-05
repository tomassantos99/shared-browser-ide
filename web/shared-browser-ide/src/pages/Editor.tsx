import { useParams, useLocation } from "react-router-dom";
import { useEffect, useState } from "react";
import SessionInfo from "@/components/SessionInfo";
import CodeEditor from "@monaco-editor/react";

export default function Editor() {
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const { name, password, language, isSessionCreator } = location.state || {};
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [editorContent, setEditorContent] = useState<string>(
    `// Welcome ${name}`
  );
  const [editorProgrammingLanguage, setEditorProgrammingLanguage] =
    useState<string>(language);

  type Message = {
    type: string;
    content: string;
    programmingLanguage: string;
  };

  function onEditorChange(editorContent: string | undefined, _: any) {
    if (editorContent === undefined) return;
    setEditorContent(editorContent);
    if (socket?.readyState === socket?.OPEN) {
      const wsMessage: Message = createWsMessage(
        editorContent,
        editorProgrammingLanguage
      );
      socket?.send(JSON.stringify(wsMessage));
    }
  }

  function onWsMessage(message: Message) {
    if (message.type === "SessionCodeUpdate") {
      setEditorContent(message.content);
      setEditorProgrammingLanguage(message.programmingLanguage);
    }
  }

  function createWsMessage(
    content: string,
    programmingLanguage: string
  ): Message {
    return {
      type: "ClientCodeUpdate",
      content,
      programmingLanguage,
    };
  }

  useEffect(() => {
    const ws = new WebSocket(
      `ws://localhost:8080/session/${id}/ws?name=${name}`
    );

    setSocket(ws);
    ws.onopen = () => {
      console.log("WebSocket connection established.");

      if (isSessionCreator) {
        ws.send(
          JSON.stringify(
            createWsMessage(editorContent, editorProgrammingLanguage)
          )
        );
      }
    };

    ws.onclose = () => {
      console.log("WebSocket connection closed.");
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
    };

    ws.onmessage = (message) => onWsMessage(JSON.parse(message.data));

    ws.addEventListener;

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
          onChange={onEditorChange}
          value={editorContent}
          language={editorProgrammingLanguage}
        />
      </div>
    </div>
  );
}
