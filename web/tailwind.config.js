/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        Poppins: ["Poppins"],
      },
      colors: { Text: "#0300A1", Background: "#F0ECFF", acsent: "#CF3DFD" },
    },
  },
  plugins: [],
};
