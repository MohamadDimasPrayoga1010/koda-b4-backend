package libs

import (
    "context"
    "errors"
    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
    "os"
)

var cld *cloudinary.Cloudinary

func InitCloudinary() (*cloudinary.Cloudinary, error) {
    if cld != nil {
        return cld, nil
    }

    cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
    apiKey := os.Getenv("CLOUDINARY_API_KEY")
    apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

    if cloudName == "" || apiKey == "" || apiSecret == "" {
        return nil, errors.New("Cloudinary credentials not set")
    }

    var err error
    cld, err = cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
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
