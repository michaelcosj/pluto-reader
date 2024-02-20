/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["views/**/*.templ"],
    theme: {
        colors: {
            'text': 'rgb(var(--text))',
            'background': 'rgb(var(--background))',
            'primary': 'rgb(var(--primary))',
            'secondary': 'rgb(var(--secondary))',
            'accent': 'rgb(var(--accent))',
        },
    },
    plugins: [
        require('@tailwindcss/forms'),
        require('@tailwindcss/typography'),
    ],
}

