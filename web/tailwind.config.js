/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        surface: {
          DEFAULT: '#0f0f0f',
          1: '#161616',
          2: '#1e1e1e',
          3: '#272727',
        },
        accent: {
          DEFAULT: '#a855f7',
          hover: '#9333ea',
        },
      },
    },
  },
  plugins: [],
}
