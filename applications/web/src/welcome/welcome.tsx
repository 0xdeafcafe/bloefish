import { Box, HStack, Input } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import Markdown from "react-markdown";
import { Button } from "~/components/ui/button";
import { NativeSelectField, NativeSelectRoot } from "~/components/ui/native-select";

interface StreamMessageFull {
  channel_id: string;
  type: 'message_full';
  message_full: string;
  message_fragment: null;
}

interface StreamMessageFragment {
  channel_id: string;
  type: 'message_fragment';
  message_full: null;
  message_fragment: string;
}

type StreamMessage = StreamMessageFull | StreamMessageFragment;

export function Welcome() {
  const [response, setResponse] = useState<string>();
  const [working, setWorking] = useState<boolean>(false);
  const [question, setQuestion] = useState<string>();
  const [model, setModel] = useState<'gpt-4o' | 'gpt-4' | 'o1-preview' | 'o1-mini'>('gpt-4o');
  const [channelId, setChannelId] = useState<string>();

  useEffect(() => {
    let ws = new WebSocket('ws://svc_stream.bloefish.local:4004/ws');

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data) as StreamMessage;

      if (message.channel_id != channelId) return;

      if (message.type === 'message_full') {
        setResponse(message.message_full);
        setWorking(false);
      } else if (message.type === 'message_fragment') {
        setResponse((prev) => {
          if (!prev) return message.message_fragment;
          return prev + message.message_fragment;
        });
      }
    };

    ws.onclose = () => {
      setTimeout(() => {
        ws = new WebSocket('ws://svc_stream.bloefish.local:4004/ws');
      }, 1000);
    };

    return () => {
      ws.close();
    };
  }, [channelId]);

  async function askQuestion() {
    if (!question || working) return;

    setWorking(true);
    setResponse(undefined);

    const userResponse = await fetch('http://svc_user.bloefish.local:4001/rpc/2025-02-12/get_or_create_default_user', { method: 'POST' });
    const { user } = await userResponse.json() as { user: { id: string; } }

    const conversationResponse = await fetch("http://svc_conversation.bloefish.local:4002/rpc/2025-02-12/create_conversation", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        idempotency_key: new Date().toISOString(),
        owner: {
          type: 'user',
          identifier: user.id,
        },
        ai_relay_options: {
          provider_id: 'open_ai',
          model_id: model,
        },
      }),
    });
    const { conversation_id, stream_channel_id } = await conversationResponse.json() as { conversation_id: string; stream_channel_id: string; };

    setChannelId(stream_channel_id);

    await fetch("http://svc_conversation.bloefish.local:4002/rpc/2025-02-12/create_conversation_message", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        conversation_id: conversation_id,
        idempotency_key: new Date().toISOString(),
        message_content: question,
        file_ids: [],
        owner: {
          type: 'user',
          identifier: user.id,
        },
        ai_relay_options: {
          provider_id: 'open_ai',
          model_id: model,
        },
        options: {
          use_streaming: true,
        },
      }),
    });
  }

  return (
    <main>
      <HStack>
        <Input
          placeholder="Question"
          value={question}
          onInput={(e) => setQuestion(e.currentTarget.value)}
          onKeyDown={(e) => e.key === 'Enter' && askQuestion()}
        />
        <Button
          disabled={working}
          onClick={async () => askQuestion()}
        >
          Ask question
        </Button>
        <NativeSelectRoot>
          <NativeSelectField
            value={model}
            onChange={(e) => setModel(e.currentTarget.value as 'gpt-4o' | 'gpt-4' | 'o1-preview' | 'o1-mini')}
          >
            <option value="gpt-4o">GPT 4o</option>
            <option value="gpt-4">GPT 4</option>
            <option value="o1-preview">o1 Preview</option>
            <option value="o1-mini">o1 Mini</option>
          </NativeSelectField>
        </NativeSelectRoot>
      </HStack>

      <br />
      <br />
      
      <span>Response:</span>
      <Box>
        <Markdown>{response}</Markdown>
      </Box>
    </main>
  );
}
