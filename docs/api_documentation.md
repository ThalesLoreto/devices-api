# Device API Documentation

## Overview

The Device API is a RESTful web service for managing device resources. It provides endpoints for creating, reading, updating, and deleting devices, along with filtering capabilities.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

Currently, the API does not require authentication. This should be implemented for production use.

## Content Type

All requests and responses use JSON format:
```
Content-Type: application/json
```

## Error Handling

The API uses standard HTTP status codes and returns error details in JSON format:

```json
{
  "error": "Bad Request",
  "message": "Detailed error description"
}
```

### Common HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `204 No Content` - Resource deleted successfully
- `400 Bad Request` - Invalid request data
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists
- `500 Internal Server Error` - Server error

## Data Models

### Device

```json
{
  "id": "string (UUID)",
  "name": "string",
  "brand": "string",
  "state": "string (available|in-use|inactive)",
  "creation_time": "string (ISO 8601 timestamp)"
}
```

#### Field Descriptions

- **id**: Unique identifier for the device (UUID format, auto-generated)
- **name**: Human-readable name of the device
- **brand**: Manufacturer or brand of the device
- **state**: Current state of the device (available, in-use, inactive)
- **creation_time**: Timestamp when the device was created (read-only)

#### Business Rules

1. **Creation time** cannot be updated after device creation
2. **Name and brand** cannot be updated if the device state is "in-use"
3. **In-use devices** cannot be deleted
4. All fields except **creation_time** are required for device creation

## API Endpoints

### 1. Create Device

Creates a new device resource.

**Endpoint:** `POST /devices`

**Request Body:**
```json
{
  "name": "iPhone 16",
  "brand": "Apple",
  "state": "available"
}
```

**Response:** `201 Created`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "iPhone 16",
  "brand": "Apple",
  "state": "available",
  "creation_time": "2025-07-16T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input data
- `409 Conflict` - Device with same ID already exists

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 16",
    "brand": "Apple",
    "state": "available"
  }'
```

### 2. Get All Devices

Retrieves all devices or filters by query parameters.

**Endpoint:** `GET /devices`

**Query Parameters:**
- `brand` (optional) - Filter devices by brand
- `state` (optional) - Filter devices by state

**Response:** `200 OK`
```json
[
  {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "name": "iPhone 16",
    "brand": "Apple",
    "state": "available",
    "creation_time": "2024-01-16T10:30:00Z"
  },
  {
    "id": "456e7890-e89b-12d3-a456-426614174001",
    "name": "Galaxy S24",
    "brand": "Samsung",
    "state": "in-use",
    "creation_time": "2024-01-16T11:00:00Z"
  }
]
```

**Examples:**

Get all devices:
```bash
curl http://localhost:8080/api/v1/devices
```

Filter by brand:
```bash
curl http://localhost:8080/api/v1/devices?brand=Apple
```

Filter by state:
```bash
curl http://localhost:8080/api/v1/devices?state=available
```

### 3. Get Single Device

Retrieves a specific device by ID.

**Endpoint:** `GET /devices/{id}`

**Path Parameters:**
- `id` - Device ID (UUID)

**Response:** `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "iPhone 16",
  "brand": "Apple",
  "state": "available",
  "creation_time": "2024-01-16T10:30:00Z"
}
```

**Error Responses:**
- `404 Not Found` - Device not found

**Example:**
```bash
curl http://localhost:8080/api/v1/devices/123e4567-e89b-12d3-a456-426614174000
```

### 4. Update Device (Full Update)

Updates all fields of an existing device.

**Endpoint:** `PUT /devices/{id}`

**Path Parameters:**
- `id` - Device ID (UUID)

**Request Body:**
```json
{
  "name": "iPhone 16 Pro",
  "brand": "Apple",
  "state": "in-use"
}
```

**Response:** `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "iPhone 16 Pro",
  "brand": "Apple",
  "state": "in-use",
  "creation_time": "2024-01-16T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input or business rule violation
- `404 Not Found` - Device not found

**Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/devices/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 16 Pro",
    "brand": "Apple",
    "state": "in-use"
  }'
```

### 5. Update Device (Partial Update)

Updates specific fields of an existing device.

**Endpoint:** `PATCH /devices/{id}`

**Path Parameters:**
- `id` - Device ID (UUID)

**Request Body:** (all fields optional)
```json
{
  "state": "inactive"
}
```

**Response:** `200 OK`
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "iPhone 16",
  "brand": "Apple",
  "state": "inactive",
  "creation_time": "2024-01-16T10:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid input or business rule violation
- `404 Not Found` - Device not found

**Example:**
```bash
curl -X PATCH http://localhost:8080/api/v1/devices/123e4567-e89b-12d3-a456-426614174000 \
  -H "Content-Type: application/json" \
  -d '{
    "state": "inactive"
  }'
```

### 6. Delete Device

Deletes a device by ID.

**Endpoint:** `DELETE /devices/{id}`

**Path Parameters:**
- `id` - Device ID (UUID)

**Response:** `204 No Content`

**Error Responses:**
- `400 Bad Request` - Cannot delete device in use
- `404 Not Found` - Device not found

**Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/devices/123e4567-e89b-12d3-a456-426614174000
```

### 7. Health Check

Checks if the API is running and healthy.

**Endpoint:** `GET /health`

**Response:** `200 OK`
```
OK
```

**Example:**
```bash
curl http://localhost:8080/health
```

## Business Rules and Validations

### Device Creation
- All fields (name, brand, state) are required
- State must be one of: "available", "in-use", "inactive"
- Name and brand cannot be empty strings

### Device Updates
- Creation time cannot be modified
- Name and brand cannot be updated if device state is "in-use"
- State transitions are allowed for all devices
- Empty values are not allowed for name and brand

### Device Deletion
- Devices with state "in-use" cannot be deleted
- Devices with states "available" or "inactive" can be deleted

## Rate Limiting

Currently, no rate limiting is implemented. For production use, consider implementing rate limiting to prevent abuse.

## Pagination

The current implementation returns all results without pagination. For large datasets, pagination should be implemented using query parameters like `limit` and `offset`.

## Sorting

Results are currently sorted by creation time in descending order (newest first). Additional sorting options could be implemented using query parameters.

## Examples

### Complete Workflow Example

1. **Create a device:**
```bash
curl -X POST http://localhost:8080/api/v1/devices \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro",
    "brand": "Apple",
    "state": "available"
  }'
```

2. **Get all devices:**
```bash
curl http://localhost:8080/api/v1/devices
```

3. **Update device state to in-use:**
```bash
curl -X PATCH http://localhost:8080/api/v1/devices/{device-id} \
  -H "Content-Type: application/json" \
  -d '{
    "state": "in-use"
  }'
```

4. **Try to update name (should fail):**
```bash
curl -X PATCH http://localhost:8080/api/v1/devices/{device-id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro M3"
  }'
```

5. **Update state back to available:**
```bash
curl -X PATCH http://localhost:8080/api/v1/devices/{device-id} \
  -H "Content-Type: application/json" \
  -d '{
    "state": "available"
  }'
```

6. **Delete the device:**
```bash
curl -X DELETE http://localhost:8080/api/v1/devices/{device-id}
```
