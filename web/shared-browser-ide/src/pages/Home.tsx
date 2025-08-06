import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Listbox,
  ListboxButton,
  ListboxOptions,
  ListboxOption,
} from "@headlessui/react";
import InputModal from "@/components/InputModal";

interface CreateSessionResponse {
  id: string;
  password: string;
}

export default function Home() {
  const [name, setName] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [joinName, setJoinName] = useState("");
  const [joinSessionId, setJoinSessionId] = useState("");
  const [joinSessionPassword, setJoinSessionPassword] = useState("");

  const supportedLanguages = [
    "javascript",
    "typescript",
    "css",
    "less",
    "scss",
    "json",
    "html",
    "xml",
    "php",
    "c#",
    "c++",
    "razor",
    "markdown",
    "diff",
    "java",
    "vb",
    "coffeescript",
    "handlebars",
    "batch",
    "pug",
    "f#",
    "lua",
    "powershell",
    "python",
    "ruby",
    "sass",
    "r",
    "objective-c",
  ];
  const [preferredLanguage, setSelectedOption] = useState(
    supportedLanguages[0]
  );
  const navigate = useNavigate();

  async function createSession(): Promise<void> {
    if (!name.trim()) {
      alert("Please enter your name");
      return;
    }

    try {
      const response = await fetch("http://localhost:8080/session/create");
      if (!response.ok) {
        throw new Error(`Failed to create session: ${response.statusText}`);
      }
      const data: CreateSessionResponse = await response.json();

      navigate(`/session/${data.id}`, {
        state: {
          name: name,
          password: data.password,
          language: preferredLanguage,
          isSessionCreator: true,
        },
      });
    } catch (err: any) {
      setError(err.message || "Unknown error");
      console.log(error);
    }
  }

  function joinSession() {
    if (!joinName.trim()) {
      alert("Please enter your name");
      return;
    }

    if (!joinSessionId.trim()) {
      alert("Please enter the session ID");
      return;
    }

    if (!joinSessionPassword.trim()) {
      alert("Please enter the session password");
      return;
    }

    navigate(`/session/${joinSessionId}`, {
      state: {
        name: joinName,
        password: joinSessionPassword,
        language: "",
        isSessionCreator: false,
      },
    });
  }

  return (
    <div className="min-h-screen w-full flex">
      <div className="w-1/3 bg-gray-1000 p-8 flex flex-col justify-center items-center">
        <h1 className="text-4xl mb-10 font-bold">Create a New Session</h1>
        <input
          type="text"
          placeholder="Your name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="border border-gray-300 rounded px-3 py-2 mb-4 w-full max-w-xs"
        />
        <Listbox value={preferredLanguage} onChange={setSelectedOption}>
          <div className="relative w-64 mb-4 w-full max-w-xs">
            <ListboxButton className="w-full py-2 px-3 border border-gray-1000 rounded-md text-left">
              {preferredLanguage}
            </ListboxButton>
            <ListboxOptions className="absolute mt-1 max-h-40 w-full overflow-auto rounded-md border bg-gray-100 shadow-lg">
              {supportedLanguages.map((option) => (
                <ListboxOption
                  key={option}
                  value={option}
                  className={({ active }) =>
                    `cursor-pointer select-none px-4 py-2 ${
                      active ? "bg-gray-300 text-gray-900" : "text-gray-700"
                    }`
                  }
                >
                  {option}
                </ListboxOption>
              ))}
            </ListboxOptions>
          </div>
        </Listbox>
        <button
          onClick={createSession}
          className="bg-white-600 text-white px-6 py-2 rounded hover:bg-white-700 w-full max-w-xs mb-4"
        >
          Create Session
        </button>
        <h4 className="text-2xl mb-4 font-bold">Or</h4>
        <button
          onClick={() => setShowModal(true)}
          className="bg-white-600 text-white px-6 py-2 rounded hover:bg-white-700 w-full max-w-xs mb-4"
        >
          Join Session
        </button>
      </div>
      <div
        className="w-2/3 bg-cover bg-center"
        style={{
          backgroundImage:
            "url('https://storage.googleapis.com/medium-feed.appspot.com/images%2F9353691196%2F3a8ed441774de-7-melhores-IDEs-e-editores-de-texto-C-para-desenvolvimento.jpg')",
        }}
        aria-label="Decorative image"
      />
      <InputModal
        isOpen={showModal}
        joinName={joinName}
        joinSessionId={joinSessionId}
        joinSessionPassword={joinSessionPassword}
        onClose={() => setShowModal(false)}
        onSubmit={joinSession}
        onJoinNameInput={setJoinName}
        onJoinSessionIdInput={setJoinSessionId}
        onJoinSessionPasswordInput={setJoinSessionPassword}
      ></InputModal>
    </div>
  );
}
