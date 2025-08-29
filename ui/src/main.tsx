import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
// import ChatbotV1 from "./ChatbotV1";
import ProsConsV1 from "./ProsConsV1";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ProsConsV1 />
    {/* <ChatbotV1 /> */}
  </StrictMode>
);
