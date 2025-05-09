/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    serverComponentsExternalPackages: [],
  },
  webpack: (config, { isServer }) => {
    // 개발 서버의 메모리 제한 증가
    if (!isServer) {
      config.optimization.splitChunks = {
        chunks: 'all',
        maxSize: 244000,
        minSize: 20000,
      };
    }
    return config;
  },
}

module.exports = nextConfig 