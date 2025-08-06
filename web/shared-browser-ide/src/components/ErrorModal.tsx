import React from "react";

type ErrorModalProps = {
  isOpen: boolean;
  message: string;
  onClose: () => void;
};

const ErrorModal: React.FC<ErrorModalProps> = ({
  isOpen,
  message,
  onClose,
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed top-1/2 left-1/2 z-50 transform -translate-x-1/2 -translate-y-1/2 w-full max-w-md bg-red-100 text-red-900 border border-red-300 p-6 rounded-lg shadow-lg">
      <h2 className="text-xl font-semibold mb-2">Connection Error</h2>
      <p className="mb-4">{message}</p>
      <div className="flex justify-end">
        <button
          className="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
          onClick={onClose}
        >
          Close
        </button>
      </div>
    </div>
  );
};

export default ErrorModal;
