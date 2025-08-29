import { Bubble, Sender, useXAgent, useXChat } from "@ant-design/x";
import { Flex, type GetProp } from "antd";
import React from "react";
import ReactMarkdown from "react-markdown";
import { UserPlus, UserMinus, User } from "lucide-react";

const roles: GetProp<typeof Bubble.List, "roles"> = {
  user: {
    placement: "end",
    avatar: { icon: <User />, style: { background: "blue" } },
  },
  agree: {
    placement: "start",
    avatar: { icon: <UserPlus />, style: { background: "green" } },
    typing: { step: 5, interval: 20 },
    style: {
      maxWidth: 600,
    },
  },
  disagree: {
    placement: "start",
    avatar: { icon: <UserMinus />, style: { background: "red" } },
    typing: { step: 5, interval: 20 },
    style: {
      maxWidth: 600,
    },
  },
};

const ChatbotV1 = () => {
  const [content, setContent] = React.useState("");

  // Agent for request
  const [ollamaAgent] = useXAgent<string, { message: string }, string>({
    request: async ({ message }, { onSuccess, onError }) => {
      try {
        const response = await fetch("http://localhost:11434/api/generate", {
          method: "POST",
          body: JSON.stringify({
            model: "llama3.2",
            prompt: message,
            stream: false,
          }),
        });
        const data = await response.json();
        onSuccess(data.response);
      } catch (error) {
        console.error(error);
        onError(new Error("Mock request failed"));
      }
    },
  });

  const [agent] = useXAgent<string, { message: string }, string>({
    request: async ({ message }, { onSuccess, onError }) => {
      try {
        const response = await fetch(
          "http://localhost:8080/v1/chat?message=" + message
        );
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
    agent: ollamaAgent,
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
          role: status === "local" ? "user" : "agree",
          content: <ReactMarkdown>{message}</ReactMarkdown>,
        }))}
      />
      <Sender
        loading={ollamaAgent.isRequesting()}
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
