import { Bubble, Sender } from "@ant-design/x";
import { App, Flex, Layout } from "antd";
import React, { useState } from "react";
import { User, UserMinus, UserPlus } from "lucide-react";

const { Content, Footer } = Layout;

const ProsConsV1: React.FC = () => {
  const [value, setValue] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const [messages, setMessages] = useState<
    { content: string; role: string; id: string }[]
  >([]);
  const refetchMessages = async () => {
    console.log("ping");
    const updates = await Promise.all(
      messages.map(async (msg) => {
        if (msg.content !== "loading...") return msg;

        try {
          const response = await fetch(
            `http://localhost:8080/v2/chat/${msg.id}`
          );
          const data = await response.json();
          if (data.status === "Completed") {
            return { ...msg, content: data.result };
          }
        } catch (error) {
          console.error("Error fetching message status:", error);
        }
        return msg;
      })
    );
    setMessages(updates);
  };

  React.useEffect(() => {
    const interval = setInterval(() => {
      refetchMessages();
    }, 3000);

    return () => clearInterval(interval);
  }, [messages]);

  const { message: notif } = App.useApp();

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
    setMessages((prev) => [
      ...prev,
      { content: value, role: "user", id: "user" },
    ]);
    for (const side of ["agree", "disagree"]) {
      const response = await fetch("http://localhost:8080/v2/chat/", {
        method: "POST",
        body: JSON.stringify({ message: value, side: side }),
      });
      const data = await response.json();
      setMessages((prev) => [
        ...prev,
        { content: "loading...", role: side, id: data.wfid },
      ]);
    }
  };

  return (
    <App>
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
                {/* <Bubble
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
                /> */}
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
              notif.error("Cancel sending!");
            }}
            autoSize={{ minRows: 2, maxRows: 6 }}
          />
        </Footer>
      </Layout>
    </App>
  );
};

export default ProsConsV1;
