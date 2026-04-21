import { BrowserRouter, Route, Routes } from "react-router-dom";

import { LobbyRoute } from "./pages/LobbyPage";
import { RoomRoute } from "./pages/RoomPage";

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LobbyRoute />} />
        <Route path="/rooms/:roomId" element={<RoomRoute />} />
      </Routes>
    </BrowserRouter>
  );
}
