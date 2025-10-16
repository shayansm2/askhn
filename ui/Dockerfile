FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile
COPY . .
RUN yarn build

FROM node:20-alpine
RUN npm install -g serve
WORKDIR /app
COPY --from=builder /app/dist .
EXPOSE 3000
CMD ["serve", "-s", ".", "-l", "5173"]
