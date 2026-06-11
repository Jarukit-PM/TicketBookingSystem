FROM node:22-alpine AS build

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM alpine:3.20

RUN apk add --no-cache rsync

COPY --from=build /app/dist /dist

VOLUME /output

CMD ["sh", "-c", "rsync -a --delete /dist/ /output/"]
