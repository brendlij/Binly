## Public Endpoints

All endpoints return JSON unless otherwise specified.

### Create Paste

**POST** `/api/p`

Create a new paste.

**Parameters** (form data):

- `content` (required) - Paste content (max 512KB)
- `syntax` - Language/syntax (default: `auto`)
- `ttl` - Time to live: `15m`, `2h`, `1d`, or `0` for never (default: never)
- `allow_edit` - Allow others to edit (default: false)
- `password` - Optional password protection

**Response:**

```json
{
  "id": "abc123def"
}
```

### Get Paste

**GET** `/api/p/{id}`

Retrieve a paste by ID.

**Response:**

```json
{
  "id": "abc123def",
  "content": "console.log('hello');",
  "syntax": "javascript",
  "allow_edit": false,
  "expires_at": 0
}
```

### Edit Paste

**POST** `/api/p/{id}`

Update a paste (only if `allow_edit` is enabled).

**Parameters** (form data):

- `content` (required) - New content

**Response:**

- `204 No Content` on success

### Authenticate

**POST** `/api/p/{id}/auth`

Submit password for protected pastes.

**Parameters** (form data):

- `password` - The password

**Response:**

- `204 No Content` on success
- Sets `auth_{id}` cookie for future access

### Get Raw

**GET** `/api/raw/{id}`

Get plain text content (no JSON wrapper).

**Response:**

```
console.log('hello');
```

## Health Check

**GET** `/_health`

Check if the server is running.

**Response:**

- `200 OK` if healthy
