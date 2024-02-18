/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["internal/views/**/*.templ"],
    theme: {
        extend: {},
    },
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/typography'),
    ],
}

