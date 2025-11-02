FROM node:20-alpine

WORKDIR /usr/src/app

COPY package*.json ./
RUN npm ci --only=production

COPY . .

ENV PORT=${PORT}
EXPOSE ${PORT}

CMD ["node", "app.js"]
