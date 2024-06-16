# About 

This backend came to life because I got tired of searching for different shady websites to convert my files to another format.
It always was time consuming to find websites that would allow you to do that without subscribing anywhere and the whole advertisements 
on these websites are just annoying on top of that.

## Services

The backend consists of 3 main parts:

1) A convert service that allows you to convert different formats
2) A compress service that allows you to compress certain image files
3) A pdf service that allows you to either split, merge or encrypt your pdfs

## Supported formats

For the convert service, following conversions are supported:
- `image/jpeg, image/webp, image/png, image/svg+xml, image/heif` -> `image/jpeg, image/webp, image/png, image/svg+xml, image/heif, application/pdf`
- `video/mp4` -> `image/gif`
- `video/mp4` -> `audio/mpeg`
- `application/vnd.openxmlformats-officedocument.wordprocessingml.document` -> `application/pdf`
- `application/pdf` -> `application/vnd.openxmlformats-officedocument.wordprocessingml.document`

For the compress service, following mimeTypes are supported:
- `image/jpeg`
- `image/webp`
- `image/png`
- `application/pdf`

