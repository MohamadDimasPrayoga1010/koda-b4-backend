package libs

import (
    "context"
    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
    "os"
)

var cld *cloudinary.Cloudinary

func InitCloudinary() (*cloudinary.Cloudinary, error) {
    if cld != nil {
        return cld, nil
    }

    cldURL := os.Getenv("CLOUDINARY_API_KEY") 
    var err error
    cld, err = cloudinary.NewFromURL(cldURL)
    if err != nil {
        return nil, err
    }
    return cld, nil
}

func UploadFile(file interface{}, folder string) (string, error) {
    cld, err := InitCloudinary()
    if err != nil {
        return "", err
    }

    ctx := context.Background()
    res, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
        Folder: folder,
    })
    if err != nil {
        return "", err
    }

    return res.SecureURL, nil
}