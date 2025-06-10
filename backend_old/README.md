# URL Shortener Backend

This is a minimal FastAPI backend stub for the URL shortener exercise.
*This is a stub backend. It’s meant for testing frontend integration and demonstrating API behavior. It is not a production-ready implementation. It returns hardcoded responses and is provided for representative purposes only.

> Happy building! 🚀

## Running Locally

1. Install dependencies:
```bash
pip install --no-cache-dir -r requirements.txt
```

2. Start the server:
```bash
uvicorn main:app --reload --port 5000
```
The backend will be available at [http://localhost:5000](http://localhost:5000).

## API Endpoints (Stubbed)

### `POST /api/shorten`
- Accepts a JSON body with a URL.
- Responds with a hardcoded short code and short URL.

Example request:
```json
{ "url": "https://www.example.com" }
```

Example response:
```json
{
  "short_code": "abc123",
  "short_url": "http://localhost:5000/abc123",
  "original_url": "http://example.com"
}
```

### `GET /{short_code}`
- Redirects to a hardcoded URL (`https://example.com`) if the code is `abc123`.
- Returns `404` if not found.

## Development Notes
- CORS is enabled for `http://localhost:3000` to support interaction with the frontend.
- The backend currently uses an in-memory stub dictionary and does not persist data.
