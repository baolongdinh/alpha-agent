# Vite + Vue 3 Frontend

## Setup

```bash
npm install
npm run dev
```

## Environment Variables

Copy `.env.example` to `.env` and configure:

```
VITE_API_URL=http://localhost:8080/api
```

## Build for Production

```bash
npm run build
npm run preview
```

## Project Structure

```
src/
├── components/       # Vue components
├── composables/      # Reusable composition functions
├── services/         # API services
├── utils/            # Utility functions
├── App.vue          # Root component
└── main.js          # Entry point
```

## Tech Stack

- Vue 3 (Composition API)
- Vite
- Tailwind CSS
- Axios
- Marked (Markdown parsing)
- DOMPurify (XSS protection)
- Heroicons
