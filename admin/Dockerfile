# Stage 1: Install dependencies
FROM node:20-alpine as deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm install

# Stage 2: Copy source and run dev server
FROM node:20-alpine
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
EXPOSE 3000
CMD ["npm", "start"] 
