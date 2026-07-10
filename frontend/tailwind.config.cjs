/** @type {import('tailwindcss').Config} */
// Tailwind config. darkMode 'class' lets us toggle dark mode by adding/removing
// the `dark` class on <html>, controlled from the ThemeContext.
module.exports = {
  darkMode: "class",
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Inter", "Nunito", "system-ui", "sans-serif"],
      },
      colors: {
        // Indigo accent used for primary actions and active nav.
        brand: {
          50: "#eef2ff",
          100: "#e0e7ff",
          500: "#6366f1",
          600: "#4f46e5",
          700: "#4338ca",
        },
      },
    },
  },
  plugins: [],
};
