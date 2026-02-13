/** @type {import('next').NextConfig} */
const nextConfig = {
    async rewrites() {
        return [
            {
                source: '/api/auth/:path*',
                destination: 'http://localhost:2021/:path*', // Proxy to Auth Service
            },
        ];
    },
};

export default nextConfig;
