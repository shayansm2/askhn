# temporallm

## projecr ideas

| application                   | data source                                                        | is RAG | hash knowledge base | can measure performance |
| ----------------------------- | ------------------------------------------------------------------ | ------ | ------------------- | ----------------------- |
| hectora chatbot               | https://docs.hectora.cloud/                                        |
|                               | https://github.com/HackerNews/API                                  |
| system design interviwer      | Acing the System Design Interview book                             |
| datatalks club project grader | https://github.com/DataTalksClub/llm-zoomcamp/blob/main/project.md |
| simple notebooklm             | -                                                                  |

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
