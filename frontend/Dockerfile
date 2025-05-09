FROM node:18-alpine AS development

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

EXPOSE 3000

ENV HOSTNAME "0.0.0.0"
ENV PORT 3000

CMD ["npm", "run", "dev", "--", "-H", "0.0.0.0"]

FROM node:18-alpine AS production

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build

EXPOSE 3000

ENV HOSTNAME "0.0.0.0"
ENV PORT 3000

CMD ["npm", "start"] 