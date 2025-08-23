# temporallm

## commands

1. elasticsearch in docker

```bash
docker run -it \
    --rm \
    --name elasticsearch \
    -p 9200:9200 \
    -p 9300:9300 \
    -e "discovery.type=single-node" \
    -e "xpack.security.enabled=false" \
    elasticsearch:8.4.3
```

2. run temporal (dev)

```bash
temporal server start-dev
```

3. es queries

- get all sotries

```json
GET /hacker_news/_search/
{
  "query": {
    "bool": {
      "filter": [
        {"term": {
          "type": "story"
        }}
      ]
    }
  }
}
```

- get all comments of a history

```json
GET /hacker_news/_search
{
  "query": {
    "bool": {
      "must": [
        {
          "match_phrase": {
            "title": "Container technologies at Coinbase: Why Kubernetes is not part of our stack"
          }
        }
      ],
      "filter": [
        {"term": {
          "type": "comment"
        }}
      ]
    }
  }
}
```

ollama

```bash
curl http://localhost:11434/api/generate -d '{
  "model": "llama3.2",
  "prompt":"Why is the sky blue?",
  "stream": false
}'
```
