#FROM node:14-alpine AS builder
#ENV NODE_ENV production
## Add a work directory
#WORKDIR /app
## Cache and Install dependencies
#COPY package.json .
#RUN npm i
#COPY . .
#RUN npm run build
#
#
#FROM nginx:1.21.0-alpine
#ENV NODE_ENV production
## Copy built assets from builder
#COPY --from=builder /app/build /usr/share/nginx/html
## Add your nginx.conf
#COPY nginx.conf /etc/nginx/conf.d/default.conf
## Expose port
#EXPOSE 80
## Start nginx
#CMD ["nginx", "-g", "daemon off;"]


FROM node:16-alpine

WORKDIR /app

COPY package*.json ./

RUN apk add --no-cache python3 py3-pip make g++  && npm install

COPY . .