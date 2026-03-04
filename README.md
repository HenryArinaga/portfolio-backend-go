# portfolio-backend-go

Go backend for a blog platform with admin dashboard, markdown support, and image uploads.

## Features

- 🔐 Admin authentication with session management
- 📝 Create, edit, and delete blog posts
- 🖼️ Image upload support (JPG, PNG, GIF, WEBP)
- 📄 Markdown content support
- 🎯 REST API for frontend integration
- 🔒 CORS configuration for secure cross-origin requests
- ⚡ Rate limiting on admin endpoints
- 💾 SQLite database

## Quick Start

1. Install dependencies:
```bash
go mod download
```

2. Configure environment variables (`.env`):
```env
PORT=8080
DB_PATH=blog.db
ADMIN_PASSWORD=your-secure-password
ALLOWED_ORIGIN=*
```

3. Run database migration (if upgrading):
```bash
go run scripts/migrate_db.go
```

4. Start the server:
```bash
go run cmd/server/main.go
```

Server runs on `http://localhost:8080`

## API Endpoints

See [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) for complete API reference.

### Public Endpoints
- `GET /api/posts` - List all published posts
- `GET /api/posts/{slug}` - Get post by slug
- `GET /api/posts/previews` - Get post previews

### Admin Endpoints (require authentication)
- `POST /api/admin/login` - Admin login
- `POST /api/admin/logout` - Admin logout
- `GET /api/admin/posts` - List all posts
- `POST /api/admin/posts` - Create post (supports multipart/form-data for images)
- `PUT /api/admin/posts/{id}` - Update post
- `POST /api/admin/posts/{id}/delete` - Delete post
- `POST /api/admin/upload` - Upload image

## Project Structure

```
.
├── cmd/server/          # Main application entry point
├── internal/
│   ├── api/            # Public API handlers
│   │   └── admin/      # Admin API handlers
│   ├── auth/           # Session management
│   ├── config/         # Configuration loading
│   ├── db/             # Database schema and connection
│   ├── handlers/       # Authentication handlers
│   ├── middleware/     # CORS, rate limiting, auth middleware
│   ├── storage/        # Image upload handling
│   └── web/            # Server-side rendered pages
├── uploads/images/     # Uploaded images directory
└── scripts/            # Database migration scripts
```

## Database Schema

```sql
CREATE TABLE posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  content TEXT NOT NULL,
  excerpt TEXT,
  image_url TEXT,
  published INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL
);
```

## Frontend Integration

The API is designed to work with modern frontend frameworks. Example with fetch:

```javascript
// Login
await fetch('http://localhost:8080/api/admin/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include',
  body: JSON.stringify({ password: 'your-password' })
});

// Create post with image
const formData = new FormData();
formData.append('title', 'My Post');
formData.append('content', '# Markdown content');
formData.append('excerpt', 'Short description');
formData.append('published', 'true');
formData.append('image', fileInput.files[0]);

await fetch('http://localhost:8080/api/admin/posts', {
  method: 'POST',
  credentials: 'include',
  body: formData
});
```

## Security Features

- Session-based authentication with HttpOnly cookies
- CORS protection with allowed origins
- Rate limiting on admin endpoints (5 req/min)
- File type validation for uploads
- File size limits (10 MB max)
- SQL injection protection via prepared statements

## Development

Build:
```bash
go build -o server ./cmd/server
```

Run:
```bash
./server
```

## License

MIT