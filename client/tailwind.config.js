/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        mono: ['"VT323"', 'monospace'], // Overrides default mono stack
      },
      colors: {
        'terminal-green': '#33ff00',
        'terminal-black': '#0a0a0a',
      }
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["black"], // Forces dark mode base
  },
}
