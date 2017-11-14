# Prog-Image

## Installation
```go
go get github.com/vedhavyas/prog-image/cmd/...
```

## Start
```
./prog-imaged --addr ":8080"
```

## API

### Upload Image

`POST /images/`

#### Form Data:

Base64 upload
```
type: base64
image: [base64 encoded image data]

Headers:
Content-Type: application/x-www-form-urlencoded
```

URL Upload
```
type: url
image: [image url]

Headers:
Content-Type: application/x-www-form-urlencoded
```

Multipart upload
```
type: file
image: [image file]

Headers:
Content-Type: multipart/form-data
```

#### JSONResponse

Successful(201)
```
{
  "id": [unique image id]   
}
```

Failed(400, 500)
```
{
  "error": [error reason]   
}
```

### Download Image

Original Image
`Get /images/{image_id}`

Formatted Image
`Get /images/{image_id}?format=[png|jpeg]`
