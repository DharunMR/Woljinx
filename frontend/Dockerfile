# BUILD STAGE
FROM node:20-alpine AS builder

# Caching
WORKDIR /app
COPY package*.json ./

RUN npm install

COPY . .
RUN npm run build


# RUNTIME STAGE
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]