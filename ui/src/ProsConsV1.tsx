import { Bubble, Sender } from "@ant-design/x";
import { App, Flex, Steps, Layout } from "antd";
import React, { useState } from "react";
import { User, UserMinus, UserPlus } from "lucide-react";

const { Content, Footer } = Layout;

const ProsConsV1: React.FC = () => {
  const [value, setValue] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [messages, setMessages] = useState<{ content: string; role: string }[]>(
    [
      { content: "is kubernetes a good choice for my company?", role: "user" },
      { content: "Yes, it is a good choice for my company.", role: "agree" },
      {
        content: "No, it is not a good choice for my company.",
        role: "disagree",
      },
      { content: "thanks for your help!", role: "user" },
    ]
  );

  const { message } = App.useApp();

  React.useEffect(() => {
    if (!loading) return;
    const timer = setTimeout(() => {
      setLoading(false);
      message.success("Send message successfully!");
      setMessages((prev) => [
        ...prev,
        { content: "new msg", role: "agree" },
        { content: "new msg", role: "disagree" },
      ]);
    }, 3000);
    return () => clearTimeout(timer);
  }, [loading, message]);

  const getAvatar = (role: string) => {
    if (role === "user")
      return { icon: <User />, style: { background: "blue" } };
    if (role === "agree")
      return { icon: <UserPlus />, style: { background: "green" } };
    return { icon: <UserMinus />, style: { background: "red" } };
  };

  const submitMessage = async () => {
    setValue("");
    setLoading(true);
    setMessages((prev) => [...prev, { content: value, role: "user" }]);
    const response = await fetch("http://localhost:8080/v2/chat/", {
      method: "POST",
      body: JSON.stringify({ message: value, side: "agree" }),
    });
    const data = await response.json();
    setMessages((prev) => [
      ...prev,
      { content: "loading...", role: "agree", id: data.wfid },
    ]);
    console.log(data);
  };

  return (
    <Layout style={{ height: "100vh" }}>
      {/* Scrollable messages */}
      <Content style={{ overflowY: "auto", padding: 16 }}>
        <Flex gap="middle" vertical>
          {messages.map((msg, idx) => (
            <React.Fragment key={idx}>
              <Bubble
                placement={msg.role === "user" ? "end" : "start"}
                content={msg.content}
                avatar={getAvatar(msg.role)}
              />
              <Bubble
                placement={msg.role === "user" ? "end" : "start"}
                content={
                  <Steps
                    progressDot
                    current={1}
                    size="small"
                    labelPlacement="vertical"
                    items={[
                      { title: "searching KB" },
                      { title: "LLM" },
                      { title: "finish" },
                    ]}
                  />
                }
                avatar={getAvatar(msg.role)}
              />
            </React.Fragment>
          ))}
        </Flex>
      </Content>

      <Footer>
        <Sender
          loading={loading}
          value={value}
          onChange={(v) => setValue(v)}
          onSubmit={submitMessage}
          onCancel={() => {
            setLoading(false);
            message.error("Cancel sending!");
          }}
          autoSize={{ minRows: 2, maxRows: 6 }}
        />
      </Footer>
    </Layout>
  );
};

export default () => (
  <App>
    <ProsConsV1 />
  </App>
);
