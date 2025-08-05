import React from "react";

type InputModalProps = {
  isOpen: boolean;
  joinName: string;
  joinSessionId: string;
  onJoinNameInput: (joinName: string) => void;
  onJoinSessionIdInput: (joinSessionName: string) => void;
  onClose: () => void;
  onSubmit: () => void;
};

const InputModal: React.FC<InputModalProps> = ({
  isOpen,
  joinName,
  joinSessionId,
  onJoinNameInput,
  onJoinSessionIdInput,
  onClose,
  onSubmit,
}) => {
  if (!isOpen) return null;

  return (
    // Centered floating modal without blocking background
    <div className="fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 z-50 w-full max-w-sm">
      <div className="grayBackground p-6 rounded-lg shadow-xl border border-gray-200">
        <h2 className="text-xl font-semibold mb-2">Join Session</h2>

        <div className="mb-4">
          <label htmlFor="join-name" className="block text-sm font-medium mb-1">
            Your Name
          </label>
          <input
            id="join-name"
            type="text"
            value={joinName}
            onChange={(e) => onJoinNameInput(e.target.value)}
            className="w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div className="mb-4">
          <label
            htmlFor="session-id"
            className="block text-sm font-medium mb-1"
          >
            Session ID
          </label>
          <input
            id="session-id"
            type="text"
            value={joinSessionId}
            onChange={(e) => onJoinSessionIdInput(e.target.value)}
            className="w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>

        <div className="flex justify-end gap-2">
          <button
            type="button"
            className="px-4 py-2 bg-gray-200 hover:bg-gray-300 rounded"
            onClick={onClose}
          >
            Cancel
          </button>
          <button
            type="button"
            className="px-4 py-2 bg-blue-600 text-white hover:bg-blue-700 rounded"
            onClick={() => {
              onSubmit();
              onClose();
            }}
          >
            Submit
          </button>
        </div>
      </div>
    </div>
  );
};

export default InputModal;
