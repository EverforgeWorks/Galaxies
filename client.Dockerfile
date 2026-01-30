# Step 1: Build Stage
FROM node:20-alpine AS builder
WORKDIR /app

# Copy dependency definitions from the client subfolder
COPY client/package*.json ./
RUN npm install

# Copy the rest of the client source code
COPY client/ .
RUN npm run build

# Step 2: Serve Stage
FROM nginx:alpine

# Copy the build output
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy your custom Nginx config (Vital for WebSocket routing)
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
