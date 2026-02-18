# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci --legacy-peer-deps

# Copy source code
COPY . .

# Set memory limit
ENV NODE_OPTIONS="--max-old-space-size=4096"

# Run type-check with verbose output
RUN npx vue-tsc --noEmit --pretty 2>&1 || (echo "=== TypeScript Errors ===" && exit 1)

# Then build (vite only, no type-check)
RUN npm run build-only

# Runtime stage with nginx
FROM nginx:alpine

# Copy built files to nginx
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
