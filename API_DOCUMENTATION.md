# Blog API Documentation

## Overview
REST API for a blog platform with admin authentication, post management, and image uploads.

## Base URL
- Development: `http://localhost:8080`
- Production: Configure via `ALLOWED_ORIGIN` env variable

## Authentication
Admin endpoints require session-based authentication. Login first to receive a session cookie.

---

## Public Endpoints

### GET /api/posts
Get all published blog posts.

**Response:**
```json
[
  {
    "id": 1,
    "title": "My First Post",
    "slug": "my-first-post",
    "content": "Full markdown content...",
    "excerpt": "Short preview text",
    "image_url": "/uploads/images/1234567890.jpg",
    "created_at": "2024-03-04T10:00:00Z",
    "published": true
  }
]
```

### GET /api/posts/{slug}
Get a specific published post by slug.

**Response:**
```json
{
  "id": 1,
  "title": "My First Post",
  "slug": "my-first-post",
  "content": "Full markdown content...",
  "excerpt": "Short preview text",
  "image_url": "/uploads/images/1234567890.jpg",
  "created_at": "2024-03-04T10:00:00Z",
  "published": true
}
```

### GET /api/posts/previews
Get post previews (limited content, default 3 posts).

**Response:**
```json
[
  {
    "id": 1,
    "title": "My First Post",
    "slug": "my-first-post",
    "excerpt": "Short preview or first 200 chars...",
    "image_url": "/uploads/images/1234567890.jpg",
    "created_at": "2024-03-04T10:00:00Z"
  }
]
```

### GET /uploads/images/{filename}
Serve uploaded images.

---

## Admin Endpoints

### POST /api/admin/login
Authenticate as admin.

**Request (JSON):**
```json
{
  "password": "your-admin-password"
}
```

**Response:**
```json
{
  "ok": true
}
```

Sets `admin_session` cookie for subsequent requests.

### POST /api/admin/logout
Logout and destroy session.

**Response:** Redirects to `/blog`

### GET /api/admin/posts
List all posts (including unpublished).

**Headers:**
- Cookie: `admin_session=...`

**Response:**
```json
[
  {
    "id": 1,
    "title": "My First Post",
    "slug": "my-first-post",
    "content": "Full markdown content...",
    "excerpt": "Short preview text",
    "image_url": "/uploads/images/1234567890.jpg",
    "created_at": "2024-03-04T10:00:00Z",
    "published": true
  }
]
```

### POST /api/admin/posts
Create a new blog post.

**Headers:**
- Cookie: `admin_session=...`
- Content-Type: `multipart/form-data` (for image upload) OR `application/json`

**Request (multipart/form-data):**
```
title: My New Post
content: # Markdown content here...
excerpt: Optional short description
published: true
image: [file upload]
```

**Request (JSON):**
```json
{
  "title": "My New Post",
  "content": "# Markdown content here...",
  "excerpt": "Optional short description",
  "image_url": "/uploads/images/1234567890.jpg",
  "published": true
}
```

**Response:**
```json
{
  "id": 2
}
```

### PUT /api/admin/posts/{id}
Update an existing post.

**Headers:**
- Cookie: `admin_session=...`
- Content-Type: `multipart/form-data` OR `application/json`

**Request:** Same as POST /api/admin/posts

**Response:** 204 No Content

### POST /api/admin/posts/{id}/delete
Delete a post and its associated image.

**Headers:**
- Cookie: `admin_session=...`

**Response:** Redirects to `/admin`

### POST /api/admin/upload
Upload an image (for use in markdown editor).

**Headers:**
- Cookie: `admin_session=...`
- Content-Type: `multipart/form-data`

**Request:**
```
image: [file upload]
```

**Response:**
```json
{
  "url": "/uploads/images/1234567890.jpg"
}
```

**Supported formats:** JPG, JPEG, PNG, GIF, WEBP  
**Max size:** 10 MB

---

## Data Models

### Post
```typescript
{
  id: number;
  title: string;
  slug: string;           // Auto-generated from title
  content: string;        // Markdown content
  excerpt?: string;       // Optional preview text
  image_url?: string;     // Optional featured image
  created_at: string;     // ISO 8601 timestamp
  published: boolean;
}
```

---

## Error Responses

All endpoints return appropriate HTTP status codes:

- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Success with no response body
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Authentication required
- `404 Not Found` - Resource not found
- `405 Method Not Allowed` - Wrong HTTP method
- `500 Internal Server Error` - Server error

Error response format:
```
Plain text error message
```

---

## CORS Configuration

Allowed origins:
- `http://localhost:5173` (development)
- `https://henryarin.github.io`
- `https://henryarin.github.io/portfolio`

Credentials (cookies) are supported for authenticated requests.

---

## Rate Limiting

Admin endpoints are rate-limited to 5 requests per minute per IP.

---

## Frontend Integration Example

```javascript
// Login
const login = async (password) => {
  const response = await fetch('http://localhost:8080/api/admin/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ password })
  });
  return response.json();
};

// Create post with image
const createPost = async (formData) => {
  const response = await fetch('http://localhost:8080/api/admin/posts', {
    method: 'POST',
    credentials: 'include',
    body: formData  // FormData with title, content, excerpt, image, published
  });
  return response.json();
};

// Get published posts
const getPosts = async () => {
  const response = await fetch('http://localhost:8080/api/posts');
  return response.json();
};

// Upload image for markdown
const uploadImage = async (file) => {
  const formData = new FormData();
  formData.append('image', file);
  
  const response = await fetch('http://localhost:8080/api/admin/upload', {
    method: 'POST',
    credentials: 'include',
    body: formData
  });
  return response.json();
};
```
