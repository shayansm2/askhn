import { Bubble, Sender, useXAgent, useXChat } from "@ant-design/x";
import { Flex, type GetProp } from "antd";
import React from "react";
import ReactMarkdown from "react-markdown";
import { User, Computer } from "lucide-react";

const roles: GetProp<typeof Bubble.List, "roles"> = {
  user: {
    placement: "end",
    avatar: { icon: <User />, style: { background: "blue" } },
  },
  ai: {
    placement: "start",
    avatar: { icon: <Computer />, style: { background: "black" } },
    typing: { step: 5, interval: 20 },
    style: {
      maxWidth: 600,
    },
  },
};

const ChatbotV1 = () => {
  const [content, setContent] = React.useState("");

  const [agent] = useXAgent<string, { message: string }, string>({
    request: async ({ message }, { onSuccess, onError }) => {
      try {
        const apiUrl = import.meta.env.BACKEND_URL || "http://localhost:8080";
        const response = await fetch(`${apiUrl}/v3/chat?message=${message}`);
        const data = await response.json();
        onSuccess(data.result);
      } catch (error) {
        console.error(error);
        onError(new Error("request failed"));
      }
    },
  });

  // Chat messages
  const { onRequest, messages } = useXChat({
    agent: agent,
    requestPlaceholder: "Waiting...",
    requestFallback: "Mock failed return. Please try again later.",
  });

  return (
    <Flex vertical gap="middle">
      <Bubble.List
        roles={roles}
        items={messages.map(({ id, message, status }) => ({
          key: id,
          loading: status === "loading",
          role: status === "local" ? "user" : "ai",
          content: <ReactMarkdown>{message}</ReactMarkdown>,
        }))}
      />
      <Sender
        loading={agent.isRequesting()}
        value={content}
        onChange={setContent}
        onSubmit={(nextContent) => {
          onRequest(nextContent);
          setContent("");
        }}
      />
    </Flex>
  );
};

export default ChatbotV1;
