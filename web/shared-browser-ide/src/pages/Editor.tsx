import { useParams, useLocation } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
import SessionInfo from "@/components/SessionInfo";
import CodeEditor from "@monaco-editor/react";
import ErrorModal from "@/components/ErrorModal";

export default function Editor() {
  const { id } = useParams<{ id: string }>();
  const location = useLocation();
  const { name, password, language, isSessionCreator } = location.state || {};
  const socketRef = useRef<WebSocket | null>(null);
  const [editorContent, setEditorContent] = useState<string>(
    `// Welcome ${name}`
  );
  const [editorProgrammingLanguage, setEditorProgrammingLanguage] =
    useState<string>(language);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [connectedClients, setConnectedClients] = useState<string[]>([]);
  const [isClientListOpen, setIsClientListOpen] = useState(false);

  type Message = {
    type: string;
    editorContent: string | undefined;
    programmingLanguage: string | undefined;
    clients: string[] | undefined;
  };

  function onEditorChange(editorContent: string | undefined, _: any) {
    if (editorContent === undefined) return;
    setEditorContent(editorContent);
    if (socketRef.current?.readyState === WebSocket.OPEN) {
      const wsMessage: Message = createWsMessage(
        editorContent,
        editorProgrammingLanguage,
        undefined
      );
      socketRef?.current?.send(JSON.stringify(wsMessage));
    }
  }

  function handleCodeUpdateMessage(message: Message) {
    setEditorContent(message.editorContent ?? "");
    setEditorProgrammingLanguage(message.programmingLanguage ?? "");
  }

  function handleClientsUpdateMessage(message: Message) {
    setConnectedClients(message.clients ?? []);
  }

  function onWsMessage(message: Message) {
    if (message.type === "SessionCodeUpdate") {
      handleCodeUpdateMessage(message);
    }

    if (message.type === "ClientsUpdate") {
      handleClientsUpdateMessage(message);
    }
  }

  function createWsMessage(
    content: string,
    programmingLanguage: string,
    clients: string[] | undefined
  ): Message {
    return {
      type: "ClientCodeUpdate",
      editorContent: content,
      programmingLanguage,
      clients,
    };
  }

  function connectWebsocket() {
    const ws = new WebSocket(
      `ws://localhost:8080/api/session/${id}/connect/ws?name=${name}&password=${password}`
    );

    socketRef.current = ws;
    ws.onopen = () => {
      console.log("WebSocket connection established.");

      if (isSessionCreator) {
        ws.send(
          JSON.stringify(
            createWsMessage(editorContent, editorProgrammingLanguage, undefined)
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
      `http://localhost:8080/api/session/${id}/connect?name=${name}&password=${password}`,
      {
        method: "GET",
        headers: { "Content-Type": "application/json" },
      }
    );

    return res.status;
  }

  useEffect(() => {
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
    });
    return () => {
      socketRef.current?.close(); // clean up on unmount
    };
  }, [id, name]);

  return (
    <div className="min-h-screen p-6 bg-gray-1000">
      <SessionInfo sessionId={id ?? ""} password={password} />
      <div className="absolute top-4 right-4 z-50">
        <button
          onClick={() => setIsClientListOpen(!isClientListOpen)}
          className="bg-blue-600 text-white px-3 py-1 rounded hover:bg-blue-700"
        >
          {isClientListOpen ? "Hide Participants" : "Show Participants"}
        </button>

        {isClientListOpen && (
          <div className="mt-2 bg-white text-black rounded shadow-lg p-4 max-w-xs">
            <ul className="text-sm max-h-60 overflow-y-auto space-y-1">
              {connectedClients.length > 0 ? (
                connectedClients.map((client, idx) => (
                  <li key={idx} className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                    {client}
                  </li>
                ))
              ) : (
                <li className="text-gray-500">No one is here...</li>
              )}
            </ul>
          </div>
        )}
      </div>
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
