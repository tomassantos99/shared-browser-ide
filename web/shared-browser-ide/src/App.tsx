import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import Editor from "./pages/Editor";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/session/:id" element={<Editor />} />
      </Routes>
    </BrowserRouter>
  );
}
