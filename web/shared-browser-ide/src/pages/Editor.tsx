import { useParams, useLocation } from "react-router-dom";
import { useEffect, useState } from "react";
import SessionInfo from "@/components/SessionInfo";
import CodeEditor from "@monaco-editor/react";
import ErrorModal from "@/components/ErrorModal";

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
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

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

  function connectWebsocket() {
    const ws = new WebSocket(
      `ws://localhost:8080/session/${id}/connect/ws?name=${name}&password=${password}`
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
      console.log("Websocket error: ", err);

      setErrorMessage("Oops! An error occured connecting to the server.");
    };

    ws.onmessage = (message) => onWsMessage(JSON.parse(message.data));

    ws.addEventListener;
  }

  async function verifySession() {
    const res = await fetch(
      `http://localhost:8080/session/${id}/connect?name=${name}&password=${password}`,
      {
        method: "GET",
        headers: { "Content-Type": "application/json" },
      }
    );

    return res.status;
  }

  useEffect(() => {
    debugger;
    verifySession().then((status) => {
      switch (status) {
        case 200:
          connectWebsocket();
          break;
        case 403:
          setErrorMessage("Invalid session password, try again buddy");
          break;
        case 404:
          setErrorMessage("Invalid name or session ID, back to the lobby");
          break;
        default:
          setErrorMessage("Unknown error. Tough luck buddy");
      }

      return () => {
        socket?.close(); // clean up on unmount
      };
    });
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
      <ErrorModal
        isOpen={!!errorMessage}
        message={errorMessage ?? ""}
        onClose={() => setErrorMessage(null)}
      />
    </div>
  );
}
