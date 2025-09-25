FROM node:20-alpine AS build
WORKDIR /app
ARG VITE_GATEWAY_URL
COPY package.json package-lock.json* ./
RUN npm ci
COPY . .
RUN npm run build

FROM gcr.io/distroless/nodejs20-debian12
WORKDIR /app
COPY --from=build /app/dist ./dist
COPY --from=build /app/node_modules ./node_modules

EXPOSE 5173
CMD ["serve", "-s", "dist"]
