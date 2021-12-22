package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", hello)
	e.GET("/buckets", buckets)

	e.Logger.Fatal(e.Start(":1323"))
}

type HelloResponse struct {
	Message string `json:"message"`
}

type BucketsResponse struct {
	Buckets []string `json:"buckets"`
}

func hello(c echo.Context) error {
	return c.JSON(http.StatusOK, HelloResponse{
		Message: "Hello ECS!",
	})
}

func buckets(c echo.Context) error {
	region := os.Getenv("AWS_DEFAULT_REGION")
	if region == "" {
		c.Echo().Logger.Error("Envirionment variable AWS_DEFAULT_REGION is undefined")
		return c.JSON(http.StatusInternalServerError, nil)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		c.Echo().Logger.Error("Failed to load AWS default config.", err.Error())
		return c.JSON(http.StatusInternalServerError, nil)
	}

	s3Client := s3.NewFromConfig(cfg)
	out, err := s3Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		c.Echo().Logger.Error("Failed to list buckets", err.Error())
		return c.JSON(http.StatusInternalServerError, nil)
	}

	var buckets []string
	for _, b := range out.Buckets {
		buckets = append(buckets, aws.ToString(b.Name))
	}

	return c.JSON(http.StatusOK, BucketsResponse{
		Buckets: buckets,
	})
}
