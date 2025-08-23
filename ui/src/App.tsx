import { UserOutlined } from "@ant-design/icons";
import { Bubble, Sender, useXAgent, useXChat } from "@ant-design/x";
import { Flex, type GetProp } from "antd";
import React from "react";
import ReactMarkdown from "react-markdown";

const roles: GetProp<typeof Bubble.List, "roles"> = {
  ai: {
    placement: "start",
    avatar: { icon: <UserOutlined />, style: { background: "#fde3cf" } },
    typing: { step: 5, interval: 20 },
    style: {
      maxWidth: 600,
    },
  },
  local: {
    placement: "end",
    avatar: { icon: <UserOutlined />, style: { background: "#87d068" } },
  },
};

const App = () => {
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
    agent: agent,
    requestPlaceholder: "Waiting...",
    requestFallback: "Mock failed return. Please try again later.",
  });

  return (
    <Flex vertical gap="middle">
      <Bubble.List
        roles={roles}
        style={{ maxHeight: 300 }}
        items={messages.map(({ id, message, status }) => ({
          key: id,
          loading: status === "loading",
          role: status === "local" ? "local" : "ai",
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

export default App;
